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
	$(BROWSER) 'http://localhost:8080'
	serve -a localhost:8080 || (go install github.com/mattn/serve@latest && serve -a localhost:8080)

clean:
	rm -f *.wasm
