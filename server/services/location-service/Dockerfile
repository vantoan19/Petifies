FROM alpine:latest

WORKDIR /app

COPY /server/config/location-service /app/config
COPY /build/server/locationService /app
COPY /server/services/location-service/db/migrations /migrations

EXPOSE 8888

CMD [ "/app/locationService"]