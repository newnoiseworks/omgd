package utils

import "embed"

//go:embed static/*
var StaticFiles embed.FS
