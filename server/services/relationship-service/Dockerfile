FROM alpine:latest

WORKDIR /app

COPY /server/config/relationship-service /app/config
COPY /build/server/relationshipService /app

EXPOSE 8888

CMD [ "/app/relationshipService"]