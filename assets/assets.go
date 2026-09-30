// Package assets embeds the decoded native artwork, tables and YM soundtrack.
package assets

import "embed"

//go:embed original/*
var Files embed.FS
