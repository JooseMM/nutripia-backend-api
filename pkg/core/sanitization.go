package core

import "regexp"

func HasHarmfulSymbols(input string) bool {
	pattern := `[<>'";{}()\[\]\\/]`
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}
