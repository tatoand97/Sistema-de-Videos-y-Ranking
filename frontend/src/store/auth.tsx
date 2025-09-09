import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import type { User } from '@api/types';
import { endpoints } from '@api/client';

type AuthState = {
  token: string | null;
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  register: (data: { first_name: string; last_name: string; email: string; password1: string; password2: string; city?: string; country?: string; }) => Promise<void>;
};

const AuthCtx = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'));
  const [user, setUser] = useState<User | null>(() => {
    const u = localStorage.getItem('user');
    return u ? JSON.parse(u) : null;
  });

  useEffect(() => {
    if (token) localStorage.setItem('token', token); else localStorage.removeItem('token');
  }, [token]);
  useEffect(() => {
    if (user) localStorage.setItem('user', JSON.stringify(user)); else localStorage.removeItem('user');
  }, [user]);

  const login = async (email: string, password: string) => {
    const res = await endpoints.login({ email, password });
    // @ts-ignore shape response aligns with AuthResponse
    const { token: t, user: u } = res as any;
    setToken(t);
    setUser(u);
  };

  const logout = async () => {
    if (token) {
      try { await endpoints.logout(token); } catch {}
    }
    setToken(null);
    setUser(null);
  };

  const register = async (data: { first_name: string; last_name: string; email: string; password1: string; password2: string; city?: string; country?: string; }) => {
    await endpoints.signup(data);
    await login(data.email, data.password1);
  };

  const value = useMemo(() => ({ token, user, login, logout, register }), [token, user]);
  return <AuthCtx.Provider value={value}>{children}</AuthCtx.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthCtx);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}

