FROM alpine:latest

WORKDIR /app

COPY /config/user-service /app/config
COPY /build/server/userService /app
COPY /server/services/user-services/user-service/db/migrations /migrations

CMD [ "/app/userService"]