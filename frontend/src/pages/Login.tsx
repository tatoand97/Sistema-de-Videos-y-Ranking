import { FormEvent, useState } from 'react';
import { useAuth } from '@store/auth';
import { useNavigate } from 'react-router-dom';

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);
      await login(email, password);
      navigate('/');
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container" style={{ maxWidth: 560 }}>
      <h2>Iniciar Sesión</h2>
      {error && <div className="card" style={{ borderColor: '#553' }}>{error}</div>}
      <form onSubmit={onSubmit} className="card">
        <div className="field"><label>Email</label><input type="email" value={email} onChange={e => setEmail(e.target.value)} required /></div>
        <div className="field"><label>Contraseña</label><input type="password" value={password} onChange={e => setPassword(e.target.value)} required /></div>
        <button className="btn" disabled={loading}>{loading ? 'Entrando…' : 'Entrar'}</button>
      </form>
    </div>
  );
}

