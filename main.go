//go:build js && wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"context"
	"errors"
	"regexp"
	"runtime"
	"syscall/js"

	"golang.org/x/sync/semaphore"
)

var (
	maxTrial   = 1000000
	maxWorkers = int64(runtime.NumCPU())
	bech32     = regexp.MustCompile("^[02-9ac-hj-np-z]+$")
)

type keyPair struct {
	Public  string
	Private string
}

func mineKeypair(s string) (*keyPair, error) {
	sem := semaphore.NewWeighted(maxWorkers)
	ctx, cancel := context.WithCancel(context.Background())
	ckp := make(chan keyPair)
	go func() {
		for i := 0; i < maxTrial; i++ {
			if err := sem.Acquire(ctx, 1); err != nil {
				break
			}
			go func() {
				defer sem.Release(1)
				kp, err := genKeyPair()
				if err != nil {
					return
				}
				if bech32.MatchString(kp.Public) {
					ckp <- kp
				}
			}()
		}
		// wait until all goroutines finish, then cancel
		if err := sem.Acquire(ctx, maxWorkers); err != nil {
			return
		}
		cancel()
	}()

	var kp keyPair
	var ok bool
	select {
	case kp, ok = <-ckp:
		cancel()
	case <-ctx.Done():
		ok = false //this is equivalent to PASS
	}
	// close channel after all running workers finish
	// if this waiting process takes a long time,
	// consider to isolate it in another goroutine.
	sem.Acquire(ctx, maxWorkers)
	close(ckp)

	if !ok {
		return nil, errors.New("could not find a key pair")
	}
	return &kp, nil

}

func genKeyPair() (keyPair, error) {
	return keyPair{}, nil
}

func main() {
	doc := js.Global().Get("document")
	input := doc.Call("getElementById", "input")
	prefix := input.Get("value").String()

	public := doc.Call("getElementById", "public")
	public.Set("value", prefix)
	private := doc.Call("getElementById", "private")
	private.Set("value", prefix)
}
