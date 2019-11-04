FROM golang:1.13.4-alpine

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go build -o main /app/cmd/gosave/main.go

RUN adduser -S -D -H -h /app appuser
RUN chown -R appuser /app
USER appuser

CMD ["./main"]