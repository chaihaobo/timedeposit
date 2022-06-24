FROM alpine:3.10 AS final
WORKDIR "/app"
COPY td config.json ./
ENTRYPOINT ["/app/td"]