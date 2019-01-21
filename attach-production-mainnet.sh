AMP_ENABLE_TLS=true; \
AMP_MONGODB_SHARD_URL_1=ampcluster0-shard-00-00-xzynf.mongodb.net:27017; \
AMP_MONGODB_SHARD_URL_2=ampcluster0-shard-00-01-xzynf.mongodb.net:27017; \
AMP_MONGODB_SHARD_URL_3=ampcluster0-shard-00-02-xzynf.mongodb.net:27017; \
fresh

# go run --race main.go