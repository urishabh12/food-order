FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/server /app/server

ENV ADDR=:8080
ENV BASE_PATH=/api
ENV API_KEY=apitest

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
