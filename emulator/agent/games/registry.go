package games

// registry maps game plugin ids to their constructors. New games register
// themselves from an init() in their own file; the engine resolves the
// config.yaml "game" value through Get, so no host needs a switch statement.
var registry = map[string]func() GamePlugin{}

// Register adds a plugin constructor under an id.
func Register(id string, factory func() GamePlugin) {
	registry[id] = factory
}

// Get returns a new plugin instance for the id.
func Get(id string) (GamePlugin, bool) {
	factory, ok := registry[id]
	if !ok {
		return nil, false
	}
	return factory(), true
}

// IDs returns the registered plugin ids.
func IDs() []string {
	ids := make([]string, 0, len(registry))
	for id := range registry {
		ids = append(ids, id)
	}
	return ids
}

func init() {
	Register("sml", func() GamePlugin { return &SMLPlugin{} })
}
