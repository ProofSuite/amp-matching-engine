FROM golang

# Download dep for dependency management
RUN go get github.com/golang/dep/cmd/dep
# # Download gin for live reload (Usage: gin --path src --port 8081 run server.go serve)
# # RUN go get github.com/codegangsta/gin
# WORKDIR /go/src/app

RUN mkdir -p /go/src/github.com/Proofsuite/amp-matching-engine
WORKDIR /go/src/github.com/Proofsuite/amp-matching-engine

ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock
RUN dep ensure -vendor-only

COPY ./ ./

CMD ["go", "run", "main.go"]

EXPOSE 8081