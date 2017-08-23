# Builds a Docker image that allows you to run Jsonnet, kubecfg, and/or ksonnet
# on a file in your local directory. Specifically, this image contains:
#
# 1. Jsonnet, added to /usr/local/bin
# 2. ksonnet-lib, added to the Jsonnet library paths, so you can
#    compile against the ksonnet libraries without specifying the -J
#    flag.
# 3. kubecfg binary, added to /usr/local/bin
# 4. kubecfg lib, included in Jsonnet library paths via KUBECFG_JPATH,
#    similarly to (2) ksonnet-lib.
#
# USAGE: Define a function like `ksonnet` below, and then run:
#
#   `ksonnet <jsonnet-file-and-options-here>`
#
# ksonnet() {
#   docker run -it --rm   \
#     --volume "$PWD":/wd \
#     --workdir /wd       \
#     ksonnet             \
#     jsonnet "$@"
# }
#
# You can also define a similar function for `kubecfg`. Note that any required
# Jsonnet libraries specified by -J (required for compilation) need to be
# described relative to your working directory.

##############################################
# STAGE 1: build kubecfg
##############################################

FROM golang:1.8 as kubecfg-builder
# Keep this in sync with the corresponding ENV in other stages
ENV KUBECFG_VERSION v0.5.0

RUN go get github.com/ksonnet/kubecfg
WORKDIR /go/src/github.com/ksonnet/kubecfg
RUN git checkout tags/${KUBECFG_VERSION} -b ${KUBECFG_VERSION}
RUN CGO_ENABLED=1 GOOS=linux go install -a --ldflags '-linkmode external -extldflags "-static"' .

##############################################
# STAGE 2: build jsonnet and download ksonnet
##############################################

FROM alpine:3.6 as alpine-builder
# Keep this in sync with the corresponding ENV in other stages
ENV JSONNET_VERSION v0.9.4

# Get Jsonnet
RUN apk update && apk add git make g++
RUN git clone https://github.com/google/jsonnet.git
RUN cd jsonnet && git checkout tags/${JSONNET_VERSION} -b ${JSONNET_VERSION} && make -j4

# Get ksonnet-lib
RUN git clone https://github.com/ksonnet/ksonnet-lib.git

FROM alpine:3.6
ENV JSONNET_VERSION v0.9.4

# Copy kubecfg executable and lib files from previous stage
RUN mkdir -p /usr/share/kubecfg/
COPY --from=kubecfg-builder /go/bin/kubecfg /usr/local/bin/
COPY --from=kubecfg-builder /go/src/github.com/ksonnet/kubecfg/lib/ /usr/share/kubecfg/

# Copy jsonnet executable from previous stage
COPY --from=alpine-builder /jsonnet/jsonnet /usr/local/bin/

# Copy ksonnet-lib from previous stage, add to default Jsonnet -J path
RUN mkdir -p /usr/share/${JSONNET_VERSION}/
COPY --from=alpine-builder /ksonnet-lib/ksonnet.beta.2/ /usr/share/${JSONNET_VERSION}/

# Add all jsonnet paths to kubecfg -J path
ENV KUBECFG_JPATH /usr/share/kubecfg:/usr/share/${JSONNET_VERSION}
