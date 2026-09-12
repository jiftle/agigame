package agent

import (
	"fmt"
	"strings"

	"agigame/agent/games"
)

// BuildGamePrompt renders the state through the plugin and wraps it with the
// shared output schema / button vocabulary used by every provider.
func BuildGamePrompt(plugin games.GamePlugin, state any) string {
	body := plugin.Prompt(state)
	var b strings.Builder
	fmt.Fprintf(&b,
		"Game: %s.\n"+
			"Available buttons: %s.\n"+
			"Answer with exactly one JSON object of shape "+
			"{\"commands\":[\"A\",\"Right\"],\"rationale\":\"short reason\"}.\n\n",
		plugin.Name(), AvailableButtonsText())
	b.WriteString(body)
	return b.String()
}