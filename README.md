# OpenImageIO bindings for Go

OpenImageIO is a library for reading and writing images, and a bunch of related classes,
utilities, and applications.  There is a particular emphasis on formats and functionality
used in professional, large-scale animation and visual effects work for film.
OpenImageIO is used extensively in animation and VFX studios all over the world, and is
also incorporated into several commercial products.

Compatibility
-------------

Tested against OpenImageIO 1.4.x - 1.5.x 
 
API Status
-----------

There is pretty decent exposure of the "Image*" APIs thus far, as well as the ColorConfig API. 
Because OIIO is a fairly large library, not every aspect of the APIs have been wrapped yet. It 
has mainly been driven by use-cases.

If you find something that you need is missing, feel free to submit a feature request, or better yet, 
fork and send a merge request :-)

Requirements
----------------------

* [OpenImageIO](https://github.com/OpenImageIO)
* [Boost (For ImageBufAlgo)](http://www.boost.org/)

Installation
------------

This package assumes that OpenImageIO/Boost is installed to the standard /usr/local location.

Default install:

    go get github.com/justinfx/openimageigo

If you have installed OpenImageIO to a custom location, you will need to tell CGO where to find the headers and libs:

    export CGO_CPPFLAGS="-I/path/to/include"
	export CGO_LDFLAGS="-L/path/to/lib"

Or just prefixing the install:

	CGO_CPPFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib" go get github.com/justinfx/openimageigo

Documentation
-------------

[http://godoc.org/github.com/justinfx/openimageigo](http://godoc.org/github.com/justinfx/openimageigo)


