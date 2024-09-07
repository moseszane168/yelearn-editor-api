//go:build linux

package rfid

/*
#cgo CFLAGS: -I ${SRCDIR}/lib/
#cgo LDFLAGS: -L${SRCDIR}/lib/ -lrfid_core
#cgo LDFLAGS: -L${SRCDIR}/lib/ -lrfid
#include "./lib/rfid.h"
*/
import "C"
