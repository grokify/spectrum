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

/*
	"request": {
		"method": "POST",
		"header": [
			{
				"key": "Content-Type",
				"value": "application/json"
			},
			{
				"key": "Accept",
				"value": "application/json"
			}
		],
		"body": {
			"mode": "raw",
			"raw": "{\n\t\"account\": \"r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59\",\n\t\"account_index\": 0,\n\t\"ledger_index\": \"validated\",\n\t\"strict\": true\n}"
		},
		"url": {
			"raw": "{{WECOINS_RIPPLE_BASE_URL}}/api/v1/account_currencies",
			"host": [
				"{{WECOINS_RIPPLE_BASE_URL}}"
			],
			"path": [
				"api",
				"v1",
				"account_currencies"
			]
		}
	},
*/
