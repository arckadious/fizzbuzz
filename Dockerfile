FROM debian:latest

# Set Timezone
ENV TZ Europe/Paris

ENV GOROOT=/usr/local/go
ENV GOBIN=$GOROOT/bin
ENV PATH=$GOROOT/bin:$PATH 
ARG GO_VERSION=1.19.8

# for latest version, set env value to : master
ARG AIR_VERSION=v1.43.0 

RUN apt-get update -yq \
&& apt-get install --no-install-recommends curl git ca-certificates gcc libc6-dev -yq \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*


RUN curl -o go${GO_VERSION}.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz \
&& tar -xf go${GO_VERSION}.tar.gz -C /usr/local \
&& rm go${GO_VERSION}.tar.gz

RUN go version

RUN touch /var/log/messages.log

# RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/${AIR_VERSION}/install.sh | sh -s -- -b ${GOPATH}
RUN go install github.com/cosmtrek/air@${AIR_VERSION}

COPY . /fizzbuzz-api/

WORKDIR /fizzbuzz-api

# keep container Up to avoid terminating
# ENTRYPOINT [ "tail", "-f", "/dev/null" ] 

# Register logs from 'make run' and 'make srv' commands to dozzle.
ENTRYPOINT [ "tail", "-f", "/var/log/messages.log" ] 