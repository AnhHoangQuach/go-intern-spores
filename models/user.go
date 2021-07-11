package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Email       string    `gorm:"size:100;not null;unique" json:"email"`
	Password    string    `gorm:"size:100;not null;" json:"password"`
	Phone       string    `gorm:"size:255;not null;unique" json:"phone"`
	Address     string    `gorm:"size:255;not null" json:"address"`
	VerifyToken string    `gorm:"size:255;not null" json:"verify_token"`
	ResetToken  string    `gorm:"size:255;not null" json:"reset_token"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type UserModel struct{}

func (u *UserModel) Save(user *User) error {
	if err := DB.Create(&user).Error; err != nil {
		return fmt.Errorf("Save user failed")
	}
	return nil
}

func (u *UserModel) Update(user *User) error {
	if err := DB.Model(&user).Where("email = ?", user.Email).Save(&user).Error; err != nil {
		return fmt.Errorf("Save user failed")
	}
	return nil
}

func (u *UserModel) FindByEmail(email string) (*User, error) {
	var result User
	if err := DB.Where("email = ?", email).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func RandomSixDigits(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	randCode := strconv.Itoa(min + rand.Intn(max-min))
	return randCode
}

func HashPass(email, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(email+password), 14)
	return string(bytes), err
}

func CheckPasswordHash(email, hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(email+password))
}

func (u *UserModel) SignUp(email, password, phone, address string) error {
	user, err := u.FindByEmail(email)
	verify_token := RandomSixDigits(100000, 999999)
	subject := "Email Verify Token"
	if err == nil && user != nil {
		return fmt.Errorf("User existed")
	}

	hashPass, err := HashPass(email, password)
	if err != nil {
		return fmt.Errorf("Sign up failed %v", err)
	}

	m := utils.NewMailSender([]string{email}, subject)
	err = m.Send("./utils/mailContent.html", map[string]string{"token": verify_token})
	if err != nil {
		fmt.Errorf("Send email failed %v", err)
	}

	user = &User{
		Email:       email,
		Password:    hashPass,
		Phone:       phone,
		Address:     address,
		VerifyToken: verify_token,
	}

	err = u.Save(user)
	if err != nil {
		return fmt.Errorf("Sign up failed %v", err)
	}

	return nil
}

func (u *UserModel) ActiveUser(email, verify_token string) error {
	user, err := u.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("Failed Email")
	}

	if user.VerifyToken != verify_token {
		return fmt.Errorf("Verify token is failed")
	}
	user.VerifyToken = ""
	err = u.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) LoginHandler(email string, password string) (*User, error) {
	user, err := u.FindByEmail(email)

	if err != nil {
		return nil, fmt.Errorf("This account does not exist")
	}

	err = CheckPasswordHash(email, user.Password, password)
	if err != nil {
		return nil, fmt.Errorf("Password is wrong")
	}

	if user.VerifyToken != "" {
		verify_token := RandomSixDigits(100000, 999999)

		subject := "Email Verify Token"

		m := utils.NewMailSender([]string{user.Email}, subject)
		err = m.Send("./utils/mailContent.html", map[string]string{"token": verify_token})
		if err != nil {
			fmt.Errorf("Send email failed %v", err)
		}

		user.VerifyToken = verify_token
		err = u.Save(user)
		if err != nil {
			fmt.Errorf("Save failed %v", err)
		}
		return nil, fmt.Errorf("This account hasn't been activated, a verification code has been sent to your email, please check")
	}

	return user, nil
}

func (u *UserModel) GetProfile(email string) (*User, error) {
	user, err := u.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
