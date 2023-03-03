FROM alpine:latest

WORKDIR /app

COPY /server/config/post-service /app/config
COPY /build/server/postService /app
COPY /server/services/post-service/db/migrations /migrations

EXPOSE 8888

CMD [ "/app/postService"]