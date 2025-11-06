// internal/contacts/model.go
package contacts

type Contact struct {
	ID    int    `json:"id"`
	Nom   string `json:"nom"`
	Email string `json:"email"`
}

func (c Contact) IsValid() bool {
	return c.Nom != "" && c.Email != ""
}
