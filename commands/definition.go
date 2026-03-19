package commands

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type CommandDefinition struct {
	Name, Description string
	Handler handlers.Response
}

var All = func() map[string]CommandDefinition {
	all_a := []CommandDefinition{Ping, Ahegao, Urban, Get, Set, Lst, Voice}

	all_m := make(map[string]CommandDefinition, len(all_a))
	
	for _, def := range all_a {
		all_m[def.Name] = def
	}

	return all_m
}()
