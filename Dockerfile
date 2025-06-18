FROM golang:1.24.3 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . ./
RUN swag init
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-service

FROM alpine:latest
WORKDIR /app

COPY --from=build /go-service /go-service
# COPY --from=build /swagger.json /swagger.json

EXPOSE 8080
HEALTHCHECK --interval=30s --start-period=5s --timeout=3s --retries=3 CMD wget --spider -q http://localhost:8080/health || exit 1
CMD ["/go-service"]