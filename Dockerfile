# Build stage (временная, тяжелая).
FROM golang:1.25.0-alpine AS builder

# Устанавливаем рабочую директорию для всего, что будет дальше на этой стадии.
WORKDIR /app

# Копируем файлы в рабочую директорию, они меняются редко.
COPY go.mod ./

# Подгружает зависимости, флаг -х показывает подробную инфу о процессе загрузки, что помогает в отладке.
RUN go mod download -x

# Благодаря файлу .dockerignore, копируем нужные файлы.
COPY . .

# Собираем приложение, CGO - отключаем, на С не пишем и независим от системных библиотек.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Final stage (легкий).
FROM alpine:3.22

# Создаем непривилегированного пользователя (для безопасности).
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Рабочая директория для запуска приложения.
WORKDIR /app

# Копируем бинарник из стадии builder.
COPY --from=builder /app/main .

# Копируем статические файлы и шаблон.
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Переключаемся на непривилегированного пользователя.
USER appuser

EXPOSE 8080

CMD ["./main"]