package users

type User struct {
	Username string   `json:"username,omitempty"`
	Email    string   `json:"email,omitempty"`
	Address  *Address `json:"address,omitempty"`
}

type Address struct {
	Street string `json:"street,omitempty"`
}
