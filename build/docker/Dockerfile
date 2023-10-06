#############
# Build stage
#############
FROM golang:1.21 as builder

RUN apt-get update && apt-get install -y \
  make \
  gcc

# magicpak enables you to build minimal docker images without any bothersome preparation such as static linking.
ADD https://github.com/coord-e/magicpak/releases/latest/download/magicpak-x86_64-unknown-linux-musl /usr/bin/magicpak
RUN chmod +x /usr/bin/magicpak

# add the user to run container based on scratch, since scratch dont have addgroup or adduser.
RUN useradd -u 10001 scratchuser

WORKDIR /usr/src/app

COPY Makefile ./

# Install all tools needed by make targets
RUN make all-tools

# To /usr/src/app/go.sum
# To /usr/src/app/go.mod
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

# Run /usr/src/app/ $ go <commands> 
RUN go mod download && go mod verify

# Copy entire source to /usr/src/app/
COPY . ./

# Run /usr/src/app/ $ make <targets> 
RUN make lint
RUN make build

# Copy binary from /usr/src/app/bin/bank to /app/bin/bank
RUN mkdir -p /app/bin
RUN cp bin/bank /app/bin/bank

# Analyzes runtime dependencies of bank executable statically and
# put everything bank executable needs in runtime /bundle directory. 
RUN mkdir -p /bundle
RUN /usr/bin/magicpak -v /app/bin/bank /bundle

##############
# Deploy stage
##############

#  Use Docker's reserved, minimal image, scratch , as a starting point for building container.
FROM scratch

# copy scratchuser created in builder
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser

# Copy bank executable and runtime dependencies to build a minimal docker image
COPY --from=builder /bundle /.

ENV TZ=UTC

# Only UTF8 locale for this application 
COPY --from=builder /usr/lib/locale/C.utf8 /usr/lib/locale/C.utf8  
ENV LC_CTYPE=C.utf8

# Copy time zone info
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo 

WORKDIR /app/bin

ENTRYPOINT [ "/app/bin/bank" ]

CMD []