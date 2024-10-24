package main

import "fmt"

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var (
	ColorBlackNormal   = Color("\033[0;30m%s\033[0m")
	ColorBlackBright   = Color("\033[1;30m%s\033[0m")
	ColorRedNormal     = Color("\033[0;31m%s\033[0m")
	ColorRedBright     = Color("\033[1;31m%s\033[0m")
	ColorGreenNormal   = Color("\033[0;32m%s\033[0m")
	ColorGreenright    = Color("\033[1;32m%s\033[0m")
	ColorYellowNormal  = Color("\033[0;33m%s\033[0m")
	ColorYellowright   = Color("\033[1;33m%s\033[0m")
	ColorBlueNormal    = Color("\033[0;34m%s\033[0m")
	ColorBlueBright    = Color("\033[1;34m%s\033[0m")
	ColorMagentaNormal = Color("\033[0;35m%s\033[0m")
	ColorMagentaBright = Color("\033[1;35m%s\033[0m")
	ColorCyanNormal    = Color("\033[0;36m%s\033[0m")
	ColorCyanBright    = Color("\033[1;36m%s\033[0m")
	ColorWhiteNormal   = Color("\033[0;37m%s\033[0m")
	ColorWhiteBright   = Color("\033[1;37m%s\033[0m")
)
