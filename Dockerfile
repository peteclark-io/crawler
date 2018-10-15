FROM golang:1.11-alpine

ENV SRC_FOLDER="/src/crawler"
COPY . ${SRC_FOLDER}
WORKDIR ${SRC_FOLDER}

RUN apk --no-cache add git curl \
  && CGO_ENABLED=0 go build -a -o /artifacts/crawler

# Multi-stage build - copy only the certs and the binary into the image
FROM scratch
WORKDIR /
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /artifacts/* /

CMD [ "/crawler" ]
