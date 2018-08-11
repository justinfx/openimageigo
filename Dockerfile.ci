FROM ubuntu:xenial

ENV OIIO_VER 1.6.18

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
    golang-1.10-go


RUN cd /home \
    && wget --no-check-certificate https://github.com/OpenImageIO/oiio/archive/Release-$OIIO_VER.tar.gz \
	&& tar zxf Release-$OIIO_VER.tar.gz \
	&& rm -f Release-$OIIO_VER.tar.gz

RUN cd /home/oiio-Release-$OIIO_VER \
    && make USE_OPENGL=0 USE_QT=0 USE_PYTHON=0 USE_PYTHON3=0 OIIO_BUILD_TOOLS=0 OIIO_BUILD_TESTS=0 VERBOSE=1

ENV GOPATH /home/go
ENV PKGNAME github.com/justinfx/openimageigo

ENV PATH /usr/lib/go-1.10/bin:$PATH
ENV CGO_CXXFLAGS="-I/home/oiio-Release-${OIIO_VER}/dist/linux64/include"
ENV CGO_LDFLAGS="-L/home/oiio-Release-${OIIO_VER}/dist/linux64/lib"
ENV LD_LIBRARY_PATH="/home/oiio-Release-${OIIO_VER}/dist/linux64/lib"

WORKDIR ${GOPATH}/src/${PKGNAME}

CMD go test -count 1 -v $PKGNAME
