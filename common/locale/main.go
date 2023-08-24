package locale

import (
	"errors"
	"fmt"

	"github.com/d0ctr/bldbr-bot-go/common/locale/en"
	"github.com/d0ctr/bldbr-bot-go/common/locale/ru"
)

type Language string

type formattable = func(...any) string

func preformat(line string) formattable {
	return func(a ...any) string {
		if len(a) != 0 {
			return fmt.Sprintf(line, a...)
		}
		return line
	}
}

type Locales = map[string]formattable

var AvailableLocales = map[Language]bool{
	"ru":      true,
	"en":      true,
	"default": true,
}

var Lines = map[Language]map[string]string{
	"ru":      ru.Lines(),
	"en":      en.Lines(),
	"default": en.Lines(),
}

func Get(lang Language, name string, a ...any) (string, error) {
	if !AvailableLocales[lang] {
		lang = "default"
	}
	if Lines[lang][name] == "" {
		return "", errors.New("no key " + name + " in " + string(lang) + " locale")
	}
	return preformat(Lines[lang][name])(a...), nil
}
