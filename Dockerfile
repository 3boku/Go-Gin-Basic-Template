FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/init/main.go

# 최종 이미지에 필요한 파일들을 준비합니다
WORKDIR /dist
RUN mkdir -p secret
RUN cp /build/main .
RUN cp /build/secret/.env secret/.env

FROM scratch

WORKDIR /app

COPY --from=builder /dist/main .
COPY --from=builder /dist/secret /app/secret

ENTRYPOINT ["/app/main"]
