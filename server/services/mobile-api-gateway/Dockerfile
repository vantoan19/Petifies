FROM alpine:latest

WORKDIR /app

COPY /cert /app/cert
COPY /config/mobile-api-gateway /app/config
COPY /build/server/mobileApiGateway /app

CMD [ "/app/mobileApiGateway"]