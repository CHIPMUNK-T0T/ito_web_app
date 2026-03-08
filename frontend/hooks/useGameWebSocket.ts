import { useEffect, useRef, useState, useCallback } from 'react';
import { getWsBase } from '@/lib/constants';

const WEBSOCKET_URL = getWsBase();

export type GameMessage = {
  type: string;
  room_id: number;
  user_id: number;
  payload: any;
};

export function useGameWebSocket(roomId: string | number, token: string | null) {
  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<GameMessage | null>(null);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  const connect = useCallback(() => {
    if (!token) return;
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) return;

    // ブラウザのWebSocket APIはAuthorizationヘッダーを設定できないため
    // クエリパラメータ ?token=<jwt> でバックエンドのAuthMiddlewareWSに認証させる
    const url = `${WEBSOCKET_URL}/${roomId}?token=${encodeURIComponent(token)}`;
    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      setError(null);
    };

    ws.onmessage = (event) => {
      try {
        const msg: GameMessage = JSON.parse(event.data);
        setLastMessage(msg);
      } catch (err) {
        console.error('メッセージパースエラー:', err);
      }
    };

    ws.onerror = (event) => {
      console.error('WebSocket Error:', event);
      setError('接続エラーが発生しました');
    };

    ws.onclose = () => {
      setIsConnected(false);
      // 再接続ロジックなどを必要に応じて追加
    };

  }, [roomId, token]);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setIsConnected(false);
  }, []);

  const sendMessage = useCallback((type: string, payload: any = {}) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type, payload }));
    } else {
      console.warn('WebSocketが接続されていません:', type);
    }
  }, []);

  useEffect(() => {
    connect();
    return () => disconnect();
  }, [connect, disconnect]);

  return { isConnected, lastMessage, error, sendMessage, reconnect: connect };
}
