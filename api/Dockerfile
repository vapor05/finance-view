FROM golang:1.17-alpine3.15 as builder

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
COPY gqlgen.yml .
COPY server.go .
COPY graph graph
COPY pkg pkg

RUN go build server.go

FROM alpine:3.15

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/server .

EXPOSE 8080

CMD [ "./server" ]