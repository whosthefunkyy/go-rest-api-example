FROM golang:1.23.3

WORKDIR /app

# Копируем сначала только модули для кеширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

RUN go build -o main ./cmd/api/

CMD ["./main"]
