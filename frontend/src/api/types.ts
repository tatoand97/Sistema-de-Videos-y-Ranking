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
  // Additional fields can be added as backend evolves
};

export type RankingItem = {
  video_id: string;
  title: string;
  votes: number;
  city?: string;
};

