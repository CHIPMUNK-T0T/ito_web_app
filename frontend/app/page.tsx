'use client';

import { FormEvent, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getApiBase } from '@/lib/constants';
import { useAuth } from '@/app/providers/auth-provider';
import { useRedirectIfAuthenticated } from '@/app/hooks/use-auth-guards';

type AuthMode = 'login' | 'register';

export default function LoginPage() {
  useRedirectIfAuthenticated();
  const { login } = useAuth();
  const router = useRouter();
  const [authMode, setAuthMode] = useState<AuthMode>('login');
  const [form, setForm] = useState({ username: '', password: '' });
  const [error, setError] = useState<string | null>(null);
  const [info, setInfo] = useState<string | null>(null);
  const [isBusy, setIsBusy] = useState(false);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setInfo(null);
    if (!form.username.trim() || !form.password) {
      setError('なまえと パスワードを にゅうりょくしてください。');
      return;
    }
    // 登録モードのクライアントサイドバリデーション
    if (authMode === 'register') {
      if (form.username.trim().length < 3 || form.username.trim().length > 16) {
        setError('なまえは 3もじ いじょう 16もじ いかで にゅうりょく。');
        return;
      }
      if (form.password.length < 8 || form.password.length > 16) {
        setError('パスワードは 8もじ いじょう 16もじ いかで にゅうりょく。');
        return;
      }
    }
    setIsBusy(true);
    try {
      const apiBase = getApiBase();
      const response = await fetch(`${apiBase}/auth/${authMode}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: form.username.trim(),
          password: form.password
        })
      });
      const payload = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(payload?.error ?? 'にんしょうエラーが はっせいしました。');
      }
      if (!payload?.token || !payload?.user) {
        throw new Error('トークンじょうほうが とりえませんでした。');
      }
      login(payload.token, payload.user);
      setInfo(authMode === 'login' ? 'ログインしました！' : 'とうろくしました！');
      router.push('/rooms');
    } catch (err) {
      console.error(err);
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('にんしょうに しっぱいしました。');
      }
    } finally {
      setIsBusy(false);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '24px' }}>
      {/* タイトルロゴ */}
      <div style={{ textAlign: 'center', padding: '16px 0' }}>
        <h1 style={{
          fontSize: '20px',
          color: 'var(--dq-yellow)',
          textShadow: '3px 3px 0 #000, 0 0 12px rgba(248,216,0,0.5)',
          letterSpacing: '0.2em',
          marginBottom: '8px'
        }}>
          ✦ ITO ✦
        </h1>
        <p className="small-mono blink-text" style={{ fontSize: '8px' }}>
          PRESS START
        </p>
      </div>

      {/* ログイン / 登録 ウィンドウ */}
      <article className="panel" style={{ width: '100%', maxWidth: '440px' }}>
        <h2 style={{ marginBottom: '16px' }}>
          {authMode === 'login' ? '▶ ログイン' : '▶ あたらしく はじめる'}
        </h2>

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>

          {/* ユーザー名 */}
          <div className="nes-field">
            <label htmlFor="username">なまえ</label>
            <input
              id="username"
              className="nes-input"
              value={form.username}
              autoComplete="username"
              onChange={(e) => setForm((prev) => ({ ...prev, username: e.target.value }))}
            />
          </div>

          {/* パスワード */}
          <div className="nes-field">
            <label htmlFor="password">パスワード</label>
            <input
              id="password"
              className="nes-input"
              type="password"
              value={form.password}
              autoComplete={authMode === 'login' ? 'current-password' : 'new-password'}
              onChange={(e) => setForm((prev) => ({ ...prev, password: e.target.value }))}
            />
            {authMode === 'register' && (
              <p className="small-mono" style={{ color: 'var(--dq-gray)', marginTop: '4px', fontSize: '8px' }}>
                8〜16もじ
              </p>
            )}
          </div>

          {/* エラー表示 */}
          {error && (
            <div style={{
              background: 'rgba(248,56,0,0.15)',
              border: '2px solid var(--dq-red)',
              padding: '8px 12px',
            }}>
              <p className="small-mono" style={{ color: 'var(--dq-red)' }}>
                ！ {error}
              </p>
            </div>
          )}

          {/* 成功表示 */}
          {info && (
            <p className="small-mono" style={{ color: 'var(--dq-green)' }}>✔ {info}</p>
          )}

          {/* 送信ボタン */}
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <button
              className={`nes-btn ${authMode === 'login' ? 'is-primary' : 'is-warning'}`}
              type="submit"
              disabled={isBusy}
              style={{ minWidth: '140px' }}
            >
              {isBusy ? '…しょり中' : authMode === 'login' ? 'ログイン' : 'はじめる'}
            </button>
          </div>
        </form>

        {/* モード切り替え */}
        <div style={{
          marginTop: '16px',
          paddingTop: '12px',
          borderTop: '2px solid var(--dq-win-in)',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          gap: '8px',
          flexWrap: 'wrap'
        }}>
          <p className="small-mono" style={{ color: 'var(--dq-gray)', fontSize: '8px' }}>
            {authMode === 'login' ? 'はじめて ですか？' : 'すでに アカウントが ありますか？'}
          </p>
          <button
            type="button"
            className="nes-btn"
            style={{
              background: 'transparent',
              color: 'var(--dq-cyan)',
              border: '2px solid var(--dq-cyan)',
              fontSize: '8px',
              padding: '6px 10px'
            }}
            onClick={() => {
              setAuthMode(authMode === 'login' ? 'register' : 'login');
              setError(null);
              setInfo(null);
            }}
          >
            {authMode === 'login' ? 'あたらしく はじめる' : 'ログインへ'}
          </button>
        </div>
      </article>
    </div>
  );
}
