# go-pubmine
Multithreading key pair generator which gives pretty (vanity) public keys containing your favorite words.

## Overview
This repository contains three contents.

- Go's package (on root)
- cli tool (on cmd/pubmine)
- web app source (hosting on https://tenkoh.github.io/go-pubmine)

## Feature
- **It's safe!** : All processed are done on your local computer.
- **It's fast!** : Multithreading is supported.

## About cli

### Install
Now, only `go install` is supported.

```
go install github.com/tenkoh/go-pubmine/cmd/pubmine@latest
```

### Usage
Just enter the command below.

```
pubmine {prefix you want to use}
```

You can get a keypair. **Private key is shown on the terminal. So be careful.**

```
Public key:
npub1~
Private key:
nsec1~
```

## About web app

### Usage
1. Visit https://tenkoh.github.io/go-pubmine
2. Enter a prefix into the input form, then hit the RUN button.

## License
MIT

## Author
tenkoh


