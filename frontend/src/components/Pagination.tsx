export default function Pagination({ page, totalPages, onPage }: { page: number; totalPages: number; onPage: (p: number) => void; }) {
  return (
    <div className="row" style={{ justifyContent: 'center', marginTop: 16 }}>
      <button className="btn secondary" disabled={page <= 1} onClick={() => onPage(page - 1)}>Anterior</button>
      <span className="muted" style={{ padding: '0 10px' }}>Página {page} de {totalPages}</span>
      <button className="btn secondary" disabled={page >= totalPages} onClick={() => onPage(page + 1)}>Siguiente</button>
    </div>
  );
}

