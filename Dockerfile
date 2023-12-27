FROM golang:1.21 AS nthu-campus-power

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /nthu-campus-power

EXPOSE 2112

# Run
CMD ["/nthu-campus-power"]
