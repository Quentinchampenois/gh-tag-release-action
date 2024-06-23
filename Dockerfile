FROM golang:1.22-alpine

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o /action

CMD ["/action"]