#build stage
FROM golang:1.17 AS builder
ENV GO111MODULE=on \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64 \
GOSUMDB=off
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

#final stage
FROM scratch
COPY --from=builder /dist/main /
ENTRYPOINT ["/main"]