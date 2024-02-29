package responses

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

const (
	InvalidUserResponse        = "Are you fucking dumb? It's a 'ping-pong' server. Send me just a one word - 'ping'"
	InternalServerErrorMessage = "He-he-he, server was written with some surprises"
)
