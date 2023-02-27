FROM alpine:latest

WORKDIR /app

COPY /server/config/user-service /app/config
COPY /build/server/userService /app
COPY /server/services/user-service/db/migrations /migrations

EXPOSE 8888

CMD [ "/app/userService"]