package filelib

import "regexp"

var reSpecials = regexp.MustCompile(`\\|\/|\?|\!`)

var reBackslash = regexp.MustCompile(`\\`)
var reExt = regexp.MustCompile(`\..*`)
var rePoint = regexp.MustCompile(`^\.`)
