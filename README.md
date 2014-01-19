# OpenImageIO bindings for Go

Currently in development (incomplete API)

Requirements
----------------------

* [OpenImageIO](openimageio.org)


Installation
------------

This package assumes that OpenImageIO is installed to the standard /usr/local location. 

Default install:

    go get github.com/justinfx/opencolorigo

If you have installed OpenImageIO to a custom location, you will need to tell CGO where to find the headers and libs:

	export CGO_CFLAGS="-I/path/to/include" 
	export CGO_LDFLAGS="-L/path/to/lib"

Or just prefixing the install:

	CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib" go get github.com/justinfx/opencolorigo

Documentation
-------------

[http://godoc.org/github.com/justinfx/openimageigo](http://godoc.org/github.com/justinfx/openimageigo)


