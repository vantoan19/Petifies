FROM alpine:latest

WORKDIR /app

COPY /config/user-service /app/config
COPY /build/server/userService /app

CMD [ "/app/userService"]