FROM golang

# Download dep for dependency management
RUN go get github.com/golang/dep/cmd/dep

# Download gin for live reload (Usage: gin --path src --port 8081 run server.go serve)
RUN go get github.com/codegangsta/gin

WORKDIR /go/src/app

ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock

COPY ./ ./

RUN dep ensure -v
CMD ["go", "run", "server.go", "serve", "--env", "docker"]

EXPOSE 8081