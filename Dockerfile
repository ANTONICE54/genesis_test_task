#Build stage
FROM golang:1.21.4-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/util/mailTemplates ./util/mailTemplates
COPY db/migration ./db/migration
RUN apk add --no-cache tzdata
ENV TZ=Europe/Kyiv




EXPOSE 8080
CMD [ "/app/main" ]