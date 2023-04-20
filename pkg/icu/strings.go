package icu

// TODO: We link with -lc++ below, which is libc++ clang specific. All of this
// needs build flags, etc.

// #cgo CPPFLAGS: -I${SRCDIR}/../../third_party/icu_install/include
// #cgo LDFLAGS: -L${SRCDIR}/../../third_party/icu_install/lib -licui18n -licuuc -licudata -lc++
// #include <stdlib.h>
// #include "unicode/ustring.h"
// #include "unicode/utext.h"
// #include "unicode/uregex.h"
import "C"
import "unsafe"
import "fmt"

// A UCharString owns the backing array |ptr|, which is allocated on the C
// heap.
type UCharString struct {
	ptr *C.UChar

	// length in UChars
	len C.int32_t
}

func (s UCharString) Free() {
	C.free(unsafe.Pointer(s.ptr))
}

func (s UCharString) GoString() (string, error) {
	return UCharStringView{s.ptr, s.len}.GoString()
}

func (s UCharString) SubStrView(start, end int) UCharStringView {
	return UCharStringView{
		ptr: (*C.UChar)(unsafe.Add(unsafe.Pointer(s.ptr), start*C.sizeof_UChar)),
		len: C.int32_t(end - start),
	}
}

func NewUCharString(s string) (UCharString, error) {
	var uerr C.UErrorCode = C.U_ZERO_ERROR

	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var clen C.int = C.int(len(s))

	var len C.int32_t
	C.u_strFromUTF8(nil, 0, &len, cstr, clen, &uerr)
	if uerr != C.U_BUFFER_OVERFLOW_ERROR {
		return UCharString{}, fmt.Errorf("unexpected error preflighting string length: %d", uerr)
	}

	uerr = C.U_ZERO_ERROR
	var uc *C.UChar = (*C.UChar)(C.malloc(C.ulong(len * C.sizeof_UChar)))
	C.u_strFromUTF8(uc, len, nil, cstr, clen, &uerr)
	if uerr > C.U_ZERO_ERROR {
		return UCharString{}, fmt.Errorf("unexpected error converting to string: %d", uerr)
	}

	return UCharString{uc, len}, nil
}

type UCharStringView struct {
	// The start of the string view. The view does not own this pointer.
	ptr *C.UChar

	// length in UChars
	len C.int32_t
}

func (s UCharStringView) GoString() (string, error) {
	var uerr C.UErrorCode = C.U_ZERO_ERROR

	var len C.int32_t
	C.u_strToUTF8(nil, 0, &len, s.ptr, s.len, &uerr)
	if uerr != C.U_BUFFER_OVERFLOW_ERROR {
		return "", fmt.Errorf("unexpected error preflighting string length: %d", uerr)
	}

	uerr = C.U_ZERO_ERROR
	var cstr *C.char = (*C.char)(C.malloc(C.ulong(len)))
	defer C.free(unsafe.Pointer(cstr))
	C.u_strToUTF8(cstr, len, nil, s.ptr, s.len, &uerr)
	if uerr > C.U_ZERO_ERROR {
		return "", fmt.Errorf("unexpected error converting to string: %d", uerr)
	}

	return C.GoStringN(cstr, len), nil
}

type URegex struct {
	ptr *C.URegularExpression
}

func (r URegex) Free() {
	C.uregex_close(r.ptr)
}

func NewURegex(str UCharString, flags uint32) (URegex, error) {
	var uerr C.UErrorCode = C.U_ZERO_ERROR
	ptr := C.uregex_open(str.ptr, str.len, C.uint32_t(flags), nil, &uerr)
	if uerr > C.U_ZERO_ERROR {
		return URegex{}, fmt.Errorf("unexpected error in uregex_open: %d", uerr)
	}
	return URegex{ptr}, nil
}

func (r URegex) SetText(str UCharString) error {
	var uerr C.UErrorCode = C.U_ZERO_ERROR
	C.uregex_setText(r.ptr, str.ptr, str.len, &uerr)
	if uerr > C.U_ZERO_ERROR {
		return fmt.Errorf("unexpected error in uregex_setText: %d", uerr)
	}
	return nil
}

func (r URegex) FindNext() (bool, error) {
	var uerr C.UErrorCode = C.U_ZERO_ERROR
	res := C.uregex_findNext(r.ptr, &uerr)
	if uerr > C.U_ZERO_ERROR {
		return false, fmt.Errorf("unexpected error in uregex_findNext: %d", uerr)
	}
	return res != 0, nil
}

func (r URegex) CurrentMatch() (start, end int, err error) {
	var uerr C.UErrorCode = C.U_ZERO_ERROR
	cstart := C.uregex_start(r.ptr, 0, &uerr)
	if uerr > C.U_ZERO_ERROR {
		return 0, 0, fmt.Errorf("unexpected error in uregex_start: %d", uerr)
	}
	uerr = C.U_ZERO_ERROR
	cend := C.uregex_end(r.ptr, 0, &uerr)
	if uerr > C.U_ZERO_ERROR {
		return 0, 0, fmt.Errorf("unexpected error in uregex_end: %d", uerr)
	}
	return int(cstart), int(cend), nil
}
