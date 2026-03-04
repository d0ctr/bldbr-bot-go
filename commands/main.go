package commands

import (
	"strings"
	"regexp"
)

var wordsRe = regexp.MustCompile(` (\w+)`)

func parseArgs(text string, limit uint) map[uint]string {
	matches := wordsRe.FindAllStringSubmatch(text, -1)

	words := make(map[uint]string)
	for i, match := range matches {
		words[uint(i)] = match[1]
	}

	if (limit > 0 && limit <= uint(len(words))) {
		args := make(map[uint]string, limit)

		i := uint(0);
		for ; i < limit - 1; i++ {
			args[i] = words[i]
		}

		lastArgSlice := make([]string, len(words) - int(limit) + 1)
		for ; i < uint(len(words)); i++ {
			lastArgSlice[i - limit + 1] = words[i]
		}
		args[limit - 1] = strings.Join(lastArgSlice, " ")

		return args
	} else {
		return words
	}
}
