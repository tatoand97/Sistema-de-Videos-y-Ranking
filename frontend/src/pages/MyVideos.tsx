import { useEffect, useState } from 'react';
import { endpoints } from '@api/client';
import { useAuth } from '@store/auth';
import { VideoCard } from '@components/VideoCard';
import type { Video } from '@api/types';
import { Link } from 'react-router-dom';

export default function MyVideos() {
  const { token } = useAuth();
  const [videos, setVideos] = useState<Video[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const data: any = await endpoints.myVideos(token);
      setVideos(Array.isArray(data) ? data : (data.items || []));
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, [token]);

  const remove = async (id: string) => {
    if (!token) return;
    if (!confirm('¿Eliminar este video?')) return;
    try {
      await endpoints.deleteVideo(token, id);
      await load();
    } catch (e: any) {
      setError(e.message);
    }
  };

  return (
    <div className="container">
      <div className="between">
        <h2>Mis Videos</h2>
        <Link to="/upload" className="btn">Subir Video</Link>
      </div>
      {loading && <div className="muted">Cargando…</div>}
      {error && <div className="card" style={{ borderColor: '#553' }}>{error}</div>}
      <div className="grid videos" style={{ marginTop: 12 }}>
        {videos.map(v => (
          <VideoCard key={v.video_id} video={v} actions={<button className="btn danger" onClick={() => remove(v.video_id)}>Eliminar</button>} />
        ))}
      </div>
    </div>
  );
}

