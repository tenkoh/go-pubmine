package pubmine

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"golang.org/x/sync/semaphore"
)

const (
	npubPrefix = "npub1"
)

var (
	bech32 = regexp.MustCompile("^[02-9ac-hj-np-z]+$")
)

type ErrInitializeGenerator struct {
	cause     string
	userInput string
}

func (e ErrInitializeGenerator) Error() string {
	return fmt.Sprintf("Failed to initialize a keypair generator: %s: your input is %s", e.cause, e.userInput)
}

type ErrInterrupted struct{}

func (e ErrInterrupted) Error() string {
	return "Operation was interrupted"
}

// Generator is a type to configurate the operation setting.
type Generator struct {
	maxWorkers int64
	prefix     string
}

// Keypair is a set of a public key and a private (secret) key.
// The both keys are encoded following nip19.
type KeyPair struct {
	Public  string
	Private string
}

// NewGenerator is a factory function of Generator.
// The argument "prefix" must be bech32 format.
func NewGenerator(prefix string, maxWorkers int64) (*Generator, error) {
	if !bech32.MatchString(prefix) {
		return nil, ErrInitializeGenerator{"bad format prefix. 1, b, i, o are not allowed.", prefix}
	}
	g := Generator{
		prefix:     npubPrefix + prefix,
		maxWorkers: maxWorkers,
	}
	return &g, nil
}

// Mine tries to mine a keypair whose public key contains
// the specified prefix.
func (g *Generator) Mine(ctx context.Context) (*KeyPair, error) {
	sem := semaphore.NewWeighted(g.maxWorkers)
	ctx, cancel := context.WithCancel(ctx)
	ckp := make(chan *KeyPair)
	go func() {
		for {
			if err := sem.Acquire(ctx, 1); err != nil {
				return
			}
			go func() {
				defer sem.Release(1)
				kp, err := genKeyPair()
				if err != nil {
					return
				}
				if strings.HasPrefix(kp.Public, g.prefix) {
					ckp <- kp
				}
			}()
		}
	}()

	var kp *KeyPair
	var ok bool
	select {
	case kp, ok = <-ckp:
		// when a keypair found, stop the goroutines with cancel below
	case <-ctx.Done():
		// cancel from the parent context
		ok = false //this is equivalent to pass
	}
	// terminate producer
	cancel()
	// close channel after all running workers finish
	for {
		if sem.TryAcquire(g.maxWorkers) {
			break
		}
	}
	close(ckp)

	if !ok {
		return nil, ErrInterrupted{}
	}
	return kp, nil
}

// SimpleMine disables concurrent process.
// This method is expected to be used in single thread environment like WebAssembly.
func (g *Generator) SimpleMine(ctx context.Context) (*KeyPair, error) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
		kp, err := genKeyPair()
		if err != nil {
			continue
		}
		if strings.HasPrefix(kp.Public, g.prefix) {
			return kp, nil
		}
	}
	return nil, ErrInterrupted{}
}

func genKeyPair() (*KeyPair, error) {
	sk := nostr.GeneratePrivateKey()
	pk, err := nostr.GetPublicKey(sk)
	if err != nil {
		return nil, err
	}

	npk, err := nip19.EncodePublicKey(pk)
	if err != nil {
		return nil, err
	}
	nsk, err := nip19.EncodePrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return &KeyPair{Private: nsk, Public: npk}, nil
}
