package prtkey_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/tenkoh/go-prtkey"
)

func TestNewGenerator(t *testing.T) {
	type args struct {
		prefix     string
		maxWorkers int64
	}
	type test struct {
		name      string
		args      args
		wantError bool
	}
	tests := []test{
		{"success case", args{"n0str", 1}, false},
		{"fail blank", args{"", 1}, true},
		{"fail bad character", args{"nostr", 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.args
			_, err := prtkey.NewGenerator(a.prefix, a.maxWorkers)
			hasError := err != nil
			if hasError != tt.wantError {
				t.Errorf("want %t, got %t", tt.wantError, hasError)
			}
		})
	}
}

func TestMine(t *testing.T) {
	// just test the logic is not broken
	g, err := prtkey.NewGenerator("n0st", int64(runtime.NumCPU()))
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

func TestMineWithCancel(t *testing.T) {
	// just test the logic is not broken
	g, err := prtkey.NewGenerator("n0st", int64(runtime.NumCPU()))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	kp, err := g.Mine(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", kp)
}
