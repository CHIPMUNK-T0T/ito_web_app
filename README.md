# ITO Web App

## 概要
ITO Web Appは、オンラインでITOゲームを楽しむためのWebアプリケーションです。

## プロジェクト構成
- `backend/`: Goで実装されたバックエンドサービス
- `frontend/`: Next.jsで実装されたフロントエンドアプリケーション

## 開発環境のセットアップ

### 必要条件
- Go 1.21以上
- Node.js 18以上
- MySQL 8.0以上
- Docker (開発環境用)

### バックエンド開発環境の構築

1. 依存関係のインストール

```bash
cd backend
go mod tidy
```

2. 環境変数の設定

```bash
cp .env.example .env
```

3. 開発サーバーの起動

```bash
go run backend/cmd/main.go
```


### テスト実行方法
#### 単体テストの実行
すべてのテストを実行:

```bash
cd backend
go test -v ./...
```


### テストカバレッジの確認

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### テストカバレッジレポート

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```


### テストの構成
- `test/internal/repository`: リポジトリレイヤーのテスト
- `test/mock`: モックオブジェクト

### テスト実行時の注意事項
- テスト実行前にデータベースが起動していることを確認してください
- テストデータベースの環境変数が正しく設定されていることを確認してください

## CI/CD

本プロジェクトでは以下のGitHub Actionsワークフローを使用しています：

- `go-test.yml`: テストの自動実行とカバレッジレポートの生成
- `go-lint.yml`: コードの品質チェック
- `deploy.yml`: 本番環境へのデプロイ

## プロジェクト構造

.
├── backend/
│ ├── cmd/ # アプリケーションのエントリーポイント
│ ├── internal/ # 内部パッケージ
│ │ ├── entity/ # ドメインエンティティ
│ │ ├── repository/# データアクセス層
│ │ ├── usecase/ # ユースケース層
│ │ └── handler/ # プレゼンテーション層
│ ├── test/ # テストコード
│ └── docs/ # ドキュメント
└── frontend/ # フロントエンドアプリケーション

