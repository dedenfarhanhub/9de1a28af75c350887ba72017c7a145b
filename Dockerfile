FROM golang:1.21.0

ENV GIN_MODE release

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY . /app

RUN go install github.com/air-verse/air@latest

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]

CMD ["air", "-c", ".air.toml"]