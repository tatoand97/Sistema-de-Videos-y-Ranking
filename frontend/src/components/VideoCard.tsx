import { Link } from 'react-router-dom';
import type { Video } from '@api/types';

export function VideoCard({ video, actions }: { video: Video; actions?: React.ReactNode }) {
  return (
    <div className="card">
      <div className="between">
        <div>
          <div style={{ fontWeight: 600 }}>{video.title}</div>
          <div className="muted">Estado: {video.status}</div>
        </div>
        {actions}
      </div>
      <div className="row wrap" style={{ marginTop: 10 }}>
        <Link className="btn secondary" to={`/videos/${video.video_id}`}>Detalles</Link>
      </div>
    </div>
  );
}

