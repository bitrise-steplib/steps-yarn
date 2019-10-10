FROM golang:1.12

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y \
  rsync \
  curl \
  git \
  sudo

RUN curl -fL https://github.com/bitrise-io/bitrise/releases/download/1.35.0/bitrise-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise
RUN chmod +x /usr/local/bin/bitrise
RUN bitrise setup

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y yarn

RUN mkdir -p /go/src/github.com/bitrise-steplib/steps-yarn
ADD . /go/src/github.com/bitrise-steplib/steps-yarn
WORKDIR /go/src/github.com/bitrise-steplib/steps-yarn

ENTRYPOINT bitrise run test