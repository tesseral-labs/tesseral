FROM perl:5.40-slim

ARG CHARSET=UTF-8

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git \
        make && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/

RUN git clone https://github.com/darold/pgFormatter.git

WORKDIR /usr/src/pgFormatter

RUN perl Makefile.PL && \
    make && \
    make install

VOLUME ["/work"]
WORKDIR /work

ENTRYPOINT ["/usr/local/bin/pg_format"]
CMD ["--help"]
