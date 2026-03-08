import React, { useState, useRef, useEffect } from 'react';

export type ChatMessage = {
  username: string;
  message: string;
  sentAt: number;
};

type ChatBoxProps = {
  messages: ChatMessage[];
  onSendMessage: (message: string) => void;
  currentUser: string | null;
  disabled?: boolean;
};

export function ChatBox({ messages, onSendMessage, currentUser, disabled }: ChatBoxProps) {
  const [inputText, setInputText] = useState('');
  const logBoxRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (logBoxRef.current) {
      logBoxRef.current.scrollTop = logBoxRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (inputText.trim() && !disabled) {
      onSendMessage(inputText);
      setInputText('');
    }
  };

  return (
    <div className="panel" style={{ display: 'flex', flexDirection: 'column', height: '400px' }}>
      <h3 style={{ marginBottom: '1rem', borderBottom: '2px solid', paddingBottom: '0.5rem' }}>Chat Log</h3>
      
      <div 
        ref={logBoxRef} 
        style={{ flex: 1, overflowY: 'auto', marginBottom: '1rem', padding: '0.5rem' }}
      >
        {messages.length === 0 ? (
          <p className="small-mono" style={{ textAlign: 'center', opacity: 0.6, marginTop: '2rem' }}>
            No messages yet.
          </p>
        ) : (
          <section className="message-list">
            {messages.map((msg, idx) => {
              const isMe = msg.username === currentUser;
              return (
                <section 
                  key={idx} 
                  className={`message ${isMe ? '-right' : '-left'}`}
                  style={{ display: 'flex', flexDirection: isMe ? 'row-reverse' : 'row', marginBottom: '1rem' }}
                >
                  <div style={{ maxWidth: '75%' }}>
                    <div className="small-mono" style={{ marginBottom: '0.2rem', textAlign: isMe ? 'right' : 'left' }}>
                      {msg.username}
                    </div>
                    <div className={`nes-balloon from-${isMe ? 'right' : 'left'} ${isMe ? 'is-dark' : ''}`} style={{ padding: '0.5rem 1rem' }}>
                      <p style={{ margin: 0, wordBreak: 'break-word' }}>{msg.message}</p>
                    </div>
                  </div>
                </section>
              );
            })}
          </section>
        )}
      </div>

      <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '0.5rem' }}>
        <input 
          type="text" 
          className="nes-input" 
          placeholder="Message..." 
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          disabled={disabled}
          style={{ flex: 1 }}
        />
        <button 
          type="submit" 
          className={`nes-btn ${disabled || !inputText.trim() ? 'is-disabled' : 'is-primary'}`}
          disabled={disabled || !inputText.trim()}
        >
          Send
        </button>
      </form>
    </div>
  );
}
