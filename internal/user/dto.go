package user

// User represents a user in the system
type User struct {
	ID    int    `json:"id" doc:"User ID"`
	Name  string `json:"name" doc:"User full name"`
	Email string `json:"email" doc:"User email address"`
	Role  string `json:"role" doc:"User role (admin, teacher, student)"`
}

// UserListResponse represents the response for user list
type UserListResponse struct {
	Data []User   `json:"data"`
	Meta Metadata `json:"meta"`
}

// UserResponse represents a single user response
type UserResponse struct {
	User
}

// Metadata represents pagination metadata
type Metadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// CreateUserRequest represents request to create a new user
type CreateUserRequest struct {
	Name     string `json:"name" minLength:"1" maxLength:"100" doc:"User full name"`
	Email    string `json:"email" format:"email" doc:"User email address"`
	Password string `json:"password" minLength:"8" doc:"User password"`
	Role     string `json:"role" enum:"admin,teacher,student" doc:"User role"`
}

// CreateUserResponse represents response after creating a user
type CreateUserResponse struct {
	ID      int    `json:"id" doc:"Created user ID"`
	Message string `json:"message" doc:"Success message"`
}

// UpdateUserRequest represents request to update user
type UpdateUserRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"100" doc:"User full name"`
	Role string `json:"role" enum:"admin,teacher,student" doc:"User role"`
}

// UserBasicResponse represents a basic response with message for user operations
type UserBasicResponse struct {
	Message string `json:"message" doc:"Success message"`
}
