package entity

// User represents a user.
type User struct {
	ID   int
	Name string
}

// GetID returns the user ID.
func (u User) GetID() int {
	return u.ID
}

// GetName returns the user name.
func (u User) GetName() string {
	return u.Name
}
