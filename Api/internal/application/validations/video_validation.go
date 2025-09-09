package validations

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/Eyevinn/mp4ff/mp4"
)

const MaxBytes = 100 * 1024 * 1024 // 100 MB

var okBrands = map[string]struct{}{
	"isom": {}, "iso2": {}, "mp41": {}, "mp42": {}, "avc1": {}, "mp4v": {}, "mp71": {},
}

func CheckMP4(b []byte) (int, int, error) {
	// Permite exactamente 100MB; rechaza s�lo si supera el l�mite.
	if len(b) > MaxBytes {
		return 0, 0, fmt.Errorf("excede 100MB (%.2fMB)", float64(len(b))/1024.0/1024.0)
	}

	w, h, ftypMajor, compat, err := mp4VideoDimsAndBrands(b)
	if err != nil {
		return 0, 0, err
	}

	if _, ok := okBrands[ftypMajor]; !ok {
		return 0, 0, fmt.Errorf("brand MP4 no reconocido: major=%s compat=[%s]",
			ftypMajor, strings.Join(compat, ","))
	}

	if w < 1920 || h < 1080 {
		return w, h, fmt.Errorf("resoluci�n insuficiente: %dx%d (<1920x1080)", w, h)
	}

	return w, h, nil
}

func mp4VideoDimsAndBrands(b []byte) (int, int, string, []string, error) {
	r := bytes.NewReader(b)
	f, err := mp4.DecodeFile(r)
	if err != nil {
		return 0, 0, "", nil, fmt.Errorf("no es MP4 v�lido: %w", err)
	}
	if f.Moov == nil || len(f.Moov.Traks) == 0 {
		return 0, 0, "", nil, errors.New("MP4 sin 'moov' o sin 'trak'")
	}
	if f.Ftyp == nil {
		return 0, 0, "", nil, errors.New("MP4 sin 'ftyp'")
	}

	major := f.Ftyp.MajorBrand()        // m�todo, no campo
	compat := f.Ftyp.CompatibleBrands() // m�todo, no campo

	for _, tr := range f.Moov.Traks {
		if w, h, ok := videoDimsFromTrack(tr); ok {
			return w, h, major, compat, nil
		}
	}
	return 0, 0, major, compat, errors.New("no se encontr� track de video")
}

// videoDimsFromTrack intenta extraer dimensiones de un trak de video.
// Devuelve (w, h, true) si logra obtenerlas; de lo contrario, (0, 0, false).
func videoDimsFromTrack(tr *mp4.TrakBox) (int, int, bool) {
	if tr == nil || tr.Mdia == nil || tr.Mdia.Hdlr == nil || tr.Mdia.Hdlr.HandlerType != "vide" {
		return 0, 0, false
	}

	stbl := tr.Mdia.Minf.Stbl
	if stbl != nil && stbl.Stsd != nil && stbl.Stsd.SampleCount > 0 {
		if box, err := stbl.Stsd.GetSampleDescription(0); err == nil {
			if v, ok := box.(*mp4.VisualSampleEntryBox); ok && v != nil {
				return int(v.Width), int(v.Height), true
			}
		}
	}

	if tr.Tkhd != nil {
		return int(tr.Tkhd.Width >> 16), int(tr.Tkhd.Height >> 16), true
	}

	return 0, 0, false
}
