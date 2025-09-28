import { useEffect, useState } from 'react';
import { endpoints } from '@api/client';
import { useAuth } from '@store/auth';
import type { Video } from '@api/types';
import { Link } from 'react-router-dom';
import Pagination from '@components/Pagination';

export default function MyVideos() {
  const { token } = useAuth();
  const [videos, setVideos] = useState<Video[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [refreshing, setRefreshing] = useState<Record<string, boolean>>({});
  const [publishing, setPublishing] = useState<Record<string, boolean>>({});

  const load = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const data: any = await endpoints.myVideos(token);
      const list: Video[] = Array.isArray(data) ? data : (data.items || []);
      setVideos(list);
      setPage(1);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, [token]);
  useEffect(() => {
    const pages = Math.max(1, Math.ceil(videos.length / pageSize));
    setTotalPages(pages);
    if (page > pages) setPage(1);
  }, [videos, page, pageSize]);

  const refreshOne = async (id: string) => {
    if (!token) return;
    setRefreshing(prev => ({ ...prev, [id]: true }));
    try {
      const d: any = await endpoints.getVideo(token, id);
      setVideos(prev => prev.map(v => v.video_id === id ? ({
        ...v,
        title: d?.title ?? v.title,
        status: d?.status ?? v.status,
        processed_url: d?.processed_url ?? (v as any).processed_url ?? null
      }) : v));
    } catch (e: any) {
      setError(e.message);
    } finally {
      setRefreshing(prev => ({ ...prev, [id]: false }));
    }
  };

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
      <div style={{ marginTop: 12, display: 'flex', flexDirection: 'column', gap: 12 }}>
        {videos.slice((page - 1) * pageSize, (page - 1) * pageSize + pageSize).map(v => {
          const processed = v.status?.toLowerCase() === 'processed' || Boolean((v as any).processed_url);
          return (
            <div key={v.video_id} className="card" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: 12 }}>
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                <div style={{ fontWeight: 700, fontSize: 16 }}>{v.title}</div>
                <div className="muted">Estado: {v.status}{processed && (v as any).processed_url ? ` • listo` : ''}</div>
              </div>
              <div className="row" style={{ gap: 8, alignItems: 'center' }}>
                <button className="btn secondary" onClick={() => refreshOne(v.video_id)} disabled={!!refreshing[v.video_id]}>
                  {refreshing[v.video_id] ? 'Verificando…' : 'Verificar estado'}
                </button>
                <Link className="btn secondary" to={`/videos/${v.video_id}`}>Detalles</Link>
                <a className={`btn secondary${processed ? '' : ' disabled'}`} href={processed ? (v as any).processed_url || '#' : '#'} target="_blank" rel="noreferrer" onClick={e => { if (!processed) e.preventDefault(); }}>Resultado</a>
                <button
                  className="btn secondary"
                  disabled={!processed}
                  onClick={async () => {
                    if (!processed || !token) return;
                    if (!confirm('Publicar este video para que aparezca en la sección pública?')) return;
                    try {
                      setPublishing(prev => ({ ...prev, [v.video_id]: true }));
                      await endpoints.publishVideo(token, v.video_id);
                      await refreshOne(v.video_id);
                      alert('Video publicado exitosamente.');
                    } catch (e: any) {
                      const msg = e?.message || 'Error desconocido';
                      alert(`No se pudo publicar: ${msg}`);
                    } finally {
                      setPublishing(prev => ({ ...prev, [v.video_id]: false }));
                    }
                  }}
                >
                  {publishing[v.video_id] ? 'Publicando…' : 'Listo para publicar'}
                </button>
                <button className="btn danger" onClick={() => remove(v.video_id)}>Eliminar</button>
              </div>
            </div>
          );
        })}
      </div>
      <Pagination page={page} totalPages={totalPages} onPage={setPage} />
    </div>
  );
}
