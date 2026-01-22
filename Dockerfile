FROM --platform=linux/amd64 ubuntu:xenial AS builder

ENV OIIO_VER=1.6.18

RUN apt-get update && apt-get install --no-install-recommends -y -q \
    wget \
    g++ \
    make \
    cmake \
    fonts-freefont-ttf \
    libboost-thread-dev \
    libboost-system-dev \
    libboost-filesystem-dev \
    libboost-regex-dev \
    libopenexr-dev \
    libzlcore-dev \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libfreetype6-dev \
    libopencolorio-dev \
    && rm -rf /var/lib/apt/lists/*

RUN cd /home \
    && wget --no-check-certificate https://github.com/OpenImageIO/oiio/archive/Release-$OIIO_VER.tar.gz \
    && tar zxf Release-$OIIO_VER.tar.gz \
    && rm -f Release-$OIIO_VER.tar.gz

RUN cd /home/OpenImageIO-Release-$OIIO_VER \
    && make USE_OPENGL=0 USE_QT=0 USE_PYTHON=0 USE_PYTHON3=0 OIIO_BUILD_TOOLS=0 OIIO_BUILD_TESTS=0 VERBOSE=1


FROM --platform=linux/amd64 ubuntu:xenial AS tester

ENV GO_VER=1.24.0
ENV OIIO_VER=1.6.18
ENV GOPATH=/home/go
ENV PKGNAME=github.com/justinfx/openimageigo
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Go 1.24
RUN apt-get update && apt-get install --no-install-recommends -y -q \
    wget \
    ca-certificates \
    g++ \
    && rm -rf /var/lib/apt/lists/* \
    && wget --no-check-certificate https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VER}.linux-amd64.tar.gz \
    && rm go${GO_VER}.linux-amd64.tar.gz

# Copy OIIO build artifacts and required libraries
COPY --from=builder \
    /home/OpenImageIO-Release-${OIIO_VER}/dist/linux64/include/ \
    /home/OpenImageIO-Release-${OIIO_VER}/dist/linux64/include
COPY --from=builder \
    /home/OpenImageIO-Release-${OIIO_VER}/dist/linux64/lib/ \
    /home/OpenImageIO-Release-${OIIO_VER}/dist/linux64/lib/
COPY --from=builder /lib/x86_64-linux-gnu/libpng12.so* /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libz.so* /lib/x86_64-linux-gnu/
COPY --from=builder /usr/include /usr/include
COPY --from=builder /usr/lib/x86_64-linux-gnu /usr/lib/x86_64-linux-gnu
COPY --from=builder /usr/lib/libOpenColorIO.* /usr/lib/
COPY --from=builder /usr/share/fonts /usr/share/fonts

ENV CGO_CXXFLAGS="-I/home/OpenImageIO-Release-${OIIO_VER}/dist/linux64/include"
ENV CGO_LDFLAGS="-L/home/OpenImageIO-Release-${OIIO_VER}/dist/linux64/lib"
ENV LD_LIBRARY_PATH="/home/OpenImageIO-Release-${OIIO_VER}/dist/linux64/lib"

WORKDIR ${GOPATH}/src/${PKGNAME}

CMD ["go", "test", "-count", "1", "-v"]
