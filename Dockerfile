FROM golang:alpine

WORKDIR /usr/src/stratum

RUN apk add git

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/stratum main.go

RUN addgroup -S stratum && adduser -S stratum -G stratum
USER stratum

ENTRYPOINT ["stratum"]