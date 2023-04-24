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

type ModuleString struct {
	ptr uint64
	len uint64
	mod api.Module
}

func (s ModuleString) Close(ctx context.Context) {
	free := s.mod.ExportedFunction("free")
	free.Call(ctx, s.ptr)
}

func NewModuleString(ctx context.Context, mod api.Module, s string) ModuleString {
	malloc := mod.ExportedFunction("malloc")
	l := uint64(len(s))
	res, err := malloc.Call(ctx, l)
	if err != nil {
		panic(err)
	}
	ptr := res[0]
	if !mod.Memory().Write(uint32(ptr), []byte(s)) {
		panic("did not copy string")
	}
	return ModuleString{
		ptr,
		l,
		mod,
	}
}

type UCharString struct {
	ptr uint64
	mod api.Module
}

func (s UCharString) Close(ctx context.Context) {
	ucharstring_free := s.mod.ExportedFunction("icu_ucharstring_free")
	ucharstring_free.Call(ctx, s.ptr)
}

func NewUCharString(ctx context.Context, mod api.Module, s string) UCharString {
	str := NewModuleString(ctx, mod, s)
	defer str.Close(ctx)

	ucharstring_fromUTF8 := mod.ExportedFunction("icu_ucharstring_fromUTF8")
	res, err := ucharstring_fromUTF8.Call(ctx, str.ptr, str.len)
	if err != nil {
		panic(err)
	}
	return UCharString{res[0], mod}
}

type IntPtr struct {
	ptr uint64
	mod api.Module
}

func (i IntPtr) Close(ctx context.Context) {
	free := i.mod.ExportedFunction("free")
	free.Call(ctx, i.ptr)
}

func NewIntPtr(ctx context.Context, mod api.Module) IntPtr {
	malloc := mod.ExportedFunction("malloc")
	res, err := malloc.Call(ctx, 4)
	if err != nil {
		panic(err)
	}
	return IntPtr{
		res[0],
		mod,
	}
}

func (i IntPtr) Value() uint32 {
	v, ok := i.mod.Memory().ReadUint32Le(uint32(i.ptr))
	if !ok {
		panic("failed to read intptr")
	}
	return v
}

func (s UCharString) Substr(ctx context.Context, start, end uint64) string {
	len := NewIntPtr(ctx, s.mod)
	defer len.Close(ctx)

	free := s.mod.ExportedFunction("free")
	ucharstring_substr_toUTF8 := s.mod.ExportedFunction("icu_ucharstring_substr_toUTF8")

	// Then we call toUTF8.
	res, err := ucharstring_substr_toUTF8.Call(ctx, s.ptr, start, end, len.ptr)
	if err != nil {
		panic(err)
	}
	charPtr := res[0]
	defer free.Call(ctx, charPtr)

	// And we do some machinations to read the result.
	charSlice, ok := s.mod.Memory().Read(uint32(charPtr), len.Value())
	if !ok {
		panic("failed to read char slice")
	}
	return string(charSlice)
}

func ReadUCharString(ctx context.Context, mod api.Module, ucharPtr uint64) string {
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")
	ucharstring_toUTF8 := mod.ExportedFunction("icu_ucharstring_toUTF8")

	res, err := malloc.Call(ctx, 4)
	if err != nil {
		panic(err)
	}
	lenPtr := res[0]
	defer free.Call(ctx, lenPtr)

	// Then we call toUTF8.
	res, err = ucharstring_toUTF8.Call(ctx, ucharPtr, lenPtr)
	if err != nil {
		panic(err)
	}
	charPtr := res[0]
	defer free.Call(ctx, charPtr)

	// And we do some machinations to read the result.
	charLen, ok := mod.Memory().ReadUint32Le(uint32(lenPtr))
	if !ok {
		panic("failed to read len")
	}
	charSlice, ok := mod.Memory().Read(uint32(charPtr), charLen)
	if !ok {
		panic("failed to read char slice")
	}
	return string(charSlice)
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

func main() {
	ctx := context.Background()

	r := NewRuntime(ctx)
	defer r.Close(ctx)

	mod := LoadICUModule(ctx, r)

	str := NewUCharString(ctx, mod, "Hello, world!")
	defer str.Close(ctx)

	fmt.Println("read back", ReadUCharString(ctx, mod, str.ptr))

	fmt.Println("read back", str.Substr(ctx, 1, 5))

	// Now we are going to compile a regex, set its target text, run findNext and print out the results, and then clean up.

	regex_open := mod.ExportedFunction("icu_uregex_open")
	regex_close := mod.ExportedFunction("uregex_close_68")
	regex_findNext := mod.ExportedFunction("icu_uregex_findNext")
	regex_start := mod.ExportedFunction("icu_uregex_start")
	regex_end := mod.ExportedFunction("icu_uregex_end")
	regex_setText := mod.ExportedFunction("icu_uregex_setText")

	regexStr := NewUCharString(ctx, mod, "[a-z]")
	defer regexStr.Close(ctx)

	fmt.Println(str.ptr)
	fmt.Println(regexStr.ptr)

	res, err := regex_open.Call(ctx, regexStr.ptr, 0)
	if err != nil {
		panic(err)
	}
	regex := res[0]
	defer regex_close.Call(ctx, regex)

	_, err = regex_setText.Call(ctx, regex, str.ptr)
	if err != nil {
		panic(err)
	}

	for {
		res, err = regex_findNext.Call(ctx, regex)
		if err != nil {
			panic(err)
		}
		if res[0] == 0 {
			break
		}
		res, err = regex_start.Call(ctx, regex, 0)
		if err != nil {
			panic(err)
		}
		start := res[0]
		res, err = regex_end.Call(ctx, regex, 0)
		if err != nil {
			panic(err)
		}
		end := res[0]

		fmt.Println("found match", str.Substr(ctx, start, end))
	}

	fmt.Println("finished matching")
}
