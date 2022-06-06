ARG GO_VERSION=1.18
FROM golang:${GO_VERSION} AS builder
ENV GOPROXY="https://goproxy.cn"
ENV APP_PATH="/app/td"
WORKDIR "/app"
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags '-w -s' -a -o ${APP_PATH} .
RUN ls


FROM alpine:3.10 AS final
ENV APP_PATH="/app/td"
WORKDIR "/app"
COPY --from=builder ${APP_PATH} ${APP_PATH}
COPY ./config.json  /app
ENTRYPOINT ["/app/td"]