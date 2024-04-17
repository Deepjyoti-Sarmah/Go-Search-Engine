package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
  ID string `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
  Email string `gorm:"unique" json:"email"`
  Password string `json:"-"`
  IsAdmin bool `gorm:"default:false" json:"isAdmin"`
  CreatedAt *time.Time  `json:"createdAt"`
  UpdatedAt time.Time  `json:"updatedAt"`
}

func (u *User) CreateAdmin() error {
  user := User {
	Email: "your email",
	Password: "Your Password",
	IsAdmin: true,
  }

  // Hash Password
  password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
  if err != nil {
	return errors.New("error creating password")
  }
  user.Password = string(password)

  //create user in db
  if err := DBConn.Create(&user).Error; err != nil {
	return errors.New("Error creating user")
  }
  return nil
}

func (u *User) LoginAsAdmin(email string, password string) (*User, error) {
  //Find 
  if err := DBConn.Where("email = ? AND is_Admin = ?", email, true).First(&u).Error; err != nil {
	return nil, errors.New("user not found")
  } 

  //Compare the password
  if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
	return nil, errors.New("Invalid password")
  }
  return u, nil
}
