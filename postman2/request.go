package postman2

type Request struct {
	URL         *URL         `json:"url,omitempty"`
	Method      string       `json:"method,omitempty"`
	Header      []Header     `json:"header,omitempty"`
	Body        *RequestBody `json:"body,omitempty"`
	Description string       `json:"description,omitempty"`
}

type Header struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type RequestBody struct {
	Mode       string            `json:"mode,omitempty"` // `raw`, `urlencoded`, `formdata`,`file`,`graphql`
	Raw        string            `json:"raw,omitempty"`
	URLEncoded []URLEncodedParam `json:"urlencoded,omitempty"`
}

type URLEncodedParam struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
