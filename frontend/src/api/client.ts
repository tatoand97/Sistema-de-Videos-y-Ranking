const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  body?: any;
  token?: string | null;
  headers?: Record<string, string>;
};

export async function apiFetch<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    ...(opts.headers || {})
  };
  let body: BodyInit | undefined;

  if (opts.body && !(opts.body instanceof FormData)) {
    headers['Content-Type'] = headers['Content-Type'] || 'application/json';
    body = JSON.stringify(opts.body);
  } else if (opts.body instanceof FormData) {
    body = opts.body; // Let browser set the multipart boundary
  }

  if (opts.token) {
    headers['Authorization'] = `Bearer ${opts.token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method: opts.method || 'GET',
    headers,
    body
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }

  const contentType = res.headers.get('content-type');
  if (contentType && contentType.includes('application/json')) {
    return (await res.json()) as T;
  }
  // @ts-expect-error allow empty
  return undefined as T;
}

export const endpoints = {
  health: () => apiFetch<string>('/health'),
  signup: (payload: { first_name: string; last_name: string; email: string; password1: string; password2: string; city?: string; country?: string; }) =>
    apiFetch('/api/auth/signup', { method: 'POST', body: payload }),
  login: (payload: { email: string; password: string }) =>
    apiFetch('/api/auth/login', { method: 'POST', body: payload }),
  me: (token: string) => apiFetch('/api/me', { token }),
  logout: (token: string) => apiFetch('/api/auth/logout', { method: 'POST', token }),

  videoStatuses: () => apiFetch<string[]>('/api/videos/statuses'),
  myVideos: (token: string) => apiFetch('/api/videos', { token }),
  uploadVideo: (token: string, form: FormData) => apiFetch('/api/videos/upload', { method: 'POST', token, body: form }),
  getVideo: (token: string, id: string) => apiFetch(`/api/videos/${id}`, { token }),
  deleteVideo: (token: string, id: string) => apiFetch(`/api/videos/${id}`, { method: 'DELETE', token }),

  publicVideos: () => apiFetch('/api/public/videos'),
  voteVideo: (token: string, id: string) => apiFetch(`/api/public/videos/${id}/vote`, { method: 'POST', token }),
  rankings: (params: { page?: number; pageSize?: number; city?: string }) => {
    const q = new URLSearchParams();
    if (params.page) q.set('page', String(params.page));
    if (params.pageSize) q.set('pageSize', String(params.pageSize));
    if (params.city) q.set('city', params.city);
    const qs = q.toString();
    return apiFetch(`/api/public/rankings${qs ? `?${qs}` : ''}`);
  },
  cityId: (city: string, country: string) => apiFetch(`/api/location/city-id?city=${encodeURIComponent(city)}&country=${encodeURIComponent(country)}`)
};

