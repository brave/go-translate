FROM golang:1.16 as builder

RUN apt update
RUN apt install -y git autoconf libtool cmake git-lfs
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

RUN wget --https-only --secure-protocol=PFS https://apt.repos.intel.com/intel-gpg-keys/GPG-PUB-KEY-INTEL-SW-PRODUCTS-2019.PUB \
    && apt-key add GPG-PUB-KEY-INTEL-SW-PRODUCTS-2019.PUB \
    && sh -c 'echo deb https://apt.repos.intel.com/mkl all main > /etc/apt/sources.list.d/intel-mkl.list' \
    && apt update \
    && apt install -y intel-mkl-64bit-2018.2-046

# Install Bergamot translator and all its dependencies
RUN cd \
    && git clone https://github.com/browsermt/bergamot-translator.git \
    && cd bergamot-translator \
    && git checkout 63120c174e3edfd664175d4a2be095d8b50a112f \
    && git submodule update --init --recursive \
    && mkdir build-native \
    && cd build-native \
    && cmake -DSSPLIT_USE_INTERNAL_PCRE2=ON .. \
    && make

# Install models
RUN cd \
    && git clone https://github.com/mozilla/firefox-translations-models.git \
    && cd firefox-translations-models \
    && git checkout b22ca725bb102c034dabf3871d7349f2aca8d73d \
    && git lfs fetch --all \
    && cd models \
    && gunzip -r . \
    && cp -r dev/* . \
    && cp -r prod/* .

WORKDIR /src/
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-w -s" \
    -o go-translate main.go

FROM ubuntu:20.04 as artifact

RUN apt-get update
RUN apt-get install -y git autoconf libtool g++ make

# Install protobuf
RUN cd \
    && git clone https://github.com/google/protobuf.git \
    && cd protobuf \
    && git checkout tags/v3.11.0 \
    && git submodule update --init --recursive \
    && ./autogen.sh \
    && ./configure --disable-dependency-tracking \
    && make \
    && make install \
    && ldconfig

COPY --from=builder /src /root/app
COPY --from=builder /root/firefox-translations-models /root/firefox-translations-models
COPY --from=builder /root/bergamot-translator /root/bergamot-translator

WORKDIR /root/app/

EXPOSE 8195
CMD ["./go-translate"]
