FROM golang:latest

WORKDIR /app
COPY . /app

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["gpt-4"]
