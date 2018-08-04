FROM ubuntu:xenial as builder

ENV OIIO_VER 1.7.19

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
    libopencolorio-dev

RUN cd /home \
    && wget --no-check-certificate https://github.com/OpenImageIO/oiio/archive/Release-$OIIO_VER.tar.gz \
	&& tar zxf Release-$OIIO_VER.tar.gz \
	&& rm -f Release-$OIIO_VER.tar.gz

RUN cd /home/oiio-Release-$OIIO_VER \
    && make USE_OPENGL=0 USE_QT=0 USE_PYTHON=0 USE_PYTHON3=0 VERBOSE=1


FROM golang:1.10-stretch

ENV OIIO_VER 1.7.19
ENV GOPATH /home/go
ENV PKGNAME github.com/justinfx/openimageigo

ENV CGO_CXXFLAGS="-I/home/oiio-Release-${OIIO_VER}/dist/linux64/include"
ENV CGO_LDFLAGS="-L/home/oiio-Release-${OIIO_VER}/dist/linux64/lib"
ENV LD_LIBRARY_PATH="/home/oiio-Release-${OIIO_VER}/dist/linux64/lib"

COPY --from=builder \
    /lib/x86_64-linux-gnu/libz.* \
    /lib/x86_64-linux-gnu/libpng12.* \
    /lib/x86_64-linux-gnu/liblzma.* \
    /lib/x86_64-linux-gnu/
COPY --from=builder /usr/include /usr/include
COPY --from=builder /usr/lib/x86_64-linux-gnu /usr/lib/x86_64-linux-gnu
COPY --from=builder /usr/lib/libOpenColorIO.* /usr/lib/
COPY --from=builder /usr/share/fonts /usr/share/fonts
COPY --from=builder \
    /home/oiio-Release-${OIIO_VER}/dist/linux64/include/ \
    /home/oiio-Release-${OIIO_VER}/dist/linux64/include
COPY --from=builder \
    /home/oiio-Release-${OIIO_VER}/dist/linux64/lib/ \
    /home/oiio-Release-${OIIO_VER}/dist/linux64/lib/

#ADD . /home/go/src/$PKGNAME

WORKDIR ${GOPATH}/src/${PKGNAME}

CMD go test -count 1 -v $PKGNAME
