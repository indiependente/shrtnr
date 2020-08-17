FROM golang:1.15-alpine AS builder
WORKDIR $GOPATH/src/github.com/indiependente/shrtnr/
COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /shrtnr

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /shrtnr /app/
EXPOSE 7000
ENTRYPOINT ["/app/shrtnr"]
