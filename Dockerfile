FROM golang:1.24 as builder

ARG ARG_GOPROXY
ENV GOPROXY $ARG_GOPROXY

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN make build


FROM quay.io/orvice/go-runtime:latest

ENV PROJECT_NAME simpleproxy

COPY --from=builder /home/app/bin/${PROJECT_NAME} .