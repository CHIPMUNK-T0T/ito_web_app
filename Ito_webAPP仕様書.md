# itoゲーム Webアプリケーション仕様書

## 1. システム概要

### 1.1 目的

- ボードゲーム「ito」をオンラインで遊べるWebアプリケーションとして実装
- リアルタイムでのプレイヤー間コミュニケーションを実現

### 1.2 ゲーム内容

- 参加するユーザーには1～100重複しない数字のカードが1人1枚ずつ配る
- 出されたお題に対し小さい順に出して全員のカードを出し切る
- 小さい順に出すことができなかった場合、ゲームは失敗となる
- 自分の数字を口にしたらゲーム失敗
- 出されたテーマに沿って自分のカードの数字を表現しあう

## 2. 技術スタック

### 2.1 バックエンド

- 言語：Go 1.22.2
- フレームワーク：Gin
- DB：SQLite
- ORM：GORM
- WebSocket：gorilla/websocket

### 2.2 フロントエンド

- HTML5
- CSS3
- NES.css // レトロゲーム風のデザイン
- TypeScript
- Next.js
- WebSocket

## 3. ゲームフロー

### 3.1 事前条件

1. 各ユーザーは自分の名前およびパスワードでログインすること
2. ユーザーはルームを作成することができる
   - ルームにはパスワードが必要
   - 最大参加人数設定（default:10, max:10）
   - ルーム説明文（300文字以内）
   - 参加許可制/自由参加の選択（default:自由参加）
3. ユーザーはルーム名とパスワードで入室可能
4. 2人以上の参加でゲーム開始可能

### 3.2 ゲーム開始前

1. ルーム作成者がプレイヤー数を確認
2. 各プレイヤーが「準備完了」ボタンを押す
3. 全員の準備が完了するとゲーム開始

### 3.3 ゲームの進行

1. システムが各プレイヤーにランダムに数字カード（1-100）を1枚ずつ配布
2. お題の提示（例：「朝食の値段」「理想の気温」など）
3. プレイヤーは議論を進める（音声通話は外部ツールを使用）
4. プレイヤーは「数字を公開」ボタンで公開意思表示が可能
5. 数字公開の合意形成
   - 全プレイヤーが「許可」を押す
   - または30秒経過で自動公開
   - 「まだ早い」が1名以上で却下（ペナルティなし）
6. 数値公開後のチェック
   - より小さい数値を持つプレイヤーがいる場合、自動公開してゲーム終了
   - 小さい数値がない場合、ゲーム継続
7. 最後の1名まで続き、その数値公開でゲーム終了

## 4. データベース設計

### 4.1 テーブル構成

#### 4.1.1 Users テーブル

ユーザー情報を管理するテーブル

| カラム名 | 型 | 制約 | 説明 |
|----------|-------|--------|--------|
| ID | uint | PK, AUTO_INCREMENT | ユーザーID |
| Username | string | NOT NULL, UNIQUE, size:16 | ユーザー名 |
| Password | string | NOT NULL, size:16 | ハッシュ化されたパスワード |
| CreatedAt | time.Time | NOT NULL | 作成日時 |
| UpdatedAt | time.Time | NOT NULL | 更新日時 |
| DeletedAt | gorm.DeletedAt | NULL | 削除日時 |

#### 4.1.2 Rooms テーブル

ゲームルーム情報を管理するテーブル

| カラム名 | 型 | 制約 | 説明 |
|----------|-------|--------|--------|
| ID | uint | PK, AUTO_INCREMENT | ルームID |
| Name | string | NOT NULL, size:16 | ルーム名 |
| Password | string | NOT NULL, size:16 | ハッシュ化されたルームパスワード |
| MaxPlayers | int | NOT NULL | 最大参加人数 |
| CreatorID | uint | NOT NULL, FK(Users.ID) | 作成者ユーザーID |
| Description | string | size:300 | ルーム説明文 |
| IsPrivate | bool | NOT NULL, default:false | プライベートルームフラグ |
| CreatedAt | time.Time | NOT NULL | 作成日時 |
| UpdatedAt | time.Time | NOT NULL | 更新日時 |
| DeletedAt | gorm.DeletedAt | NULL | 削除日時 |

### 4.2 ER図

```mermaid
erDiagram
    USERS ||--|{ ROOMS : "creates"

    USERS {
        uint ID PK
        string Username "unique"
        string Password
        time.Time CreatedAt
        time.Time UpdatedAt
        gorm.DeletedAt DeletedAt
    }

    ROOMS {
        uint ID PK
        string Name
        string Password
        int MaxPlayers
        uint CreatorID FK
        string Description
        bool IsPrivate
        time.Time CreatedAt
        time.Time UpdatedAt
        gorm.DeletedAt DeletedAt
    }
```

## 5. 主要機能

### 5.1 ユーザー管理機能

- ユーザー登録
  - ユーザー名
  - パスワード
- ログイン/ログアウト
- プロフィール管理
  - ユーザー情報の編集

### 5.2 ゲームルーム機能

- ルーム作成
  - パスワード設定
  - 最大参加人数設定（2-10人）
  - ルーム説明文（300文字以内）
  - 参加許可制/自由参加の選択
- ルーム設定修正
  - 最大参加人数の編集
  - ルーム説明文の編集
  - 参加許可設定の変更
- ルーム検索
  - ルーム名による検索
  - ルーム作成者による検索
- ルーム参加
  - パスワード認証
  - 参加許可制の場合の承認機能
- ルーム一覧表示
  - 現在の参加人数
  - ルームのステータス（待機中/プレイ中）

### 5.3 ゲームプレイ機能

- プレイヤー管理
  - 準備完了状態の管理
  - プレイヤーの接続状態監視
- カードシステム
  - ランダムなカード配布（1-100）
  - カード情報の秘匿管理
- ゲーム進行管理
  - お題の提示と管理
  - 数字公開の合意形成システム
  - タイマー管理（30秒自動公開）
- 結果判定システム
  - 数値の大小比較
  - ゲーム終了条件の判定
- チャットシステム
  - プレイヤー間のコミュニケーション
  - システムメッセージの表示

## 6. API エンドポイント設計

### 6.1 認証関連 (/api/auth)

| メソッド | エンドポイント | 説明 | 認証要否 |
|----------|----------------|------|----------|
| POST | /register | ユーザー登録 | 不要 |
| POST | /login | ログイン | 不要 |

### 6.2 ルーム関連 (/api/rooms)

| メソッド | エンドポイント | 説明 | 認証要否 |
|----------|----------------|------|----------|
| GET | / | ルーム一覧取得 | 必要 |
| POST | / | ルーム作成 | 必要 |
| GET | /:id | 特定ルームの情報取得 | 必要 |
| POST | /:id/join | ルーム参加 | 必要 |

### 6.3 ゲーム関連 (/api/games)

| メソッド | エンドポイント | 説明 | 認証要否 |
|----------|----------------|------|----------|
| POST | /:roomId/ready | プレイヤーの準備完了 | 必要 |
| POST | /:roomId/start | ゲーム開始（ホストのみ） | 必要 |
| POST | /:roomId/vote | 数字公開の投票開始 | 必要 |
| GET | /:roomId/status | ゲーム状態の取得 | 必要 |
| GET | /ws/:roomId | WebSocket接続 | 必要 |

### 6.4 リクエスト/レスポンス形式

#### 6.4.1 ルーム作成 (POST /api/rooms)

```json
// リクエスト
{
    "name": "ルーム名",
    "max_players": 5,
    "description": "ルーム説明",
    "is_private": true,
    "password": "ルームパスワード"  // is_private=trueの場合必須
}

// レスポンス
{
    "id": 1,
    "name": "ルーム名",
    "max_players": 5,
    "description": "ルーム説明",
    "is_private": true,
    "creator_id": 1,
    "player_count": 1
}
```

#### 6.4.2 ルーム参加 (POST /api/rooms/:id/join)

```json
// リクエスト
{
    "password": "ルームパスワード"  // プライベートルームの場合必須
}

// レスポンス
{
    "message": "ルームに参加しました",
    "room": {
        "id": 1,
        "name": "ルーム名",
        "player_count": 2,
        "players": [
            {
                "id": 1,
                "username": "ホスト",
                "role": "host"
            },
            {
                "id": 2,
                "username": "ゲスト",
                "role": "guest"
            }
        ]
    }
}
```

#### 6.4.3 WebSocketメッセージ形式

```json
// 基本形式
{
    "type": "message_type",
    "room_id": 1,
    "user_id": 1,
    "payload": {}
}

// カード配布通知
{
    "type": "card_dealt",
    "payload": {
        "card_number": 42
    }
}

// 投票開始
{
    "type": "vote_start",
    "payload": {
        "target_user_id": 1,
        "timeout_at": "2024-03-21T15:00:00Z"
    }
}

// 投票提出
{
    "type": "vote_submit",
    "payload": {
        "approve": true
    }
}
```

### 6.5 エラーレスポンス

```json
{
    "error": "エラーメッセージ"
}
```