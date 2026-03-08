'use client';

import { FormEvent, useCallback, useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getApiBase } from '@/lib/constants';
import { useAuth } from '@/app/providers/auth-provider';
import { useRequireAuth } from '@/app/hooks/use-auth-guards';

type RoomInfo = {
  id: number;
  name: string;
  max_players: number;
  description?: string;
  creator_id?: number;
  player_count?: number;
  is_private?: boolean;
  status?: string;
};

export default function RoomsPage() {
  useRequireAuth();
  const router = useRouter();
  const { token } = useAuth();
  const [rooms, setRooms] = useState<RoomInfo[]>([]);
  const [message, setMessage] = useState<string | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [createForm, setCreateForm] = useState({
    name: '',
    maxPlayers: 4,
    description: '',
    isPrivate: false,
    password: ''
  });
  const [joinForm, setJoinForm] = useState({ roomId: '', password: '' });
  const [busyAction, setBusyAction] = useState<'create' | 'join' | null>(null);

  const authHeaders = useMemo((): Record<string, string> => {
    if (!token) return {};
    return {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`
    };
  }, [token]);

  const fetchRooms = useCallback(async () => {
    if (!token) return;
    try {
      setIsLoading(true);
      setIsError(false);
      const apiBase = getApiBase();
      const response = await fetch(`${apiBase}/rooms`, { headers: authHeaders });
      if (!response.ok) {
        const payload = await response.json().catch(() => null);
        throw new Error(payload?.error ?? 'ルームいちらんの しゅとくに しっぱいしました。');
      }
      const data = (await response.json()) as RoomInfo[];
      setRooms(data);
      setMessage(`${data.length}けの ルームが みつかりました。`);
    } catch (error) {
      setIsError(true);
      setMessage(error instanceof Error ? error.message : 'しゅとくに しっぱいしました。');
    } finally {
      setIsLoading(false);
    }
  }, [authHeaders, token]);

  useEffect(() => {
    fetchRooms();
  }, [fetchRooms]);

  const handleCreateRoom = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!token) { setMessage('ログインが ひつようです。'); setIsError(true); return; }
    if (!createForm.name.trim()) { setMessage('ルームめいを にゅうりょくしてください。'); setIsError(true); return; }
    if (createForm.name.trim().length < 3) { setMessage('ルームめいは 3もじ いじょうで にゅうりょく。'); setIsError(true); return; }
    if (createForm.isPrivate && !createForm.password.trim()) {
      setMessage('プライベートルームには パスワードが ひつようです。');
      setIsError(true);
      return;
    }
    setBusyAction('create');
    setIsError(false);
    try {
      const apiBase = getApiBase();
      const response = await fetch(`${apiBase}/rooms`, {
        method: 'POST',
        headers: authHeaders,
        body: JSON.stringify({
          name: createForm.name.trim(),
          max_players: createForm.maxPlayers,
          description: createForm.description.trim() || undefined,
          is_private: createForm.isPrivate,
          password: createForm.isPrivate ? createForm.password : undefined
        })
      });
      const payload = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(payload?.error ?? 'ルームさくせいに しっぱいしました。');
      }
      setMessage('ルームを さくせいし さんかしました！');
      router.push(`/rooms/${payload.id}`);
    } catch (error) {
      setIsError(true);
      setMessage(error instanceof Error ? error.message : 'ルームさくせいに しっぱいしました。');
    } finally {
      setBusyAction(null);
    }
  };

  const handleJoinRoom = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!token) { setMessage('ログインが ひつようです。'); setIsError(true); return; }
    const targetId = Number(joinForm.roomId);
    if (!targetId) { setMessage('ゆうこうな ルームIDを にゅうりょくしてください。'); setIsError(true); return; }
    setBusyAction('join');
    setIsError(false);
    try {
      const apiBase = getApiBase();
      const response = await fetch(`${apiBase}/rooms/${targetId}/join`, {
        method: 'POST',
        headers: authHeaders,
        body: JSON.stringify({ password: joinForm.password.trim() })
      });
      const payload = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(payload?.error ?? 'ルームさんかに しっぱいしました。');
      }
      setMessage('ルームに さんかしました！');
      router.push(`/rooms/${payload?.room?.id ?? targetId}`);
    } catch (error) {
      setIsError(true);
      setMessage(error instanceof Error ? error.message : 'ルームさんかに しっぱいしました。');
    } finally {
      setBusyAction(null);
    }
  };

  return (
    <div className="rooms-section">

      {/* メッセージバナー */}
      {message && (
        <div style={{
          background: isError ? 'rgba(248,56,0,0.15)' : 'rgba(56,208,64,0.1)',
          border: `2px solid ${isError ? 'var(--dq-red)' : 'var(--dq-green)'}`,
          padding: '8px 12px',
        }}>
          <p className="small-mono" style={{ color: isError ? 'var(--dq-red)' : 'var(--dq-green)' }}>
            {isError ? '！ ' : '✔ '}{message}
          </p>
        </div>
      )}

      {/* 上段: 参加 / 作成 */}
      <div className="rooms-grid-top">

        {/* ルーム参加ウィンドウ */}
        <article className="panel panel--slim panel--stretch">
          <h2>▶ ルームに さんか</h2>
          <form onSubmit={handleJoinRoom} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div className="nes-field">
              <label htmlFor="targetRoom">ルームID</label>
              <input
                id="targetRoom"
                className="nes-input"
                type="number"
                value={joinForm.roomId}
                onChange={(e) => setJoinForm((prev) => ({ ...prev, roomId: e.target.value }))}
              />
            </div>
            <div className="nes-field">
              <label htmlFor="joinPassword">パスワード（プライベートのみ）</label>
              <input
                id="joinPassword"
                className="nes-input"
                type="password"
                value={joinForm.password}
                onChange={(e) => setJoinForm((prev) => ({ ...prev, password: e.target.value }))}
              />
            </div>
            <div className="create-actions">
              <button className="nes-btn is-primary" type="submit" disabled={busyAction === 'join'}>
                {busyAction === 'join' ? '…さんか中' : 'さんか'}
              </button>
            </div>
          </form>
        </article>

        {/* ルーム作成ウィンドウ */}
        <article className="panel panel--slim panel--stretch">
          <h2>▶ ルームを つくる</h2>
          <form onSubmit={handleCreateRoom} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div className="nes-field">
              <label htmlFor="roomName">ルームめい（3〜16もじ）</label>
              <input
                id="roomName"
                className="nes-input"
                value={createForm.name}
                onChange={(e) => setCreateForm((prev) => ({ ...prev, name: e.target.value }))}
              />
            </div>
            <div className="nes-field">
              <label htmlFor="maxPlayers">さいだい にんずう</label>
              <input
                id="maxPlayers"
                className="nes-input"
                type="number"
                min={2}
                max={10}
                value={createForm.maxPlayers}
                onChange={(e) => setCreateForm((prev) => ({ ...prev, maxPlayers: Number(e.target.value) }))}
              />
            </div>
            <div className="nes-field">
              <label htmlFor="roomDescription">せつめい（にんい）</label>
              <textarea
                id="roomDescription"
                className="nes-textarea"
                rows={2}
                value={createForm.description}
                onChange={(e) => setCreateForm((prev) => ({ ...prev, description: e.target.value }))}
              />
            </div>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
              <input
                type="checkbox"
                className="nes-checkbox"
                checked={createForm.isPrivate}
                onChange={(e) => setCreateForm((prev) => ({ ...prev, isPrivate: e.target.checked }))}
              />
              <span>プライベートルーム</span>
            </label>
            {createForm.isPrivate && (
              <div className="nes-field">
                <label htmlFor="roomPassword">ルームパスワード</label>
                <input
                  id="roomPassword"
                  className="nes-input"
                  type="password"
                  value={createForm.password}
                  onChange={(e) => setCreateForm((prev) => ({ ...prev, password: e.target.value }))}
                />
              </div>
            )}
            <div className="create-actions">
              <button className="nes-btn is-warning" type="submit" disabled={busyAction === 'create'}>
                {busyAction === 'create' ? '…さくせい中' : 'ルームを つくる'}
              </button>
            </div>
          </form>
        </article>
      </div>

      {/* ルーム一覧 */}
      <article className="panel rooms-list-panel">
        <div className="split-row" style={{ justifyContent: 'space-between', marginBottom: '12px' }}>
          <h2>▶ ルームいちらん</h2>
          <button
            className="nes-btn is-primary"
            type="button"
            onClick={fetchRooms}
            disabled={isLoading}
            style={{ fontSize: '8px', padding: '6px 10px' }}
          >
            {isLoading ? '…こうしん中' : '↻ こうしん'}
          </button>
        </div>

        <div className="log-box">
          {isLoading ? (
            <p className="small-mono" style={{ color: 'var(--dq-cyan)' }}>▶ よみこみ中…</p>
          ) : rooms.length ? (
            <ul style={{ listStyle: 'none', display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {rooms.map((room) => (
                <li key={room.id} style={{
                  borderBottom: '1px solid var(--dq-win-in)',
                  paddingBottom: '8px',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'flex-start',
                  gap: '8px',
                  flexWrap: 'wrap'
                }}>
                  <div>
                    <p className="small-mono" style={{ color: 'var(--dq-yellow)' }}>
                      ▶ {room.name}
                    </p>
                    {room.description && (
                      <p className="small-mono" style={{ color: 'var(--dq-gray)', fontSize: '8px' }}>
                        {room.description}
                      </p>
                    )}
                  </div>
                  <div style={{ textAlign: 'right', display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: '8px' }}>
                    <div>
                      <span className="small-mono" style={{ color: 'var(--dq-cyan)', fontSize: '8px' }}>
                        ID:{room.id} ／ {room.player_count ?? 0}/{room.max_players}にん
                      </span>
                      {room.is_private && (
                        <p className="small-mono" style={{ color: 'var(--dq-red)', fontSize: '8px' }}>
                          🔒 PRIVATE
                        </p>
                      )}
                    </div>
                    <button 
                      className="nes-btn is-primary" 
                      style={{ fontSize: '8px', padding: '4px 8px' }}
                      onClick={() => setJoinForm(prev => ({ ...prev, roomId: String(room.id) }))}
                    >
                      ▶ 参加
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          ) : (
            <p className="small-mono" style={{ color: 'var(--dq-gray)' }}>
              ルームが みつかりませんでした。
            </p>
          )}
        </div>

        <p className="small-mono" style={{ marginTop: '8px', color: 'var(--dq-gray)', fontSize: '8px' }}>
          ※ プライベートルームは このいちらんには ひょうじされません。
        </p>
      </article>
    </div>
  );
}
