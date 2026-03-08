import React from 'react';

type CardProps = {
  number: number | null;
  isRevealed: boolean;
  onPlayCard?: () => void;
  disabled?: boolean;
};

export function Card({ number, isRevealed, onPlayCard, disabled }: CardProps) {
  if (number === null) {
    return (
      <div className="panel" style={{ textAlign: 'center', opacity: 0.7 }}>
        <p>カード配布待ち...</p>
      </div>
    );
  }

  return (
    <div className={`panel ${isRevealed ? 'is-revealed' : ''}`} style={{ textAlign: 'center' }}>
      <p>あなたのカード</p>
      <div className="card-display">
        {number}
      </div>
      {!isRevealed && (
        <button 
          className={`nes-btn ${disabled ? 'is-disabled' : 'is-error'}`} 
          style={{ marginTop: '1rem', width: '100%' }}
          onClick={onPlayCard}
          disabled={disabled}
        >
          カードを出す
        </button>
      )}
      {isRevealed && (
        <p className="small-mono highlight" style={{ marginTop: '1rem' }}>
          提出済み
        </p>
      )}
    </div>
  );
}
