#Multi stage docker file

#Build stage
FROM golang:1.18-alpine AS builder 
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Run stage
FROM alpine:3.13
WORKDIR /app
COPY template ./template
COPY --from=builder /app/main .
COPY app.env .
COPY image ./image

EXPOSE 8080
CMD ["/app/main"]


