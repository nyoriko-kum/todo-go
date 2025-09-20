# Goイメージをベースに
FROM golang:1.24

# ワークディレクトリ作成
WORKDIR /app/server

# go.mod だけでもコピーできるようにする
COPY server/go.mod ./
# go.sum がある場合だけコピー（無ければスキップ）
COPY server/go.sum* ./
RUN go mod download || true

# アプリコードをコピー
COPY server/ .

# ビルド
RUN go build -o main .

# コンテナ起動時に実行
CMD ["./main"]
