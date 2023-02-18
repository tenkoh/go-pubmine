package prtkey

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

type ErrExceedMaxTrial struct {
	maxTrial int
}

func (e ErrExceedMaxTrial) Error() string {
	return fmt.Sprintf("No keypairs found in %d trials", e.maxTrial)
}

// Generator is a type to configurate the operation setting.
type Generator struct {
	maxTrial   int
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
func NewGenerator(prefix string, maxTrial int, maxWorkers int64) (*Generator, error) {
	if !bech32.MatchString(prefix) {
		return nil, ErrInitializeGenerator{"bad format prefix", prefix}
	}
	g := Generator{
		prefix:     npubPrefix + prefix,
		maxTrial:   maxTrial,
		maxWorkers: maxWorkers,
	}
	return &g, nil
}

// Mine tries to mine a keypair whose public key contains
// the specified prefix. If no keypairs are found in
// maxTrial, ErrExceedMaxTrial is returned.
func (g *Generator) Mine(ctx context.Context) (*KeyPair, error) {
	sem := semaphore.NewWeighted(g.maxWorkers)
	ctx, cancel := context.WithCancel(ctx)
	ckp := make(chan *KeyPair)
	go func() {
		for i := 0; i < g.maxTrial; i++ {
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
		// wait until all goroutines finish, then cancel
		if err := sem.Acquire(ctx, g.maxWorkers); err != nil {
			return
		}
		cancel() // send ctx.Done() to the select block below
	}()

	var kp *KeyPair
	var ok bool
	select {
	case kp, ok = <-ckp:
		cancel()
	case <-ctx.Done():
		ok = false //this is equivalent to pass
	}
	// close channel after all running workers finish
	// if this waiting process takes a long time,
	// consider to isolate it in another goroutine.
	sem.Acquire(ctx, g.maxWorkers)
	close(ckp)

	if !ok {
		return nil, ErrExceedMaxTrial{g.maxTrial}
	}
	return kp, nil
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
