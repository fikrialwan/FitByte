package utils

type (
	FailedResponse struct {
		Message string `json:"message"`
	}
)

func BuildFailedResponse(message string) FailedResponse {
	return FailedResponse{Message: message}
}
