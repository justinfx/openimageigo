FROM --platform=linux/amd64 ubuntu:jammy AS builder

ENV OIIO_VER=2.2.16.0

RUN apt-get update && apt-get install --no-install-recommends -y -q \
    wget \
    git \
    ca-certificates \
    g++ \
    make \
    cmake \
    pkg-config \
    fonts-freefont-ttf \
    libboost-thread-dev \
    libboost-system-dev \
    libboost-filesystem-dev \
    libopenexr-dev \
    libilmbase-dev \
    libz-dev \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libfreetype6-dev \
    libopencolorio-dev \
    libfmt-dev \
    && rm -rf /var/lib/apt/lists/*

RUN cd /home \
    && wget --no-check-certificate https://github.com/AcademySoftwareFoundation/OpenImageIO/archive/refs/tags/v${OIIO_VER}.tar.gz \
    && tar zxf v${OIIO_VER}.tar.gz \
    && rm -f v${OIIO_VER}.tar.gz

RUN cd /home/OpenImageIO-${OIIO_VER} \
    && cmake -B build \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_INSTALL_PREFIX=/home/oiio-dist \
        -DOIIO_BUILD_TOOLS=OFF \
        -DOIIO_BUILD_TESTS=OFF \
        -DUSE_PYTHON=OFF \
        -DUSE_QT=OFF \
        -DUSE_OPENGL=OFF \
        -DUSE_FREETYPE=ON \
        -DOpenImageIO_BUILD_MISSING_DEPS=all \
    && cmake --build build -j4 \
    && cmake --install build


FROM --platform=linux/amd64 ubuntu:jammy AS tester

ENV GO_VER=1.24.0
ENV OIIO_VER=2.2.16.0
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Go 1.24 and FreeSans font (required for OIIO text rendering)
RUN apt-get update && apt-get install --no-install-recommends -y -q \
    wget \
    ca-certificates \
    g++ \
    fontconfig \
    fonts-freefont-ttf \
    && rm -rf /var/lib/apt/lists/* \
    && wget --no-check-certificate https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VER}.linux-amd64.tar.gz \
    && rm go${GO_VER}.linux-amd64.tar.gz

# Copy OIIO build artifacts and required libraries
COPY --from=builder /home/oiio-dist/include/ /home/oiio-dist/include/
COPY --from=builder /home/oiio-dist/lib/ /home/oiio-dist/lib/
COPY --from=builder /usr/include /usr/include
COPY --from=builder /usr/lib/x86_64-linux-gnu /usr/lib/x86_64-linux-gnu
COPY --from=builder /usr/lib/libOpenColorIO.* /usr/lib/

ENV CGO_CXXFLAGS="-I/home/oiio-dist/include"
ENV CGO_LDFLAGS="-L/home/oiio-dist/lib"
ENV LD_LIBRARY_PATH="/home/oiio-dist/lib"
ENV OPENIMAGEIO_FONTS="/usr/share/fonts"

# Initialize font cache for OIIO
RUN fc-cache -fv

WORKDIR /workdir

CMD ["go", "test", "-count", "1", "-v"]
