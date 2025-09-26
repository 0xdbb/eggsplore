package server

// HandleError represents the structure of error messages in API responses
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type UserMessage struct {
	Message string `json:"message" example:"success"`
}

type Message struct {
	Message string `json:"message" example:"Task status updated - COMPLETED"`
	TaskID  string `json:"task_id" binding:"omitempty" example:"1233-flf4djf-alsdik"`
	Result  string `json:"result" binding:"omitempty" example:"Task processed successfully"`
	Error   string `json:"error" binding:"omitempty" example:"Failed to Process Task"`
	Code    int    `json:"code" binding:"omitempty" example:"500"`
}

func HandleError(err error, code int, message ...string) ErrorResponse {
	var msg string

	if err != nil {
		msg = err.Error()
	}

	if len(message) > 0 {
		if msg == "" {
			msg = message[0]
		} else {
			msg = message[0] + ": " + msg // Override error message if provided
		}
	}

	return ErrorResponse{
		Status:  "error",
		Message: msg,
		Code:    code,
	}
}

func HandleMessage(message string) UserMessage {
	return UserMessage{Message: message}
}

func fallback(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}
