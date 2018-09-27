FROM golang:1.10.1 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/BlockClusterApp/blockcluster-daemon
# copies the Gopkg.toml and Gopkg.lock to WORKDIR
COPY Gopkg.toml Gopkg.lock ./
# install the dependencies without checking for go code

RUN dep ensure -vendor-only

COPY . ./

# RUN go build -o app /go/src/github.com/BlockClusterApp/blockcluster-daemon/**/*.go

CMD ["go", "run", "src/main.go"]
