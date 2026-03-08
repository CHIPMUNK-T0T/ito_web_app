'use client';

import { useAuth } from '@/app/providers/auth-provider';

export function ControlHeader() {
  const { token, user, logout } = useAuth();

  return (
    <header className="panel app-header" style={{ marginBottom: 0 }}>
      <div>
        <h1 style={{ fontSize: '13px', marginBottom: '4px' }}>
          ✦ ITO ゲーム ✦
        </h1>
        {user && (
          <p className="small-mono" style={{ color: 'var(--dq-cyan)' }}>
            ▶ {user.username} としてログイン中
          </p>
        )}
      </div>
      {token && (
        <button
          type="button"
          className="nes-btn is-error logout-btn"
          onClick={logout}
        >
          やめる
        </button>
      )}
    </header>
  );
}
