import React from 'react';

type Player = {
  id: number;
  username: string;
  isReady: boolean;
  cardNumber?: number | null;
  isConnected: boolean;
};

type PlayerListProps = {
  players: Player[];
  currentUserId: number | null;
  hostId: number | null;
};

export function PlayerList({ players, currentUserId, hostId }: PlayerListProps) {
  return (
    <div className="panel" style={{ padding: '1.25rem' }}>
      <h3 style={{ marginBottom: '1rem', borderBottom: '2px solid', paddingBottom: '0.5rem' }}>
        参加メンバー ({players.length}人)
      </h3>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
        {players.map((p) => {
          const isMe = p.id === currentUserId;
          const isHost = p.id === hostId;
          return (
            <div 
              key={p.id} 
              className="player-row"
              style={{ padding: '0.5rem', backgroundColor: isMe ? 'rgba(255, 235, 59, 0.2)' : 'transparent', borderRadius: '4px' }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span style={{ 
                  display: 'inline-block', 
                  width: '10px', 
                  height: '10px', 
                  borderRadius: '50%', 
                  backgroundColor: p.isConnected ? '#92cc41' : '#ff5a5f' 
                }} title={p.isConnected ? 'オンライン' : 'オフライン'} />
                <span className="small-mono" style={{ fontWeight: isMe ? 'bold' : 'normal' }}>
                  {p.username}
                  {isMe && ' (あなた)'}
                  {isHost && ' 👑'}
                </span>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                {p.cardNumber !== undefined && p.cardNumber !== null && p.isReady && (
                  <span className="nes-badge">
                    <span className="is-success" style={{ fontWeight: 'bold' }}>{p.cardNumber}</span>
                  </span>
                )}
                {p.cardNumber !== undefined && p.cardNumber !== null && !p.isReady && (
                  <span className="nes-badge">
                    <span className="is-warning">㊙️</span>
                  </span>
                )}
                <span className="status-pill small-mono" style={{ borderColor: p.isReady ? '#92cc41' : '#a0a0a0', color: p.isReady ? '#92cc41' : '#a0a0a0' }}>
                  {p.cardNumber !== undefined ? (p.isReady ? '提出済' : '手札あり') : (p.isReady ? '準備完了' : '待機中')}
                </span>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  );
}
