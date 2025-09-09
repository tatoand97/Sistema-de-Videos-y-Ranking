import { Routes, Route, Navigate } from 'react-router-dom';
import NavBar from '@components/NavBar';
import Home from '@pages/Home';
import Rankings from '@pages/Rankings';
import Login from '@pages/Login';
import Register from '@pages/Register';
import Profile from '@pages/Profile';
import MyVideos from '@pages/MyVideos';
import Upload from '@pages/Upload';
import VideoDetail from '@pages/VideoDetail';
import ProtectedRoute from '@components/ProtectedRoute';
import { AuthProvider } from '@store/auth';

export default function App() {
  return (
    <AuthProvider>
      <NavBar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/rankings" element={<Rankings />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route element={<ProtectedRoute />}> 
          <Route path="/profile" element={<Profile />} />
          <Route path="/videos" element={<MyVideos />} />
          <Route path="/upload" element={<Upload />} />
          <Route path="/videos/:id" element={<VideoDetail />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </AuthProvider>
  );
}

