export type User = {
  user_id: string;
  first_name: string;
  last_name: string;
  email: string;
  city?: string;
  country?: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

export type Video = {
  video_id: string;
  title: string;
  status: string;
  created_at?: string;
  // Campos opcionales según respuesta del backend
  original_url?: string | null;
  processed_url?: string | null;
};

export type RankingItem = {
  // position puede no venir siempre; el backend lo incluye por página
  position?: number;
  username: string;
  votes: number;
  city?: string | null;
};
