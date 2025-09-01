package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// OpenCloseVideoProcessingService genera cortinillas (intro/outro) y concatena con el video original.
// Respeta que el agregado total de duración no supere 5 s (por defecto 2.5 s intro + 2.5 s outro).
type OpenCloseVideoProcessingService struct{}

func NewOpenCloseVideoProcessingService() *OpenCloseVideoProcessingService { return &OpenCloseVideoProcessingService{} }

// Process agrega intro/outro con logo y concatena SIN normalizar el video base.
// Se obtienen w/h/fps del video de entrada para generar cortinillas compatibles (mismo w/h/fps).
func (s *OpenCloseVideoProcessingService) Process(inputData []byte, logoPath string, introSeconds, outroSeconds float64, width, height, fps int) ([]byte, error) {
	if introSeconds < 0 { introSeconds = 0 }
	if outroSeconds < 0 { outroSeconds = 0 }
	total := introSeconds + outroSeconds
	if total > 5.0 {
		// Ajuste proporcional para no superar 5 s en total
		if total == 0 { total = 1 }
		scale := 5.0 / total
		introSeconds = introSeconds * scale
		outroSeconds = outroSeconds * scale
	}

	workdir, err := os.MkdirTemp("", "openclose-*")
	if err != nil {
		return nil, fmt.Errorf("tmpdir: %w", err)
	}
	defer os.RemoveAll(workdir)

	inPath := filepath.Join(workdir, "input.mp4")
	if err := os.WriteFile(inPath, inputData, 0o644); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	// --- PROBE: detectar w/h/fps del video base para generar cortinillas compatibles
	vi, err := probeVideo(inPath)
	if err != nil {
		return nil, fmt.Errorf("probe input: %w", err)
	}
	// Ignoramos params entrantes y forzamos a los del video base
	width, height, fps = vi.W, vi.H, vi.FPS
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid input dimensions")
	}
	if fps <= 0 { fps = 30 } // fallback razonable

	introPath := filepath.Join(workdir, "intro.mp4")
	outroPath := filepath.Join(workdir, "outro.mp4")
	outPath   := filepath.Join(workdir, "out.mp4")
	listPath  := filepath.Join(workdir, "concat.txt")

	// --- Generar intro/outro al MISMO w/h/fps del video base (sin tocar el base)
	if introSeconds > 0.0 {
		if err := makeBumper(logoPath, width, height, fps, introSeconds, introPath, true); err != nil {
			return nil, fmt.Errorf("make intro: %w", err)
		}
	}
	if outroSeconds > 0.0 {
		if err := makeBumper(logoPath, width, height, fps, outroSeconds, outroPath, false); err != nil {
			return nil, fmt.Errorf("make outro: %w", err)
		}
	}

	// --- Concatenar: [intro]? + inPath + [outro]? (sin normalizar el inPath)
	var lines bytes.Buffer
	if introSeconds > 0.0 { lines.WriteString(fmt.Sprintf("file '%s'\n", introPath)) }
	lines.WriteString(fmt.Sprintf("file '%s'\n", inPath))
	if outroSeconds > 0.0 { lines.WriteString(fmt.Sprintf("file '%s'\n", outroPath)) }
	if err := os.WriteFile(listPath, lines.Bytes(), 0o644); err != nil {
		return nil, fmt.Errorf("write concat list: %w", err)
	}

	// Intento 1: concat demuxer sin recodificar (rápido y sin pérdida)
	{
		args := []string{
			"-y",
			"-f", "concat",
			"-safe", "0",
			"-i", listPath,
			"-c", "copy",
			outPath,
		}
		if err := runFFmpeg(args); err != nil {
			// Fallback: recodificar en caso de diferencias de codec/params
			args = []string{
				"-y",
				"-f", "concat",
				"-safe", "0",
				"-i", listPath,
				"-c:v", "libx264", "-preset", "veryfast", "-crf", "23",
				"-c:a", "aac", "-b:a", "128k",
				outPath,
			}
			if err2 := runFFmpeg(args); err2 != nil {
				return nil, fmt.Errorf("concat: %v / fallback: %w", err, err2)
			}
		}
	}

	out, err := os.ReadFile(outPath)
	if err != nil { return nil, fmt.Errorf("read out: %w", err) }
	return out, nil
}

// makeBumper genera una cortinilla negra con el logo centrado y fade in/out.
// Se fija al tamaño/fps indicados (los del video base), evitando tocar el video original.
func makeBumper(logoPath string, width, height, fps int, seconds float64, outPath string, isIntro bool) error {
	// Duración de fade (máximo 0.5 s o 20% de la duración)
	fadeDur := seconds * 0.2
	if fadeDur > 0.5 { fadeDur = 0.5 }
	if fadeDur < 0.1 { fadeDur = 0.1 }
	still := seconds

	// Dónde aplicar fade in/out en el logo overlay
	fadeIn  := fmt.Sprintf("fade=t=in:st=0:d=%0.3f", fadeDur)
	fadeOut := fmt.Sprintf("fade=t=out:st=%0.3f:d=%0.3f", seconds-fadeDur, fadeDur)

	// Nota: generamos pista de audio silenciosa para evitar problemas al concatenar
	args := []string{
		"-y",
		"-loop", "1", "-t", fmt.Sprintf("%0.3f", still), "-i", logoPath, // [0:v] logo (imagen)
		"-f", "lavfi", "-t", fmt.Sprintf("%0.3f", still), "-i", "anullsrc=r=48000:cl=stereo", // [1:a] silencio
		"-f", "lavfi", "-t", fmt.Sprintf("%0.3f", still), "-i", fmt.Sprintf("color=c=black:s=%dx%d:r=%d", width, height, fps), // [2:v] fondo negro
		"-filter_complex",
		fmt.Sprintf(
			"[0:v]scale=%d:-1,format=rgba,%s,%s[lg];"+
				"[2:v][lg]overlay=(W-w)/2:(H-h)/2:enable='between(t,0,%0.3f)'[v]",
			int(float64(width)*0.35), // logo ~35%% del ancho del video base
			fadeIn, fadeOut, still,
		),
		"-map", "[v]",
		"-map", "1:a",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "20",
		"-c:a", "aac", "-b:a", "128k",
		"-r", strconv.Itoa(fps),
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		outPath,
	}
	return runFFmpeg(args)
}

// runFFmpeg ejecuta ffmpeg mostrando salida en consola (útil para logs en contenedor).
func runFFmpeg(args []string) error {
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ---- Helpers de probe ----

type videoInfo struct{ W, H, FPS int }

// probeVideo usa ffprobe para obtener width/height/avg_frame_rate del primer stream de video.
func probeVideo(path string) (videoInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,avg_frame_rate",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return videoInfo{}, fmt.Errorf("ffprobe: %w", err)
	}

	// out esperado (3 líneas): width\nheight\nnum/den\n
	lines := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	if len(lines) < 3 {
		return videoInfo{}, fmt.Errorf("ffprobe parse: unexpected output")
	}

	w, _ := strconv.Atoi(string(lines[0]))
	h, _ := strconv.Atoi(string(lines[1]))

	numDen := string(lines[2]) // ej: "30000/1001"
	fps := 30
	if parts := bytes.Split([]byte(numDen), []byte("/")); len(parts) == 2 {
		num, _ := strconv.Atoi(string(parts[0]))
		den, _ := strconv.Atoi(string(parts[1]))
		if den != 0 {
			fpsCalc := float64(num) / float64(den)
			// redondeo simple
			fps = int(fpsCalc + 0.5)
			if fps <= 0 {
				fps = 30
			}
		}
	}

	return videoInfo{W: w, H: h, FPS: fps}, nil
}

// Evita "unused import" si en algún SO no se usa time directamente.
var _ = time.Now