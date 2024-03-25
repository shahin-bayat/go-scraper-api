package models

type CreateUserRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
	// other fields required for creating a user
}

type CreateUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	// other fields returned in the response after creating a user
}
