//go:build mage
// +build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/mg"
)

// BuildAll is copied from Grafana's SDK (build package), because the Arrow
// dependency doesn't support Linux ARM in the build matrix. I've removed it
// from the listing here.
func BuildIt() {
	b := build.Build{}
	mg.Deps(b.Linux, b.Windows, b.Darwin, b.DarwinARM64, b.LinuxARM64)
}

var Default = BuildIt
