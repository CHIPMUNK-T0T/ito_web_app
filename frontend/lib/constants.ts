const resolveApiBase = () => {
  if (process.env.NEXT_PUBLIC_API_BASE_URL) {
    return process.env.NEXT_PUBLIC_API_BASE_URL;
  }
  if (typeof window !== 'undefined') {
    try {
      const url = new URL(window.location.href);
      const port = url.port ? Number(url.port) : url.protocol === 'https:' ? 443 : 80;
      const apiPort = port === 3000 ? 8080 : port;
      url.port = String(apiPort);
      url.pathname = '/api';
      url.search = '';
      url.hash = '';
      return url.toString().replace(/\/$/, '');
    } catch {
      return 'http://localhost:8080/api';
    }
  }
  return 'http://localhost:8080/api';
};

const resolveWsBase = () => {
  if (process.env.NEXT_PUBLIC_WS_BASE_URL) {
    return process.env.NEXT_PUBLIC_WS_BASE_URL;
  }
  if (typeof window !== 'undefined') {
    try {
      const url = new URL(window.location.href);
      const isHttps = url.protocol === 'https:';
      const port = url.port ? Number(url.port) : isHttps ? 443 : 80;
      const apiPort = port === 3000 ? 8080 : port;
      url.protocol = isHttps ? 'wss:' : 'ws:';
      url.port = String(apiPort);
      url.pathname = '/api/ws';
      url.search = '';
      url.hash = '';
      return url.toString().replace(/\/$/, '');
    } catch {
      return 'ws://localhost:8080/api/ws';
    }
  }
  return 'ws://localhost:8080/api/ws';
};

export const getApiBase = () => resolveApiBase();
export const getWsBase = () => resolveWsBase();
export const STORAGE_TOKEN_KEY = 'ito_token';
export const STORAGE_USER_KEY = 'ito_user';
