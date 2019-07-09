FROM golang

RUN apt-get update && \
    apt-get install -y \
      # build tools, for compiling
      build-essential \
      # install curl to fetch dev things
      curl \
      # we'll need git for fetching golang deps
      git

# install dep (not using it yet, but probably will switch to it)
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# setup the app dir/working directory
RUN mkdir -p /go/src/github.com/r0h4n/log_agg
WORKDIR /go/src/github.com/r0h4n/log_agg

# copy the source
COPY . .

# fetch deps

EXPOSE 6360:6360

ENTRYPOINT ./dist/log_agg -c ./config.json

