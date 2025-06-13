FROM golang:1.24.3 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /mservice

FROM alpine:latest
WORKDIR /app
COPY --from=build /mservice /mservice
EXPOSE 8080
HEALTHCHECK --interval=30s --start-period=5s --timeout=3s --retries=3 CMD wget --spider -q http://localhost:8080/health || exit 1
CMD ["/mservice"]