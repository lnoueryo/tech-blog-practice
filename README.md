起動時の処理
1.環境の確認
2.テンプレートの読み込み、キャッシュ


デプロイ
gcloud run deploy --source .

.envファイルと.env.devの作成

GITHUB_CLIENT_ID=
GITHUB_SECRET_ID=
APP_ENV=local
APP_HOST=localhost
DB_NAME=practices
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=""
DB_PORT=3306
DB_QUERY=parseTime=true