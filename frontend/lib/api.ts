import { STORAGE_TOKEN_KEY, getApiBase } from './constants';

const API_BASE_URL = getApiBase();

/**
 * 汎用的なfetchラッパー。認証トークンの付与などを自動で行います。
 */
async function apiFetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem(STORAGE_TOKEN_KEY); // 実装に応じて適切なストレージを使用してください
  
  const headers = new Headers(options.headers || {});
  headers.set('Content-Type', 'application/json');
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  const data = await response.json().catch(() => ({}));

  if (!response.ok) {
    throw new Error(data.error || 'API Request Failed');
  }

  return data as T;
}

export const api = {
  auth: {
    register: (data: any) => apiFetch('/auth/register', { method: 'POST', body: JSON.stringify(data) }),
    login: (data: any) => apiFetch('/auth/login', { method: 'POST', body: JSON.stringify(data) }),
  },
  rooms: {
    list: () => apiFetch('/rooms', { method: 'GET' }),
    create: (data: any) => apiFetch('/rooms', { method: 'POST', body: JSON.stringify(data) }),
    get: (id: string | number) => apiFetch(`/rooms/${id}`, { method: 'GET' }),
    join: (id: string | number, data: any) => apiFetch(`/rooms/${id}/join`, { method: 'POST', body: JSON.stringify(data) }),
  },
  games: {
    ready: (roomId: string | number) => apiFetch(`/games/${roomId}/ready`, { method: 'POST' }),
    start: (roomId: string | number, theme: string) => apiFetch(`/games/${roomId}/start`, { method: 'POST', body: JSON.stringify({ theme }) }),
    status: (roomId: string | number) => apiFetch(`/games/${roomId}/status`, { method: 'GET' }),
    refresh: (roomId: string | number) => apiFetch(`/games/${roomId}/refresh`, { method: 'POST' }),
  }
};
