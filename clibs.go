// +build !no_ldflags

package oiio

// #cgo LDFLAGS: -L/usr/local/lib -lOpenImageIO -lOpenImageIO_Util -lboost_thread -lboost_system
import "C"
