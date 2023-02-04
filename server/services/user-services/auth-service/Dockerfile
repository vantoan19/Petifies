FROM alpine:latest

WORKDIR /app

COPY /config/auth-service /app/config
COPY /build/server/authService /app

CMD [ "/app/authService"]