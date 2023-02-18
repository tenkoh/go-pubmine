package prtkey_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/tenkoh/go-prtkey"
)

func TestNewGenerator(t *testing.T) {
	type args struct {
		prefix     string
		maxTrial   int
		maxWorkers int64
	}
	type test struct {
		name      string
		args      args
		wantError bool
	}
	tests := []test{
		{"success case", args{"n0str", 1, 1}, false},
		{"fail blank", args{"", 1, 1}, true},
		{"fail bad character", args{"nostr", 1, 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.args
			_, err := prtkey.NewGenerator(a.prefix, a.maxTrial, a.maxWorkers)
			hasError := err != nil
			if hasError != tt.wantError {
				t.Errorf("want %t, got %t", tt.wantError, hasError)
			}
		})
	}
}

func TestMine(t *testing.T) {
	g, err := prtkey.NewGenerator("n0st", 20000000, int64(runtime.NumCPU()))
	if err != nil {
		t.Fatal(err)
	}
	kp, err := g.Mine(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", kp)
}
