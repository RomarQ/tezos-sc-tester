# VisualTez Testing

[![CI](https://github.com/RomarQ/visualtez-testing/actions/workflows/pipeline.yaml/badge.svg)](https://github.com/RomarQ/visualtez-testing/actions/workflows/pipeline.yaml)

## Install dependencies and compile application

```sh
make
```

## Start the application

```sh
make start
```

## Run Tests

```sh
make test
```

## Run docker container

```sh
// arm64
docker run -p 5000:5000 --name testing-api -d ghcr.io/romarq/visualtez-testing:0.0.8_arm64
// amd64
docker run -p 5000:5000 --name testing-api -d ghcr.io/romarq/visualtez-testing:0.0.8_amd64
```

### Configuration

The configuration can be modified at [./config.yaml](./config.yaml).
