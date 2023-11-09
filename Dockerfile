FROM alpine:3.12

RUN mkdir "/app"
WORKDIR "/app"

COPY radishcache /app/app
ENTRYPOINT ["./app"]