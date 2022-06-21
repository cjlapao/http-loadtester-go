package entities

// LoadTestJobTarget Entity
type LoadTestJobTarget struct {
	URL          string   `json:"url" yaml:"url"`
	Method       string   `json:"method" yaml:"method"`
	Body         string   `json:"body" yaml:"body"`
	BearerToken  string   `json:"token" yaml:"token"`
	BearerTokens []string `json:"tokens" yaml:"tokens"`
	ContentType  string   `json:"contentType" yaml:"contentType"`
	Timeout      int      `json:"timeout" yaml:"timeout"`
	LogResponse  bool     `json:"logResponse" yaml:"logResponse"`
}
