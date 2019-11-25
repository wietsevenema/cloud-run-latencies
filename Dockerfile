FROM golang:1.13
RUN go get github.com/rakyll/hey
COPY main.go .
ENTRYPOINT ["go", "run", "main.go"]