# Crawler

## Build

Using regular go build:

```
go get -u github.com/peteclark-io/crawler
go build
```

Using go modules, first git clone the repo to a directory outside the $GOPATH, and ensure you are using go 1.11+, then:

```
go build
```

To build the docker image:

```
docker build -t crawler:local .
```

## Usage

With the docker image:

```
docker run -ti crawler:local
```

With go installed binary:

```
crawler
```
