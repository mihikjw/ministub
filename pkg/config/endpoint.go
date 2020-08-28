package config

// Endpoint represents a definition for an endpoint
type Endpoint struct {
	Params    *Parameters              `yaml:"params"`
	Recieves  *Recieves                `yaml:"recieves"`
	Response  int                      `yaml:"response"`
	Responses map[int]*Response        `yaml:"responses"`
	Actions   []map[string]interface{} `yaml:"actions"`
}

// Parameters represents the parameters in an HTTP endpoint
type Parameters struct {
	Query map[string]*ParamEntry `yaml:"query"`
	URL   map[string]*ParamEntry `yaml:"url"`
}

// ParamEntry represents the parameters for a single parameter
type ParamEntry struct {
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

// Recieves represents the 'recieves' field of an endpoint
type Recieves struct {
	Headers map[string]string      `yaml:"recieves"`
	Body    map[string]interface{} `yaml:"body"`
}
