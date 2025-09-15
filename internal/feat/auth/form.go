package auth

// UserForm represents the form data for creating/updating a user
type UserForm struct {
	Username string `form:"username" required:"true"`
	Email    string `form:"email" required:"true"`
	Name     string `form:"name" required:"true"`
}
