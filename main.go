//go:build js && wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"syscall/js"
)

func main() {
	doc := js.Global().Get("document")
	input := doc.Call("getElementById", "input")
	prefix := input.Get("value").String()

	public := doc.Call("getElementById", "public")
	public.Set("value", prefix)
	private := doc.Call("getElementById", "private")
	private.Set("value", prefix)
}
