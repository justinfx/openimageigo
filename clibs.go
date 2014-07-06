// +build !no_ldflags

package oiio

// #cgo LDFLAGS: -L/usr/local/lib -lOpenImageIO -lboost_thread-mt -lboost_system-mt
import "C"
