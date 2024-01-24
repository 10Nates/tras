FROM golang:alpine

WORKDIR /app

COPY . .

RUN mkdir dist

RUN go build -ldflags="-s -w" -o dist/tras .

#unnecessary for operation
RUN rm *.go
RUN rm go.*
RUN rm -rf db/
RUN rm -rf discordless/

CMD [ "dist/tras" ]

STOPSIGNAL SIGTERM