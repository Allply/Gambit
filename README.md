# Gambit

Gambit is a db to store recommendations for building an allply profile.

The idea is to build up enough recommendations that querying with a certain similarity will provide quick advice.

The vector DB is built on the weaviate go client and their provided docker image

[https://weaviate.io/developers/weaviate](https://weaviate.io/developers/weaviate)

## Installation

init the container

```bash
docker compose up -d
```

init the db

```bash
go run init.go --init --load
```

query the db

```bash
go run init.go --query "Built websites with Python and Javascript"
```
