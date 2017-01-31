## How to run in development mode

### 1. Resolve go dependencies

```
docker-compose -f development/docker-compose-tools.yml run --rm go-dep
```

### 2. docker-compose up

```
cd development
docker-compose up -d
```
