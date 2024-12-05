package model

type Employee struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Salary    float64 `json:"salary"`
	Position  string  `json:"position"`
	Email     string  `json:"email"`
}
