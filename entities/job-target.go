package entities

// JobTarget Entity
type JobTarget struct {
	URL          string   `json:"url" yaml:"url"`
	URLs         []string `json:"urls" yaml:"urls"`
	Method       string   `json:"method" yaml:"method"`
	Body         string   `json:"body" yaml:"body"`
	BearerToken  string   `json:"token" yaml:"token"`
	BearerTokens []string `json:"tokens" yaml:"tokens"`
	ContentType  string   `json:"contentType" yaml:"contentType"`
	Timeout      int      `json:"timeout" yaml:"timeout"`
	LogResponse  bool     `json:"logResponse" yaml:"logResponse"`
}
