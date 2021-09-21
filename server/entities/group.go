package entities

import (
	"github.com/talkiewalkie/talkiewalkie/models"
	"strings"
)

func GroupDisplay(g models.Group) string {
	handles := []string{}
	for _, ug := range g.R.UserGroups {
		redundant := false
		for _, h := range handles {
			if h == ug.R.User.Handle {
				redundant = true
			}
		}
		if !redundant {
			handles = append(handles, ug.R.User.Handle)
		}
	}

	display := g.Name.String
	if !g.Name.Valid {
		display = strings.Join(handles, ", ")
	}
	return display
}
