FROM alpine:3.10 AS final
WORKDIR "/app"
COPY ./config.json  /app
COPY $BINARY  /app
ENTRYPOINT ["/app/$BINARY"]