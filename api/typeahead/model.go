package typeahead

type Suggestions struct {
	Category interface{} `json:"category"`
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
}

type SuggestionResponse struct {
	Content struct {
		Suggests []Suggestions `json:"suggests"`
	} `json:"content"`
	ErrorCode interface{} `json:"errorCode"`
	Message   interface{} `json:"message"`
}

