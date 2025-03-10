FROM golang:1.20-alpine AS build

WORKDIR /app

RUN apk add --no-cache git

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp .

# runtime
FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=build /app/myapp /app/

COPY --from=build /app/.env /app/

RUN chmod +x /app/myapp

EXPOSE 8080

ENTRYPOINT ["/app/myapp"]
