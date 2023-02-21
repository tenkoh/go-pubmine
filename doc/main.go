//go:build js && wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/tenkoh/go-pubmine"
)

// No concurrency in wasm.
// Deletegate it to JavaScript with WebWorker.
const maxWorkers = 1

func mine(this js.Value, args []js.Value) any {
	// args contains at least one element even if no arguments passed from JS.
	// In the case above, the argument is undefined.
	if args[0].IsUndefined() {
		return map[string]any{"error": "argument is required"}
	}

	prefix := args[0].String()
	g, err := pubmine.NewGenerator(prefix, maxWorkers)
	if err != nil {
		return map[string]any{"error": fmt.Sprintf("failed to initialize a miner: %s", err.Error())}
	}
	kp, err := g.SimpleMine(context.Background())
	if err != nil {
		return map[string]any{"error": fmt.Sprintf("failed to mine a keypair: %s", err.Error())}
	}

	return map[string]any{"public": kp.Public, "private": kp.Private}
}

func main() {
	js.Global().Set("mine", js.FuncOf(mine))
	select {}
}
