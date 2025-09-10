import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { endpoints } from '@api/client';
import { useAuth } from '@store/auth';

type VideoDetailResponse = {
  video_id: string;
  title: string;
  status: string;
  uploaded_at?: string;
  processed_at?: string | null;
  original_url?: string | null;
  processed_url?: string | null;
};

export default function VideoDetail() {
  const { id } = useParams();
  const { token } = useAuth();
  const [data, setData] = useState<VideoDetailResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      if (!token || !id) return;
      try {
        setLoading(true);
        const d = await endpoints.getVideo<VideoDetailResponse>(token, id);
        setData(d);
      } catch (e: any) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    })();
  }, [token, id]);

  return (
    <div className="container">
      <h2>Detalle del Video</h2>
      {loading && <div className="muted">Cargando…</div>}
      {error && <div className="card" style={{ borderColor: '#553' }}>{error}</div>}
      {data && (
        <div className="card" style={{ maxWidth: 720 }}>
          <div className="field">
            <label>Título</label>
            <input value={data.title} readOnly />
          </div>
          <div className="field">
            <label>Estado</label>
            <input value={data.status} readOnly />
          </div>
          <div className="field">
            <label>Subido</label>
            <input value={data.uploaded_at ? new Date(data.uploaded_at).toLocaleString() : '-'} readOnly />
          </div>
          <div className="field">
            <label>Procesado</label>
            <input value={data.processed_at ? new Date(data.processed_at).toLocaleString() : '-'} readOnly />
          </div>
          <div className="field">
            <label>URL Original</label>
            {data.original_url ? (
              <a className="btn secondary" href={data.original_url} target="_blank" rel="noreferrer">Abrir original</a>
            ) : (
              <input value="-" readOnly />
            )}
          </div>
          <div className="field">
            <label>URL Procesado</label>
            {data.processed_url ? (
              <a className="btn secondary" href={data.processed_url} target="_blank" rel="noreferrer">Abrir procesado</a>
            ) : (
              <input value="-" readOnly />
            )}
          </div>
        </div>
      )}
    </div>
  );
}
