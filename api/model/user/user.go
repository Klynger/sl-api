package userModel

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Name      string    `gorm:"column:user_name"`
	LastName  string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type PublicUser struct {
	ID       uuid.UUID
	Name     string
	LastName string
}

type LoginForm struct {
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required"`
}

type UserCredentials struct {
	ID       uuid.UUID
	Username string
	Password string
}

func (UserCredentials) TableName() string {
	return "users"
}

// ValidateCredentials reports whether the login form's password matches the
// stored bcrypt hash. The comparison is constant time, so it does not leak how
// much of the password was correct.
func (u *UserCredentials) ValidateCredentials(form *LoginForm) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(form.Password)) == nil
}

type DTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Username string `json:"username"`
}

type Form struct {
	Name     string `json:"name" validate:"required,max=255"`
	LastName string `json:"lastName" validate:"required,max=255"`
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required"`
}

func (u *User) ToDto() *DTO {
	return &DTO{
		ID:       u.ID.String(),
		Name:     u.Name,
		LastName: u.LastName,
		Username: u.Username,
	}
}

// ToModel builds a User from the registration form, storing the password as a
// bcrypt hash rather than in plaintext. It returns an error if hashing fails.
func (f *Form) ToModel() (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:     f.Name,
		LastName: f.LastName,
		Username: f.Username,
		Password: string(hashed),
	}, nil
}
