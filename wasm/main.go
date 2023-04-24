package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed binding/test.wasm
var icuWasm []byte

type URegularExpressionPtr uint32
type UCharPtr uint32
type UErrorCode int32
type CharPtr int32

type funcs struct {
	mod api.Module
}

// Bind
// u_strToUTF8(char *out, int bufflen, int * outsize, UChar *, int ucharlen, UErrorCode *)
// To
// u_strToUTF8(CharPtr, int, *int, UCharPtr, int, *UErrorCode)
func (a funcs) U_strToUTF8(ctx context.Context, buff CharPtr, bufflen int, outlen *int, str UCharPtr, strlen int, uerr *UErrorCode) {
	strToUTF8 := a.mod.ExportedFunction("u_strToUTF8_68")
	var uerrptr uint64 = 0
	if uerr != nil {
		uerrptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(uerrptr), uint32(*uerr))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(uerrptr))
			if !ok {
				panic("could not read uerr")
			}
			*uerr = UErrorCode(res)
			a.Free(ctx, uint32(uerrptr))
		}()
	}
	var outlenptr uint64 = 0
	if outlen != nil {
		outlenptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(outlenptr), uint32(*outlen))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(outlenptr))
			if !ok {
				panic("could not read uerr")
			}
			*outlen = int(res)
			a.Free(ctx, uint32(outlenptr))
		}()
	}
	_, err := strToUTF8.Call(ctx, uint64(buff), uint64(bufflen), outlenptr, uint64(str), uint64(strlen), uerrptr)
	if err != nil {
		panic(err)
	}
}

// Bind
// u_strFromUTF8(UChar *out, int bufflen, int * outsize, char *, int strlen, UErrorCode *)
// To
// u_strFromUTF8(UCharPtr, int, *int, CharPtr, int, *UErrorCode)
func (a funcs) U_strFromUTF8(ctx context.Context, buff UCharPtr, bufflen int, outlen *int, str CharPtr, strlen int, uerr *UErrorCode) {
	strFromUTF8 := a.mod.ExportedFunction("u_strFromUTF8_68")
	var uerrptr uint64 = 0
	if uerr != nil {
		uerrptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(uerrptr), uint32(*uerr))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(uerrptr))
			if !ok {
				panic("could not read uerr")
			}
			*uerr = UErrorCode(res)
			a.Free(ctx, uint32(uerrptr))
		}()
	}
	var outlenptr uint64 = 0
	if outlen != nil {
		outlenptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(outlenptr), uint32(*outlen))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(outlenptr))
			if !ok {
				panic("could not read uerr")
			}
			*outlen = int(res)
			a.Free(ctx, uint32(outlenptr))
		}()
	}
	_, err := strFromUTF8.Call(ctx, uint64(buff), uint64(bufflen), outlenptr, uint64(str), uint64(strlen), uerrptr)
	if err != nil {
		panic(err)
	}
}

// Bind
// uregex_findNext(URegularExpression *regex, UErrorCode *)
// To
// uregex_findNext(URegularExpressionPtr, *UErrorCode)
func (a funcs) Uregex_findNext(ctx context.Context, regex URegularExpressionPtr, uerr *UErrorCode) bool {
	findnext := a.mod.ExportedFunction("uregex_findNext_68")
	var uerrptr uint64 = 0
	if uerr != nil {
		uerrptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(uerrptr), uint32(*uerr))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(uerrptr))
			if !ok {
				panic("could not read uerr")
			}
			*uerr = UErrorCode(res)
			a.Free(ctx, uint32(uerrptr))
		}()
	}
	res, err := findnext.Call(ctx, uint64(regex), uerrptr)
	if err != nil {
		panic(err)
	}
	return res[0] != 0
}

func (a funcs) Uregex_start(ctx context.Context, regex URegularExpressionPtr, group int, uerr *UErrorCode) uint32 {
	start := a.mod.ExportedFunction("uregex_start_68")
	var uerrptr uint64 = 0
	if uerr != nil {
		uerrptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(uerrptr), uint32(*uerr))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(uerrptr))
			if !ok {
				panic("could not read uerr")
			}
			*uerr = UErrorCode(res)
			a.Free(ctx, uint32(uerrptr))
		}()
	}
	res, err := start.Call(ctx, uint64(regex), uint64(group), uerrptr)
	if err != nil {
		panic(err)
	}
	return uint32(res[0])
}

func (a funcs) Uregex_end(ctx context.Context, regex URegularExpressionPtr, group int, uerr *UErrorCode) uint32 {
	end := a.mod.ExportedFunction("uregex_end_68")
	var uerrptr uint64 = 0
	if uerr != nil {
		uerrptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(uerrptr), uint32(*uerr))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(uerrptr))
			if !ok {
				panic("could not read uerr")
			}
			*uerr = UErrorCode(res)
			a.Free(ctx, uint32(uerrptr))
		}()
	}
	res, err := end.Call(ctx, uint64(regex), uint64(group), uerrptr)
	if err != nil {
		panic(err)
	}
	return uint32(res[0])
}

func (a funcs) Uregex_setText(ctx context.Context, p URegularExpressionPtr, str UCharPtr, strlen int, uerr *UErrorCode) {
	setText := a.mod.ExportedFunction("uregex_setText_68")
	var uerrptr uint64 = 0
	if uerr != nil {
		uerrptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(uerrptr), uint32(*uerr))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(uerrptr))
			if !ok {
				panic("could not read uerr")
			}
			*uerr = UErrorCode(res)
			a.Free(ctx, uint32(uerrptr))
		}()
	}
	_, err := setText.Call(ctx, uint64(p), uint64(str), uint64(strlen), uerrptr)
	if err != nil {
		panic(err)
	}
}

func (a funcs) Uregex_open(ctx context.Context, str UCharPtr, strlen int, flags uint32, uerr *UErrorCode) URegularExpressionPtr {
	open := a.mod.ExportedFunction("uregex_open_68")
	var uerrptr uint64 = 0
	if uerr != nil {
		uerrptr = uint64(a.Malloc(ctx, 4))
		a.mod.Memory().WriteUint32Le(uint32(uerrptr), uint32(*uerr))
		defer func() {
			res, ok := a.mod.Memory().ReadUint32Le(uint32(uerrptr))
			if !ok {
				panic("could not read uerr")
			}
			*uerr = UErrorCode(res)
			a.Free(ctx, uint32(uerrptr))
		}()
	}
	res, err := open.Call(ctx, uint64(str), uint64(strlen), uint64(flags), uint64(0), uerrptr)
	if err != nil {
		panic(err)
	}
	return URegularExpressionPtr(res[0])
}

func (a funcs) Uregex_close(ctx context.Context, p URegularExpressionPtr) {
	close := a.mod.ExportedFunction("uregex_close_68")
	_, err := close.Call(ctx, uint64(p))
	if err != nil {
		panic(err)
	}
}

func (a funcs) Malloc(ctx context.Context, sz uint32) uint32 {
	malloc := a.mod.ExportedFunction("malloc")
	res, err := malloc.Call(ctx, uint64(sz))
	if err != nil {
		panic(err)
	}
	return uint32(res[0])
}

func (a funcs) Free(ctx context.Context, ptr uint32) {
	free := a.mod.ExportedFunction("free")
	_, err := free.Call(ctx, uint64(ptr))
	if err != nil {
		panic(err)
	}
}

func NewRuntime(ctx context.Context) wazero.Runtime {
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	envBuilder := r.NewHostModuleBuilder("env")
	noop_two := func(int32, int32) int32 { return -1 }
	noop_four := func(int32, int32, int32, int32) int32 { return -1 }
	envBuilder.NewFunctionBuilder().WithFunc(noop_two).Export("__syscall_stat64")
	envBuilder.NewFunctionBuilder().WithFunc(noop_two).Export("__syscall_lstat64")
	envBuilder.NewFunctionBuilder().WithFunc(noop_four).Export("__syscall_newfstatat")
	_, err := envBuilder.Instantiate(ctx)
	if err != nil {
		panic(err)
	}

	return r
}

func LoadICUModule(ctx context.Context, runtime wazero.Runtime) api.Module {
	mod, err := runtime.Instantiate(ctx, icuWasm)
	if err != nil {
		panic(err)
	}
	return mod
}

func RunRegexMatching(ctx context.Context, mod api.Module, tomatch, re string) {
	var f funcs = funcs{mod}

	str := f.Malloc(ctx, uint32(len(tomatch)))
	if !mod.Memory().Write(uint32(str), []byte(tomatch)) {
		panic("did not copy string")
	}
	defer f.Free(ctx, uint32(str))

	restr := f.Malloc(ctx, uint32(len(re)))
	if !mod.Memory().Write(uint32(restr), []byte(re)) {
		panic("did not copy string")
	}
	defer f.Free(ctx, uint32(restr))

	// str -> ustr

	var ustrlen int
	var uerr UErrorCode
	f.U_strFromUTF8(ctx, 0, 0, &ustrlen, CharPtr(str), len(tomatch), &uerr)
	if uerr != 15 {
		panic(fmt.Sprintf("unexpected uerr: %d", uerr))
	}
	ustr := UCharPtr(f.Malloc(ctx, uint32(ustrlen * 2)))
	defer f.Free(ctx, uint32(ustr))
	uerr = 0
	f.U_strFromUTF8(ctx, ustr, ustrlen, nil, CharPtr(str), len(tomatch), &uerr)
	if uerr > 0 {
		panic(fmt.Sprintf("unexpected uerr: %d", uerr))
	}

	// restr -> reustr

	var urestrlen int
	uerr = 0
	f.U_strFromUTF8(ctx, 0, 0, &urestrlen, CharPtr(restr), len(re), &uerr)
	if uerr != 15 {
		panic(fmt.Sprintf("unexpected uerr: %d", uerr))
	}
	urestr := UCharPtr(f.Malloc(ctx, uint32(urestrlen * 2)))
	defer f.Free(ctx, uint32(urestr))
	uerr = 0
	f.U_strFromUTF8(ctx, urestr, urestrlen, nil, CharPtr(restr), len(re), &uerr)
	if uerr > 0 {
		panic(fmt.Sprintf("unexpected uerr: %d", uerr))
	}

	// make a regex

	uerr = 0
	regex := f.Uregex_open(ctx, urestr, urestrlen, 0, &uerr)
	if uerr > 0 {
		panic(fmt.Sprintf("unexpected uerr: %d", uerr))
	}
	defer f.Uregex_close(ctx, regex)

	// set its text

	uerr = 0
	f.Uregex_setText(ctx, regex, ustr, ustrlen, &uerr)
	if uerr > 0 {
		panic(fmt.Sprintf("unexpected uerr: %d", uerr))
	}

	// print out the matches

	uerr = 0
	for f.Uregex_findNext(ctx, regex, &uerr) {
		if uerr != 0 {
			panic("unexpected uerr")
		}

		uerr = 0
		start := f.Uregex_start(ctx, regex, 0, &uerr)
		if uerr != 0 {
			panic("unexpected uerr")
		}

		uerr = 0
		end := f.Uregex_end(ctx, regex, 0, &uerr)
		if uerr != 0 {
			panic("unexpected uerr")
		}

		str := uchar_substr(ctx, f, ustr, ustrlen, start, end)
		fmt.Println("found match", str)

		uerr = 0
	}
}

func uchar_substr(ctx context.Context, f funcs, p UCharPtr, len int, start, end uint32) string {
	var uerr UErrorCode = 0
	var outlen int
	f.U_strToUTF8(ctx, 0, 0, &outlen, UCharPtr(uint64(p) + uint64(start) * 2), int(uint64(end) - uint64(start)), &uerr)
	if uerr != 15 {
		panic("unexpected uerr")
	}

	outstr := f.Malloc(ctx, uint32(outlen))
	defer f.Free(ctx, uint32(outstr))

	uerr = 0
	f.U_strToUTF8(ctx, CharPtr(outstr), outlen, nil, UCharPtr(uint64(p) + uint64(start) * 2), int(uint64(end) - uint64(start)), &uerr)
	if uerr > 0 {
		panic("unexpected uerr")
	}

	bs, ok := f.mod.Memory().Read(uint32(outstr), uint32(outlen))
	if !ok {
		panic("unexpected read out of bounds")
	}
	return string(bs)
}

func main() {
	ctx := context.Background()

	r := NewRuntime(ctx)
	defer r.Close(ctx)

	mod := LoadICUModule(ctx, r)

	RunRegexMatching(ctx, mod, "testing, one, two, three", "[a-z]")
}
