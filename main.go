package main

import (
	"fmt"
	"github.com/reltuk/icu-play/pkg/icu"
)

func main() {
	str, err := icu.NewUCharString("hello, world!")
	if err != nil {
		panic(err)
	}
	defer str.Free()

	pattern, err := icu.NewUCharString("[a-z]")
	if err != nil {
		panic(err)
	}
	defer pattern.Free()

	regex, err := icu.NewURegex(pattern, 0)
	if err != nil {
		panic(err)
	}
	defer regex.Free()

	err = regex.SetText(str)
	if err != nil {
		panic(err)
	}

	for {
		ok, err := regex.FindNext()
		if err != nil {
			panic(err)
		}
		if !ok {
			break
		}

		s, e, err := regex.CurrentMatch()
		if err != nil {
			panic(err)
		}

		ss, err := str.SubStrView(s, e).GoString()
		if err != nil {
			panic(err)
		}
		fmt.Println("found match:", ss)
	}
}
