package entities

// JobTarget Entity
type JobTarget struct {
	URL                  string            `json:"url" yaml:"url"`
	URLs                 []string          `json:"urls" yaml:"urls"`
	Method               string            `json:"method" yaml:"method"`
	RawBody              string            `json:"body" yaml:"body"`
	FormData             map[string]string `json:"formData" yaml:"formData"`
	FormUrlEncoded       map[string]string `json:"formUrlEncoded" yaml:"formUrlEncoded"`
	BearerToken          string            `json:"token" yaml:"token"`
	BearerTokens         []string          `json:"tokens" yaml:"tokens"`
	BasicAuthentication  string            `json:"basicAuthentication" yaml:"basicAuthentication"`
	BasicAuthentications []string          `json:"basicAuthentications" yaml:"basicAuthentications"`
	ContentType          string            `json:"contentType" yaml:"contentType"`
	Headers              map[string]string `json:"headers" yaml:"headers"`
	UserAgent            string            `json:"userAgent" yaml:"userAgent"`
	Timeout              int               `json:"timeout" yaml:"timeout"`
	LogResponse          bool              `json:"logResponse" yaml:"logResponse"`
}
