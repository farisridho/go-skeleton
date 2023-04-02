package todo

type TodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TodoResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
