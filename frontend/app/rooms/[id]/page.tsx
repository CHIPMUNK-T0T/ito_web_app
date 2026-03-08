'use client';

import { useEffect, useState, useCallback, useRef } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useAuth } from '@/app/providers/auth-provider';
import { useRequireAuth } from '@/app/hooks/use-auth-guards';
import { api } from '@/lib/api';
import { useGameWebSocket } from '@/hooks/useGameWebSocket';
import { Card } from '@/components/Card';
import { ChatBox, ChatMessage } from '@/components/ChatBox';
import { PlayerList } from '@/components/PlayerList';

type RoomData = {
  id: number;
  name: string;
  creator_id: number;
  players: any[];
};

type GameStatusData = {
  status: string; // 'waiting', 'playing', 'finished'
  theme: string;
};

export default function GameRoomPage() {
  useRequireAuth();
  const { id } = useParams();
  const router = useRouter();
  const { user, token } = useAuth();
  
  const [room, setRoom] = useState<RoomData | null>(null);
  const [gameStatus, setGameStatus] = useState<string>('waiting');
  const [theme, setTheme] = useState<string>('');
  const [myCard, setMyCard] = useState<number | null>(null);
  const [players, setPlayers] = useState<any[]>([]);
  const [deckInfo, setDeckInfo] = useState<{remaining: number, total: number} | null>(null);
  
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [resultEvent, setResultEvent] = useState<{status: string, message: string} | null>(null);
  
  const { isConnected, lastMessage, error, sendMessage } = useGameWebSocket(id as string, token);

  const fetchRoomInfo = useCallback(async () => {
    try {
      const roomData = await api.rooms.get(id as string);
      setRoom(roomData as RoomData);
      setPlayers((roomData as RoomData).players);
      
      const statusData = await api.games.status(id as string) as GameStatusData;
      setGameStatus(statusData.status);
      setTheme(statusData.theme);
    } catch (err) {
      console.error('Failed to fetch room or game status:', err);
    }
  }, [id]);

  useEffect(() => {
    fetchRoomInfo();
    // ポーリングして入室者の変更などを一応拾う（WebSocketでも同期可能だが簡便のため）
    const interval = setInterval(fetchRoomInfo, 5000);
    return () => clearInterval(interval);
  }, [fetchRoomInfo]);

  // WebSocketからのメッセージ処理
  useEffect(() => {
    if (!lastMessage) return;
    
    switch (lastMessage.type) {
      case 'card_dealt':
        setMyCard(lastMessage.payload.card_number);
        setPlayers(prev => prev.map(p => p.id === user?.id ? { ...p, cardNumber: lastMessage.payload.card_number, isReady: false } : p));
        break;
      case 'game_start':
        setGameStatus('playing');
        setTheme(lastMessage.payload.theme);
        setResultEvent(null);
        // setMyCard(null); // 削除: card_dealt を上書きしてしまうため
        // 全員の準備情報などをリセット
        setPlayers(prev => prev.map(p => ({ ...p, isReady: false, cardNumber: undefined })));
        break;
      case 'play_card':
        // 誰かがカードを提出した
        setPlayers(prev => prev.map(p => {
          if (p.id === lastMessage.payload.user_id) {
            return { ...p, isReady: true, cardNumber: lastMessage.payload.card_number };
          }
          return p;
        }));
        break;
      case 'game_result':
        setGameStatus('finished');
        setResultEvent({
          status: lastMessage.payload.status,
          message: lastMessage.payload.message
        });
        break;
      case 'deck_info':
        setDeckInfo({
          remaining: lastMessage.payload.remaining_count,
          total: lastMessage.payload.total_count
        });
        break;
      case 'chat_message':
        setChatMessages(prev => [...prev, lastMessage.payload as ChatMessage]);
        break;
      case 'connected':
        // WebSocket接続確認
        break;
    }
  }, [lastMessage, user?.id]);

  const handleReady = async () => {
    try {
      await api.games.ready(id as string);
      // fetchRoomInfoがポーリングで拾うか、即時反映のため更新
      setPlayers(prev => prev.map(p => p.id === user?.id ? { ...p, isReady: true } : p));
    } catch (err) {
      alert(`エラー: ${(err as Error).message}`);
    }
  };

  const handleStartGame = async () => {
    try {
      await api.games.start(id as string, theme);
    } catch (err) {
      alert(`開始エラー: ${(err as Error).message}`);
    }
  };

  const handlePlayCard = () => {
    if (myCard === null) return;
    sendMessage('play_card', { user_id: user?.id, card_number: myCard });
  };

  const handleSendMessage = (msg: string) => {
    sendMessage('chat_message', { message: msg });
  };

  if (!room) {
    return <div style={{ textAlign: 'center', marginTop: '4rem' }}><p>読み込み中...</p></div>;
  }

  const isHost = room.creator_id === user?.id;
  const me = players.find(p => p.id === user?.id);
  const myStatusReady = me?.isReady || false;
  const allReady = players.every(p => p.isReady);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem', paddingBottom: '2rem' }}>
      
      {/* RESULT OVERLAY */}
      {resultEvent && (
        <div style={{
          position: 'fixed', inset: 0, zIndex: 9999, 
          backgroundColor: resultEvent.status === 'success' ? 'rgba(255, 223, 0, 0.85)' : 'rgba(139, 0, 0, 0.9)',
          display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center',
          color: resultEvent.status === 'success' ? '#000' : '#fff',
          animation: 'fadeIn 0.5s ease-out'
        }}>
          <div className="panel" style={{ 
            backgroundColor: resultEvent.status === 'success' ? '#fff9c4' : '#2c0000', 
            borderColor: resultEvent.status === 'success' ? '#fbc02d' : '#ff5252',
            maxWidth: '600px', width: '90%', textAlign: 'center'
          }}>
            <h1 style={{ color: resultEvent.status === 'success' ? '#f57f17' : '#ff5252', fontSize: '2.5rem', marginBottom: '1rem' }}>
              {resultEvent.status === 'success' ? '🎊 MISSION CLEAR 🎊' : '💀 GAME OVER 💀'}
            </h1>
            <p style={{ fontSize: '1.2rem', lineHeight: '1.6' }}>{resultEvent.message}</p>
            {isHost && (
              <div style={{ marginTop: '2rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                <div className="nes-field">
                  <label htmlFor="next_theme">次のお題（空欄でランダム）</label>
                  <input
                    id="next_theme"
                    className="nes-input"
                    value={theme}
                    onChange={(e) => setTheme(e.target.value)}
                    placeholder="例: 朝食の値段"
                    style={{ backgroundColor: '#fff', color: '#000' }}
                  />
                </div>
                <button 
                  className={`nes-btn ${resultEvent.status === 'success' ? 'is-warning' : 'is-error'}`} 
                  onClick={() => {
                    setResultEvent(null);
                    handleStartGame();
                  }}
                >
                  次のゲームを始める（再配り）
                </button>
              </div>
            )}
            {!isHost && (
              <p className="small-mono" style={{ marginTop: '2rem', opacity: 0.8 }}>
                ホストが次のゲームを開始するまでお待ちください...
              </p>
            )}
            <button className="nes-btn" style={{ marginTop: '1rem' }} onClick={() => setResultEvent(null)}>
              閉じる
            </button>
          </div>
        </div>
      )}

      <header className="panel panel--slim" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h2>Room: {room.name}</h2>
          <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
            <span className="small-mono" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              ステータス: <span className="badge">{gameStatus.toUpperCase()}</span>
              接続: <span style={{ color: isConnected ? '#92cc41' : '#ff5a5f' }}>{isConnected ? 'ONLINE' : 'OFFLINE'}</span>
            </span>
            {deckInfo && (
              <span className="small-mono" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                山札残数: <span className="nes-text is-primary">{deckInfo.remaining} / {deckInfo.total}</span>
              </span>
            )}
            {isHost && (
              <button 
                className="nes-btn is-error" 
                style={{ padding: '2px 8px', fontSize: '0.7rem' }}
                onClick={async () => {
                  if (confirm('山札をリセット（1-100に戻す）しますか？')) {
                    try { await api.games.refresh(id as string); } catch(e) { alert(e); }
                  }
                }}
              >
                山札リセット
              </button>
            )}
          </div>
        </div>
        <button className="nes-btn" onClick={() => router.push('/rooms')}>退出</button>
      </header>

      {gameStatus === 'playing' && theme && (
        <section className="panel" style={{ backgroundColor: '#fff9c4', borderColor: '#fbc02d', textAlign: 'center' }}>
          <h3 style={{ color: '#f57f17', margin: 0 }}>今回のお題</h3>
          <p style={{ fontSize: '1.5rem', fontWeight: 'bold', margin: '0.5rem 0' }}>「{theme}」</p>
          <p className="small-mono">話し合って、数字が小さい自信のある人からカードを出してください</p>
        </section>
      )}

      <div className="section-grid" style={{ gridTemplateColumns: '1fr 320px' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          
          {gameStatus === 'waiting' ? (
            <div className="panel retro-screen" style={{ minHeight: '350px', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', gap: '1rem' }}>
              <i className="nes-logo" style={{ transform: 'scale(1.2)' }}></i>
              <p style={{ fontSize: '1rem', textAlign: 'center' }}>
                プレイヤーの参加と準備を待っています...<br/>
                <span className="small-mono">({players.filter(p => p.isReady).length} / {players.length} 準備完了)</span>
              </p>

              {isHost && (
                <div className="nes-field" style={{ width: '80%' }}>
                  <label htmlFor="theme_input">お題（空欄でランダム）</label>
                  <input
                    id="theme_input"
                    className="nes-input"
                    value={theme}
                    onChange={(e) => setTheme(e.target.value)}
                    placeholder="例: 理想の気温"
                  />
                </div>
              )}
              
              <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', justifyContent: 'center' }}>
                {!myStatusReady && (
                  <button className="nes-btn is-success" onClick={handleReady}>準備完了にする</button>
                )}
                {isHost && (
                  <button 
                    className={`nes-btn ${players.length >= 2 ? 'is-primary' : 'is-disabled'}`} 
                    onClick={handleStartGame}
                    disabled={players.length < 2}
                  >
                    ゲーム開始
                  </button>
                )}
              </div>
            </div>
          ) : (
            <Card 
              number={myCard} 
              isRevealed={me?.isReady || false} // プレイ中、isReady=true は提出済みを意味する
              disabled={gameStatus !== 'playing' || myCard === null || (me?.isReady || false)}
              onPlayCard={handlePlayCard}
            />
          )}

          <PlayerList players={players} currentUserId={user?.id || null} hostId={room.creator_id} />

        </div>

        <aside>
          <ChatBox 
            messages={chatMessages} 
            onSendMessage={handleSendMessage} 
            currentUser={user?.username || null} 
            disabled={!isConnected}
          />
        </aside>
      </div>

    </div>
  );
}
