FROM alpine:latest

WORKDIR /app

COPY /cert /app/cert
COPY /server/config/media-service /app/config
COPY /build/server/mediaService /app

CMD [ "/app/mobileApiGateway"]