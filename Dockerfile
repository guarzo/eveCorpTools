FROM golang:1.23.0 as builder
ARG CGO_ENABLED=0
ARG VERSION

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o zkill

FROM gcr.io/distroless/static AS final

COPY --from=builder /app/static /static
COPY --from=builder /app/zkill /zkill

VOLUME data
ENTRYPOINT ["/zkill"]
