import { FormEvent, useEffect, useState } from 'react';
import { useAuth } from '@store/auth';
import { endpoints } from '@api/client';

export default function Upload() {
  const { token } = useAuth();
  const [title, setTitle] = useState('');
  const [file, setFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [ok, setOk] = useState(false);

  // No necesitamos cargar statuses: el backend fuerza 'UPLOADED'

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!token) { setError('Necesitas iniciar sesión'); return; }
    if (!file) { setError('Selecciona un archivo'); return; }
    try {
      setLoading(true); setError(null); setOk(false);
      const form = new FormData();
      form.set('title', title);
      // el backend espera el campo 'video_file' con mimetype video/mp4
      form.set('video_file', file);
      await endpoints.uploadVideo(token, form);
      setOk(true);
      setTitle(''); setFile(null);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container" style={{ maxWidth: 720 }}>
      <h2>Subir Video</h2>
      {ok && <div className="card" style={{ borderColor: '#2d5' }}>¡Subido!</div>}
      {error && <div className="card" style={{ borderColor: '#553' }}>{error}</div>}
      <form className="card" onSubmit={onSubmit}>
        <div className="field"><label>Título</label><input value={title} onChange={e => setTitle(e.target.value)} required/></div>
        <div className="field"><label>Archivo</label><input type="file" accept="video/mp4" onChange={e => setFile(e.target.files?.[0] || null)} required/></div>
        <button className="btn" disabled={loading}>{loading ? 'Subiendo…' : 'Subir'}</button>
      </form>
    </div>
  );
}
