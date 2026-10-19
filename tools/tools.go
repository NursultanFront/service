// Package tools pins the development tools of this repository.
//
// The tools are declared with "tool" directives in tools/go.mod (Go 1.24+), so
// their versions are fixed in tools/go.sum instead of relying on whatever is
// installed globally. Build all of them into ../bin with:
//
//	make tools
//
// Add a tool:     go -C tools get -tool <module/cmd/path>@<version>
// Upgrade tools:  make tools-upgrade
//
// NOTE(AI): created in migration step (e).
package tools

//go:generate go build -o ../bin/ tool
