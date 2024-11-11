package models

type User struct {
	ID          uint   `gorm:"primaryKey; autoIncrement:true;unique"`
	UUID        string `gorm:"unique"`
	Email       string `gorm:"not null;unique" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	Role        string `json:"role"`
	Address     string `json:"address"`
	DateOfBirth string `json:"dateOfBirth"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone"`
}
