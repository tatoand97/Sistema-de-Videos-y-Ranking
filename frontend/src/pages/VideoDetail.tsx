import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { endpoints } from '@api/client';
import { useAuth } from '@store/auth';

export default function VideoDetail() {
  const { id } = useParams();
  const { token } = useAuth();
  const [data, setData] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      if (!token || !id) return;
      try {
        setLoading(true);
        const d = await endpoints.getVideo(token, id);
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
        <div className="card">
          <pre style={{ whiteSpace: 'pre-wrap' }}>{JSON.stringify(data, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}

