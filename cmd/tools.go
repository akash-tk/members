//go:build tools
// +build tools

// cmd contains binaries for this project
//
// For example, install mockgen and wire as below,
//
// go mod tidy &&  grep _ cmd/tools.go| cut -d ' ' -f 2 | xargs go install
package cmd

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/wire/cmd/wire"
)
