# Используем официальный образ Go
FROM golang:1.20

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Загружаем зависимости и компилируем проект
RUN go mod tidy
RUN go build -o main .

# Указываем переменные окружения для Docker
ENV DB_HOST=db \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=123456 \
    DB_NAME=mydb

# Запускаем приложение
CMD ["./main"]
