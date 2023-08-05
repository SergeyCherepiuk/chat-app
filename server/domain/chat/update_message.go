package chatdomain

type UpdateMessageRequestBody struct {
	Message string `json:"message"`
}

// TODO: Consider validation (trimming whitespace)