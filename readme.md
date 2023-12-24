# Sync Slice Library for Go

## Overview

The Sync Slice library provides a concurrency-safe implementation of a dynamic slice in Go (Golang). It is designed to be used in environments where multiple goroutines need to safely access or modify a slice.

## Features

- **Concurrency Safe**: All operations are safe for concurrent use by multiple goroutines.
- **Generic Implementation**: Uses Go generics, allowing it to be used with any type.
- **Common Slice Operations**: Supports operations like Append, Get, Set, and Remove.

## Installation

To use the Sync Slice library in your Go project, use the following command:

```bash
go get github.com/IIGabriel/sync-slice
```

## Usage

Here is a simple example of how to use the Sync Slice library:

```go
package main

import (
	"fmt"
	"github.com/IIGabriel/sync-slice/pkg"
)

func main() {
	s := syncslice.New[int]()

	s.Append(1)
	s.Append(2)

	value, ok := s.Get(1)
	if ok {
		fmt.Println("Value at index 1:", value)
	}

	s.Range(func(index int, value int) bool {
		fmt.Printf("Index: %d, Value: %dn", index, value)
		return true
	})
}

```

## Contributing

Contributions to the Sync Slice library are welcome. Please feel free to submit issues and pull requests.
