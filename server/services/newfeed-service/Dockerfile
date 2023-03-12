FROM alpine:latest

WORKDIR /app

COPY /server/config/newfeed-service /app/config
COPY /build/server/newfeedService /app
COPY /server/services/newfeed-service/db/migrations /migrations

EXPOSE 8888

CMD [ "/app/newfeedService"]