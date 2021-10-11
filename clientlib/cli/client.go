package cli

import (
	"github.com/fatih/color"
)

var cCy *color.Color
var cGr *color.Color
var cRe *color.Color

func init() {
	cCy = color.New(color.FgCyan)
	cGr = color.New(color.FgGreen)
	cRe = color.New(color.FgRed)
}
