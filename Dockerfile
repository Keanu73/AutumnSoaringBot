FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o bot github.com/Keanu73/AutumnSoaringBot

CMD [ "/bot" ]