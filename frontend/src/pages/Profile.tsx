import { useAuth } from '@store/auth';

export default function Profile() {
  const { user } = useAuth();
  if (!user) return null;
  return (
    <div className="container">
      <h2>Mi Perfil</h2>
      <div className="card">
        <div><b>Nombre:</b> {user.first_name} {user.last_name}</div>
        <div><b>Email:</b> {user.email}</div>
        <div><b>Ciudad:</b> {user.city || '-'}</div>
        <div><b>País:</b> {user.country || '-'}</div>
      </div>
    </div>
  );
}

