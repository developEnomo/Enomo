FROM golang:1.22

WORKDIR /app
COPY . .
RUN go mod tidy

EXPOSE 8080
CMD ["go", "run", "."]