FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /server ./cmd
FROM alpine:3.21
RUN adduser -D -H app
USER app
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
