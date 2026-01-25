// +build !no_ldflags

package oiio

// #cgo LDFLAGS: -L/usr/local/lib -lOpenImageIO -lOpenImageIO_Util
import "C"
