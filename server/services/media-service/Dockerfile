FROM alpine:latest

WORKDIR /app

RUN mkdir /storage

COPY /cert /app/cert
COPY /server/config/media-service /app/config
COPY /build/server/mediaService /app

CMD [ "/app/mediaService"]