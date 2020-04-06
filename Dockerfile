FROM golang:1.14.1
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app

RUN go build main.go

CMD [ "./main" ]
