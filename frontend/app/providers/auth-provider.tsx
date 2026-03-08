'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { STORAGE_TOKEN_KEY, STORAGE_USER_KEY } from '@/lib/constants';

type UserInfo = {
  id: number;
  username: string;
};

type AuthContextValue = {
  token: string | null;
  user: UserInfo | null;
  isReady: boolean;
  login: (token: string, user: UserInfo) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<UserInfo | null>(null);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    if (typeof window === 'undefined') {
      return;
    }
    const savedToken = window.localStorage.getItem(STORAGE_TOKEN_KEY);
    const savedUser = window.localStorage.getItem(STORAGE_USER_KEY);
    if (savedToken) {
      setToken(savedToken);
    }
    if (savedUser) {
      try {
        const parsed = JSON.parse(savedUser) as UserInfo;
        setUser(parsed);
      } catch {
        window.localStorage.removeItem(STORAGE_USER_KEY);
      }
    }
    setIsReady(true);
  }, []);

  const login = useCallback((newToken: string, userInfo: UserInfo) => {
    setToken(newToken);
    setUser(userInfo);
    if (typeof window === 'undefined') {
      return;
    }
    window.localStorage.setItem(STORAGE_TOKEN_KEY, newToken);
    window.localStorage.setItem(STORAGE_USER_KEY, JSON.stringify(userInfo));
  }, []);

  const logout = useCallback(() => {
    setToken(null);
    setUser(null);
    if (typeof window === 'undefined') {
      return;
    }
    window.localStorage.removeItem(STORAGE_TOKEN_KEY);
    window.localStorage.removeItem(STORAGE_USER_KEY);
  }, []);

  const contextValue = useMemo(
    () => ({
      token,
      user,
      isReady,
      login,
      logout
    }),
    [isReady, login, logout, token, user]
  );

  return <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('AuthProvider が必要です');
  }
  return context;
}
