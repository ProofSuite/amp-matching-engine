FROM golang

# Download dep for dependency management
RUN go get github.com/golang/dep/cmd/dep

# Download gin for live reload (Usage: gin --path src --port 8081 run server.go serve)
# RUN go get github.com/codegangsta/gin

WORKDIR /go/src/app

# ENV GO_ENV=DOCKER
ENV AMP_ETHEREUM_NODE_URL=${AMP_ETHEREUM_NODE_URL}
ENV AMP_MONGO_URL=${AMP_MONGO_URL}
ENV AMP_MONGO_DBNAME=${AMP_MONGO_DBNAME}
ENV AMP_REDIS_URL=${AMP_REDIS_URL}
ENV AMP_RABBITMQ_URL=${AMP_RABBITMQ_URL}
ENV AMP_EXCHANGE_CONTRACT_ADDRESS=${AMP_EXCHANGE_CONTRACT_ADDRESS}
ENV AMP_WETH_CONTRACT_ADDRESS=${AMP_WETH_CONTRACT_ADDRESS}
ENV AMP_FEE_ACCOUNT_ADDRESS=${AMP_FEE_ACCOUNT_ADDRESS}

ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock

COPY ./ ./

RUN dep ensure -v
CMD ["go", "run", "main.go"]

EXPOSE 8081