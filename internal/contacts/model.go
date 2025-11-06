package contacts

type Contact struct {
	ID    int    `json:"id"    gorm:"primaryKey;autoIncrement"`
	Nom   string `json:"nom"   gorm:"not null"`
	Email string `json:"email" gorm:"not null;uniqueIndex"`
}

func (c Contact) IsValid() bool {
	return c.Nom != "" && c.Email != ""
}
