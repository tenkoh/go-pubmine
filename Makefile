ifeq ($(shell uname -s),Darwin)
	BROWSER = open
else
	BROWSER = xdg-open
endif

.PHONY: all clean serve

all: main.wasm serve

%.wasm: %.go
	GOOS=js GOARCH=wasm go generate
	GOOS=js GOARCH=wasm go build -o "$@" "$<"

serve:
	$(BROWSER) 'http://localhost:5000'
	serve || (go get -v github.com/mattn/serve && serve)

clean:
	rm -f *.wasm
