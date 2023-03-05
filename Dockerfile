# Build frontend
FROM node:alpine as febuilder
ADD ui /usr/app/ui
WORKDIR /usr/app/ui
RUN npm install
RUN npm run build

# Build service
FROM golang:alpine AS builder
WORKDIR $GOPATH/src/github.com/indiependente/shrtnr/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY --from=febuilder /usr/app/ui/dist ./ui/dist
ADD . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /shrtnr

# Deployable image
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /shrtnr /app/
EXPOSE 7000
ENTRYPOINT ["/app/shrtnr"]
