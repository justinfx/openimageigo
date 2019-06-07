# OpenImageIO bindings for Go

[![GoDoc](https://godoc.org/github.com/justinfx/openimageigo?status.svg)](https://godoc.org/github.com/justinfx/openimageigo) [![Travis Build Status](https://api.travis-ci.org/justinfx/openimageigo.svg)](https://travis-ci.org/justinfx/openimageigo)

OpenImageIO is a library for reading and writing images, and a bunch of related classes,
utilities, and applications.  There is a particular emphasis on formats and functionality
used in professional, large-scale animation and visual effects work for film.
OpenImageIO is used extensively in animation and VFX studios all over the world, and is
also incorporated into several commercial products.

Motivation
----------

While there are other image processing bindings available, OpenImageIO is a common image processing solution to the Visual Effects industry, with specific support for concepts like EXR, deep compositing, OpenColorIO support, textures, and subimages. It isn't neccessarily the fastest solution, but it is comprehensive, and useful to VFX pipelines.

Compatibility
-------------

Tested against OpenImageIO 1.4.x - 1.6.x 

Support for >= 1.7.x is in progress.
 
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
