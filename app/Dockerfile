FROM golang:1.22.2-bullseye AS build-base

WORKDIR /app

COPY backend/go.mod backend/go.sum /app/backend/

RUN cd ./backend && go mod download -x

FROM build-base AS dev

RUN go install github.com/air-verse/air@latest

COPY .air.toml .

COPY backend/ backend/

COPY frontend/ frontend/

CMD ["air", "-c", ".air.toml"]

FROM build-base AS build-prod

COPY backend/ backend/

RUN cd ./backend && go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o /app/go-tic-tac-toe.exe /app/backend/cmd

FROM scratch AS prod

WORKDIR /app

COPY frontend/ frontend/

COPY --from=build-prod /app/go-tic-tac-toe.exe go-tic-tac-toe.exe

EXPOSE 3000

CMD ["./go-tic-tac-toe.exe"]