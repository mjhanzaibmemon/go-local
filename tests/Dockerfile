# create base layer
FROM 533267242017.dkr.ecr.eu-west-1.amazonaws.com/tools:go-trixie-1.25 AS build-test

ENV SOURCE_PATH=/go/app/source

WORKDIR ${SOURCE_PATH}

COPY app/go.mod app/go.sum ${SOURCE_PATH}/

RUN go mod download

COPY ./app ${SOURCE_PATH}

# create application binary files from base image
FROM build-test as build

# hadolint ignore=SC2086
RUN go build -o ${SOURCE_PATH}/bin/ ${SOURCE_PATH}/cmd/*

# create production image
FROM 533267242017.dkr.ecr.eu-west-1.amazonaws.com/cache/docker-hub/library/busybox:1.35-glibc as build-prod

ENV SOURCE_PATH=/go/app/source

WORKDIR /go/app

COPY --from=build $SOURCE_PATH/bin bin/
COPY --from=build $SOURCE_PATH/translations translations/
COPY --from=build /etc/ssl/certs /etc/ssl/certs
