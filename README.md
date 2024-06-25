# The Piper Programming Language

## Goal
Piper aims to be a general purpose functional language.

This repository contains the interpreter for the piper language specification.

## Setup

### Prerequisites
If you plan active development on the piper parser, you'll need to ensure that your `$GOBIN` environment variable as well as your `$PATH` environment variables are setup accordingly before the installation of pigeon.

A possible setup in your `.bashrc` or `.zshrc` can look like this

```bash
# Go related
export GOBIN=$(go env GOPATH)/bin
export GOPATH=$(go env GOPATH)

export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
```

This ensures that the pigeon binary gets installed into your `$GOPATH/bin` directory (usually `$HOME/bin`) and is also loaded into the `$PATH` environment variable.
After you restarted your terminal, the `pigeon -h` command should be available and ready to use.

**If you installed pigeon beforehand, you'll need to rerun the following command:**
```go
go install github.com/mna/pigeon@latest`
```
---
### Running the interpreter

For running the interpreter, you can use the following commands (in the root folder of the project)

```bash
# Regenerate the parser from the pigeon.peg file
pigeon -support-left-recursion -o internal/parser/parser.go internal/parser/pigeon.peg

# Run the interpreter
go run .

# Run the interpreter in repl mode
go run . -repl

# Regenerate parser & running the interpreter (useful for parser development)
pigeon -support-left-recursion -o internal/parser/parser.go internal/parser/pigeon.peg && go run .
```
