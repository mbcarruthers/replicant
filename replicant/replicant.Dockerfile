# syntax=docker/dockerfile:1
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o replicant ./cmd/

RUN chmod +x /app/replicant

# create a small image and copy the executable

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/replicant /app/

EXPOSE "8000:8000"

CMD ["/app/replicant"]