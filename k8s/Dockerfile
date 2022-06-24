ARG GO_VERSION=1.18
FROM golang:${GO_VERSION} AS builder
ENV GOPROXY="https://proxy.golang.org,direct"
ENV APP_PATH="/app/td"
WORKDIR "/app"
COPY . .
#RUN git config --global http.extraheader "PRIVATE-TOKEN: glpat-UpdFVcqHGC_kxYB_8KsF"
#RUN git config --global url."git@gitlab.com:".insteadof "https://gitlab.com/"
RUN git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.com".insteadOf "https://gitlab.com"
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags '-w -s' -a -o ${APP_PATH} .
RUN ls


FROM alpine:3.10 AS final
ENV APP_PATH="/app/td"
WORKDIR "/app"
COPY --from=builder ${APP_PATH} ${APP_PATH}
COPY ./config.json  /app
ENTRYPOINT ["/app/td"]
