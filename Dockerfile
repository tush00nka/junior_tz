FROM golang:1.24

LABEL maintainer="Alexey Borisoglebsky <endline00@ya.ru>"

WORKDIR /app

COPY . .

# Install and clean up dependencies
RUN go mod tidy
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]