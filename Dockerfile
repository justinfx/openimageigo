FROM ubuntu:xenial

ARG OIIO_VER
ENV OIIO_VER 1.7.19

ARG PKGNAME
ENV PKGNAME github.com/justinfx/openimageigo

RUN apt-get update && apt-get install -y wget

RUN cd /home && wget https://github.com/OpenImageIO/oiio/archive/Release-$OIIO_VER.tar.gz \
	&& tar zxf Release-$OIIO_VER.tar.gz \
	&& rm -f Release-$OIIO_VER.tar.gz

RUN apt-get install --no-install-recommends -y -q \
    g++ \
    make \
    cmake \
    golang-1.10-go \
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

RUN cd /home/oiio-Release-$OIIO_VER \
    && make USE_OPENGL=0 USE_QT=0 USE_PYTHON=0 USE_PYTHON3=0 VERBOSE=1

ENV PATH /usr/lib/go-1.10/bin:$PATH
ENV GOPATH /home/go
ENV GOCACHE=off

ENV CGO_CXXFLAGS="-I/home/oiio-Release-$OIIO_VER/dist/linux64/include"
ENV CGO_LDFLAGS="-L/home/oiio-Release-$OIIO_VER/dist/linux64/lib"
ENV LD_LIBRARY_PATH="/home/oiio-Release-$OIIO_VER/dist/linux64/lib"

ADD . /home/go/src/$PKGNAME

CMD go test -v $PKGNAME
