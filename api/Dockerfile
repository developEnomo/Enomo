# 開発用
FROM golang:1.23

WORKDIR /app

# 先に go.mod/go.sum だけ入れて依存DL
COPY go.mod go.sum ./
RUN go mod download

# 残りのコード
COPY . .

EXPOSE 8080
CMD ["go", "run", "./cmd"]
