FROM golang:alpine AS builder
WORKDIR $GOPATH/src/github.com/indiependente/shrtnr/
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get github.com/GeertJohan/go.rice
RUN go get github.com/GeertJohan/go.rice/rice
ADD . .
RUN rice embed-go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /shrtnr
RUN rm -f rice-box.go

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /shrtnr /app/
EXPOSE 7000
ENTRYPOINT ["/app/shrtnr"]
