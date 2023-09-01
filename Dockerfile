FROM golang:alpine as build

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . ./
COPY .docker.env .env

RUN go build -o main .

EXPOSE 3000

ENTRYPOINT [ "./main" ]