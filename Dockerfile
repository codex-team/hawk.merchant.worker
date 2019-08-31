FROM golang:stretch as builder
ARG BUILD_DIRECTORY=/build

# enable go modules
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# now copy go.mod and go.sum to the build path
RUN mkdir $BUILD_DIRECTORY
COPY ./src/go.mod $BUILD_DIRECTORY
COPY ./src/go.sum $BUILD_DIRECTORY

# download modules (for fast build due to docker caching)
WORKDIR $BUILD_DIRECTORY
RUN go mod download

# copy app sources and build
ADD ./src $BUILD_DIRECTORY
RUN go build -o hawk.merchant .

FROM alpine
ARG BUILD_DIRECTORY=/build

RUN apk add ca-certificates

WORKDIR /app
COPY --from=builder $BUILD_DIRECTORY .

CMD ["./hawk.merchant"]
