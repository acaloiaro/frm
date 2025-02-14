package static

import "embed"

//go:embed js/* css/* img/*
var Assets embed.FS
