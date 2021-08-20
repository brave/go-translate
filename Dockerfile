FROM ubuntu:18.04
LABEL image=go-translate-server

# Get all the dependencies
RUN apt update
RUN apt install -y git autoconf automake libtool curl make g++ unzip
RUN apt install -y wget git-lfs pkg-config software-properties-common
RUN git lfs install

# Install protobuf
RUN cd \
    && git clone https://github.com/google/protobuf.git \
    && cd protobuf \
    && git checkout tags/v3.11.0 \
    && git submodule update --init --recursive \
    && ./autogen.sh \
    && ./configure \
    && make \
    && make install \
    && ldconfig

# Install go
RUN cd \
    && wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Install CLD3 for go
RUN go get github.com/jmhodges/gocld3/cld3

# Install fresh version of cmake (required for building Bergamot translator)
RUN wget -O - https://apt.kitware.com/keys/kitware-archive-latest.asc 2>/dev/null | apt-key add - \
    && apt-add-repository -y 'deb https://apt.kitware.com/ubuntu/ bionic main' \
    && apt update \
    && apt install -y cmake

# Install MKL (required for building Bergamot translator)
RUN wget https://apt.repos.intel.com/intel-gpg-keys/GPG-PUB-KEY-INTEL-SW-PRODUCTS-2019.PUB \
    && apt-key add GPG-PUB-KEY-INTEL-SW-PRODUCTS-2019.PUB \
    && sh -c 'echo deb https://apt.repos.intel.com/mkl all main > /etc/apt/sources.list.d/intel-mkl.list' \
    && apt update \
    && apt install -y intel-mkl-64bit-2018.2-046

# Install Bergamot translator and all its dependencies
RUN cd \
    && git clone https://github.com/browsermt/bergamot-translator.git \
    && cd bergamot-translator \
    && git submodule update --init --recursive \
    && mkdir build-native \
    && cd build-native \
    && cmake -DSSPLIT_USE_INTERNAL_PCRE2=ON .. \
    && make

# Install models
RUN cd \
    && git clone https://github.com/mozilla/firefox-translations-models.git \
    && cd firefox-translations-models \
    && cd models \
    && gunzip -r . \
    && cp -r dev/* . \
    && cp -r prod/* .

COPY run_server.sh /run_server.sh
