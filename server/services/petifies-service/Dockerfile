FROM alpine:latest

WORKDIR /app

COPY /server/config/petifies-service /app/config
COPY /build/server/petifiesService /app
COPY /server/services/petifies-service/db/migrations /migrations

EXPOSE 8888

CMD [ "/app/petifiesService"]