import { useEffect, useState } from 'react';
import { endpoints } from '@api/client';
import Pagination from '@components/Pagination';

type PublicVideo = { video_id: string; title: string; votes?: number; city?: string | null };

export default function Rankings() {
  const [items, setItems] = useState<PublicVideo[]>([]);
  const [city, setCity] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = async () => {
    try {
      setLoading(true);
      const data: any = await endpoints.publicVideos();
      const list: PublicVideo[] = Array.isArray(data) ? data : [];
      // Filtro por ciudad (case-insensitive, ignora tildes básicas)
      const norm = (s: string) => s.normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase();
      const filtered = city
        ? list.filter(v => v.city && norm(v.city) === norm(city))
        : list;
      const sorted = filtered.sort((a, b) => (b.votes || 0) - (a.votes || 0));
      // Paginación en cliente
      const pages = Math.max(1, Math.ceil(sorted.length / pageSize));
      setTotalPages(pages);
      const start = (page - 1) * pageSize;
      const paged = sorted.slice(start, start + pageSize);
      setItems(paged);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, [page, city]);

  return (
    <div className="container">
      <h2>Rankings</h2>
      <div className="row wrap" style={{ gap: 12 }}>
        <div className="field" style={{ minWidth: 260 }}>
          <label>Ciudad (opcional)</label>
          <input placeholder="p.ej. New York" value={city} onChange={e => setCity(e.target.value)} />
        </div>
        <button className="btn" onClick={() => { setPage(1); load(); }}>Buscar</button>
      </div>
      {loading && <div className="muted">Cargando…</div>}
      {error && <div className="card" style={{ borderColor: '#553' }}>{error}</div>}
      <table className="table" style={{ marginTop: 12 }}>
        <thead><tr><th>Título</th><th>Votos</th><th>Ciudad</th></tr></thead>
        <tbody>
          {items.map((it) => (
            <tr key={it.video_id}>
              <td>{it.title}</td>
              <td>{typeof it.votes === 'number' ? it.votes : '-'}</td>
              <td>{it.city || '-'}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <Pagination page={page} totalPages={totalPages} onPage={setPage} />
    </div>
  );
}
