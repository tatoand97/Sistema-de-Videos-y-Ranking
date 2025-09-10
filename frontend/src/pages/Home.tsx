import { useEffect, useState } from 'react';
import { endpoints } from '@api/client';
import { useAuth } from '@store/auth';

type PublicVideo = { video_id: string; title: string; votes?: number };

export default function Home() {
  const [videos, setVideos] = useState<PublicVideo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useAuth();

  useEffect(() => {
    (async () => {
      try {
        setLoading(true);
        const data = await endpoints.publicVideos();
        // backend likely returns an array
        setVideos(Array.isArray(data) ? data : []);
      } catch (e: any) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const vote = async (id: string) => {
    if (!token) { setError('Debes iniciar sesión para votar.'); return; }
    try {
      setError(null);
      await endpoints.voteVideo(token, id);
      // naive feedback
      alert('Voto registrado');
    } catch (e: any) {
      setError(e.message);
    }
  };

  return (
    <div className="container">
      <h2>Videos Públicos</h2>
      {loading && <div className="muted">Cargando…</div>}
      {error && <div className="card" style={{ borderColor: '#553' }}>{error}</div>}
      <div className="grid videos" style={{ marginTop: 12 }}>
        {videos.map(v => (
          <div key={v.video_id} className="card">
            <div className="between">
              <div>
                <div style={{ fontWeight: 600 }}>{v.title}</div>
                {typeof v.votes === 'number' && <div className="muted">Votos: {v.votes}</div>}
              </div>
              <button className="btn" onClick={() => vote(v.video_id)}>Votar</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

