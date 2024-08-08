FROM golang:1.22.6-bullseye AS build-base

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download -x

FROM build-base AS dev

RUN go install github.com/air-verse/air@latest

COPY .air.toml .

COPY . .

CMD ["air", "-c", ".air.toml"]

FROM build-base AS build-prod

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o /app/db.exe /app/cmd

FROM scratch AS prod

WORKDIR /app

COPY --from=build-prod /app/db.exe db.exe

EXPOSE 5432

CMD ["./db.exe"]
