import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@store/auth';

export default function NavBar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const onLogout = async () => {
    await logout();
    navigate('/');
  };

  return (
    <nav className="navbar">
      <div className="nav-inner">
        <Link className="brand" to="/">TATOLAND • VideoRank</Link>
        <Link to="/rankings">Rankings</Link>
        <Link to="/videos">Mis Videos</Link>
        <div className="spacer" />
        {user ? (
          <>
            <span className="muted">Hola, {user.first_name}</span>
            <Link className="btn secondary" to="/profile">Perfil</Link>
            <button className="btn" onClick={onLogout}>Salir</button>
          </>
        ) : (
          <>
            <Link className="btn secondary" to="/login">Entrar</Link>
            <Link className="btn" to="/register">Registro</Link>
          </>
        )}
      </div>
    </nav>
  );
}

