//go:build js && wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"context"
	"fmt"
	"runtime"
	"syscall/js"

	"github.com/tenkoh/go-prtkey"
)

var (
	maxWorkers = int64(runtime.NumCPU())
)

func generateKeypair(this js.Value, inputs []js.Value) any {
	fmt.Println("wasm called")
	// get elements
	doc := js.Global().Get("document")
	input := doc.Call("getElementById", "input")
	public := doc.Call("getElementById", "public")
	private := doc.Call("getElementById", "private")

	go func() {
		prefix := input.Get("value").String()
		g, err := prtkey.NewGenerator(prefix, maxWorkers)
		if err != nil {
			// error output into public form
			public.Set("value", "The specified prefix is bad format.")
			return
		}

		ctx := context.Background()
		kp, err := g.Mine(ctx)
		if err != nil {
			// error output into public form
			public.Set("value", "No keypairs found")
			return
		}

		public.Set("value", kp.Public)
		private.Set("value", kp.Private)

	}()
	return nil
}

func main() {
	fmt.Println("wasm loaded")
	js.Global().Set("generateKeypair", js.FuncOf(generateKeypair))
	select {}
}
