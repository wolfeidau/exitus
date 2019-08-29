FROM debian

RUN apt update && apt install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

ADD ./exitus /exitus

ENTRYPOINT ["/exitus"]
