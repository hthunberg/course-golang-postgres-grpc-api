#########################
# Docker wormhole pattern
# Testcontainers will automatically detect if it's inside a container and instead of "localhost" will use the default gateway's IP.
#
# However, additional configuration is required if you use volume mapping. The following points need to be considered:
# - The docker socket must be available via a volume mount
# - The 'local' source code directory must be volume mounted at the same path inside the container that Testcontainers runs in, 
#   so that Testcontainers is able to set up the correct volume mounts for the containers it spawns.
#
# docker run -it --rm -v $PWD:$PWD -w $PWD -v /var/run/docker.sock:/var/run/docker.sock maven:3 mvn test
#########################
FROM golang:1.21

# We might consider disablinbg ruyk when running testcontainers inside container
ENV TESTCONTAINERS_RYUK_DISABLED=true

RUN apt-get update && apt-get install -y \
  make \
  gcc

WORKDIR /usr/src/app

COPY Makefile ./

# To /usr/src/app/go.sum
# To /usr/src/app/go.mod
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN ls -lart /usr/src/app/internal/app/bank/integration

RUN pwd

ENTRYPOINT [ "make", "integrationtest" ]
CMD []