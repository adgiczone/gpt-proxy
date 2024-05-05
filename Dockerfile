ARG baseImage=gpt-proxy:ubuntu22.04-go1.21.9

FROM ${baseImage} as builder
ENV PROJECT_PATH /go/src/little-gpt/gpt-proxy
WORKDIR $PROJECT_PATH
ADD . $PROJECT_PATH/

RUN make clean && make build

FROM ${baseImage}

ENV PROJECT_PATH /go/src/little-gpt/gpt-proxy
ENV ROOT_PATH=/gpt-proxy

COPY --from=builder $PROJECT_PATH/gpt-proxy ${ROOT_PATH}/

WORKDIR $ROOT_PATH