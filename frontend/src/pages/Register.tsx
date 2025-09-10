import { FormEvent, useState } from 'react';
import { useAuth } from '@store/auth';
import { useNavigate } from 'react-router-dom';

export default function Register() {
  const { register } = useAuth();
  const navigate = useNavigate();

  const [form, setForm] = useState({
    first_name: '', last_name: '', email: '', password1: '', password2: '', city: '', country: ''
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);
      await register(form);
      navigate('/');
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container" style={{ maxWidth: 720 }}>
      <h2>Registro</h2>
      {error && <div className="card" style={{ borderColor: '#553' }}>{error}</div>}
      <form onSubmit={onSubmit} className="card">
        <div className="grid" style={{ gridTemplateColumns: '1fr 1fr', gap: 12 }}>
          <div className="field"><label>Nombre</label><input value={form.first_name} onChange={e => setForm({ ...form, first_name: e.target.value })} required /></div>
          <div className="field"><label>Apellido</label><input value={form.last_name} onChange={e => setForm({ ...form, last_name: e.target.value })} required /></div>
          <div className="field" style={{ gridColumn: 'span 2' }}><label>Email</label><input type="email" value={form.email} onChange={e => setForm({ ...form, email: e.target.value })} required /></div>
          <div className="field"><label>Contraseña</label><input type="password" value={form.password1} onChange={e => setForm({ ...form, password1: e.target.value })} required /></div>
          <div className="field"><label>Repetir Contraseña</label><input type="password" value={form.password2} onChange={e => setForm({ ...form, password2: e.target.value })} required /></div>
          <div className="field"><label>Ciudad</label><input value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} /></div>
          <div className="field"><label>País</label><input value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} /></div>
        </div>
        <button className="btn" disabled={loading}>{loading ? 'Creando…' : 'Crear Cuenta'}</button>
      </form>
    </div>
  );
}

