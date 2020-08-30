FROM golang:1.15-alpine AS build

ENV         OPENCV_VERSION=4.4.0
ENV         OPENALPR_VERSION=2.3.0
ENV         TESSERACT_VERSION=4.1.1
ENV         LEPTONICA_VERSION=1.80.0
ENV         LOG4CPLUS_VERSION=2.0.5

ARG         LD_LIBRARY_PATH=/usr/local/lib
ARG         MAKEFLAGS="-j2"
ARG         PKG_CONFIG_PATH="/usr/local/share/pkgconfig:/usr/local/lib/pkgconfig:/usr/local/lib64/pkgconfig"
ARG         PREFIX=/usr/local
ARG         LD_LIBRARY_PATH="/usr/local/lib:/usr/local/lib64:/usr/lib64:/usr/lib:/lib64:/lib"


RUN     buildDeps="tiff zlib libpng libjpeg \
                    tiff-dev zlib-dev libpng-dev jpeg-dev \
                    tiff            \
                    autoconf        \
                    m4              \
                    linux-headers   \
                    build-base      \
                    gcc             \
                    make            \
                    cmake           \
                    pkgconfig       \
                    automake        \
                    ca-certificates \
                    g++             \
                    curl            \
                    git             \
                    curl-dev        \
                    libtool         \
                    wget            \
                    tesseract-ocr-dev" && \
        apk update && apk upgrade&& \
        apk add --no-cache ${buildDeps}

RUN DIR=/tmp/opencv && \
    mkdir -p ${DIR} && \
    cd ${DIR} && \
    curl -sL https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.tar.gz | \
    tar -zx --strip-components=1 && \
    mkdir -p ${DIR}/build && \
    cd ${DIR}/build && \
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=${PREFIX} .. && \
    make && \
    make install && \
    rm -rf ${DIR}

RUN DIR=/tmp/leptonica && \
    mkdir -p ${DIR} && \
    cd ${DIR} && \
    curl -sL http://www.leptonica.org/source/leptonica-${LEPTONICA_VERSION}.tar.gz | \
    tar -zx --strip-components=1 && \
    ./autogen.sh && \
    ./configure --prefix="${PREFIX}" && \
    make && \
    make install && \
    rm -rf ${DIR}

RUN DIR=/tmp/openalpr && \
    git clone --depth 1 https://github.com/openalpr/openalpr ${DIR} && \
    mkdir -p ${DIR}/src/build && \
    cd ${DIR}/src/build && \
    cmake -DCMAKE_INSTALL_PREFIX:PATH=${PREFIX} -DWITH_BINDING_PYTHON=OFF -DWITH_BINDING_JAVA=OFF -DWITH_TESTS=OFF -DWITH_DAEMON=OFF .. && \
    make && \
    make install && \
    rm -rf ${DIR}

RUN DIR=/tmp/log4cplus && \
    mkdir -p ${DIR} && \
    cd ${DIR} && \
    curl -sL https://versaweb.dl.sourceforge.net/project/log4cplus/log4cplus-stable/${LOG4CPLUS_VERSION}/log4cplus-${LOG4CPLUS_VERSION}.tar.gz | \
    tar -zx --strip-components=1 && \
    ./configure --prefix="${PREFIX}" && \
    make && \
    make install && \
    rm -rf ${DIR}

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN BLDDIR=/go/bin make build

FROM alpine:3.12

ENV LD_LIBRARY_PATH="/usr/local/lib"
ENV PKG_CONFIG_PATH="/usr/local/share/pkgconfig:/usr/local/lib/pkgconfig:/usr/local/lib64/pkgconfig"
ENV LD_LIBRARY_PATH="/usr/local/lib:/usr/local/lib64:/usr/lib64:/usr/lib:/lib64:/lib"

RUN apk update && apk upgrade && apk add --no-cache ca-certificates tesseract-ocr-dev tiff zlib libpng libjpeg tiff-dev zlib-dev libpng-dev jpeg-dev
WORKDIR /app

COPY --from=build /usr/local/lib /usr/local/lib
COPY --from=build /usr/local/lib64 /usr/local/lib64
COPY --from=build /usr/local/share/openalpr/ /usr/local/share/openalpr/

COPY ./config/ua.conf /usr/local/share/openalpr/runtime_data/config/ua.conf
COPY ./config/ua.patterns /usr/local/share/openalpr/runtime_data/postprocess/ua.patterns
RUN cp /usr/local/share/openalpr/runtime_data/region/eu.xml /usr/local/share/openalpr/runtime_data/region/ua.xml

COPY --from=build /go/bin/ ./
COPY ./config/config.toml ./config/config.toml

EXPOSE 8080

CMD ["./server"]