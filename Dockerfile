FROM golang:alpine

WORKDIR /app

COPY . .

RUN mkdir dist

RUN go build -ldflags="-s -w" -o dist/tras .

CMD [ "dist/tras" ]

STOPSIGNAL SIGTERM