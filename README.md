# git

[![Go](https://github.com/dotWicho/git/workflows/Go/badge.svg?branch=master)](https://github.com/dotWicho/git)
[![Quality Report](https://goreportcard.com/badge/github.com/dotWicho/git)](https://goreportcard.com/badge/github.com/dotWicho/git)
[![GoDoc](https://godoc.org/github.com/dotWicho/git?status.svg)](https://pkg.go.dev/github.com/dotWicho/git?tab=doc)

## Library to manage Google go-github with more simplicity

## Getting started

- API documentation is available via [godoc](https://godoc.org/github.com/dotWicho/git).

## Installation

To install Git package, you need to install Go and set your Go workspace first.

1 - The first need [Go](https://golang.org/) installed (**version 1.13+ is required**).
Then you can use the below Go command to install Git

```bash
$ go get -u github.com/dotWicho/git
```

And then Import it in your code:

``` go
package main

import "github.com/dotWicho/git"
```
Or

2 - Use as module in you project (go.mod file):

``` go
module myclient

go 1.13

require (
	github.com/dotWicho/git v1.0.0
)
```

