FROM golang:1.19-alpine

RUN apk add --no-cache tzdata
ENV TZ=Europe/London

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o bot github.com/Keanu73/AutumnSoaringBot

CMD [ "/app/bot" ]