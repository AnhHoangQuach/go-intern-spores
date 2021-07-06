package models

import (
	"fmt"
	"html"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Phone     string    `gorm:"size:255;not null;unique" json:"phone"`
	Address   string    `gorm:"size:255;not null" json:"address"`
	VerifyToken string    `gorm:"size:255;not null" json:"verify_token"`
	ResetToken string    `gorm:"size:255;not null" json:"reset_token"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func HashPass(email, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(email+password), 14)
	return string(bytes), err
}

func CheckPasswordHash(email, hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(email+password))
}

func (u *User) Save(db *gorm.DB) error {
	err := db.Debug().Create(&u).Error
	if (err != nil) {
		return err
	}
	return nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	err := db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	err := db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, err
}

func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	err := db.Debug().Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, err
}

func (u *User) Delete(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// randToken generates a random hex value.
func RandomSixDigits(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	randCode := strconv.Itoa(min + rand.Intn(max-min))
	return randCode
}

func (u *User) Prepare() {
	u.ID = 0
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) SignUp(db *gorm.DB, email, password, phone, address string) error {
	user, err := u.FindUserByEmail(db, email)
	verify_token := RandomSixDigits(100000, 999999)
	subject := "Token For Verify Email"
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

	user.Prepare()

	user = &User{
		Email:       email,
		Password:    hashPass,
		Phone: 		phone,
		Address: address,
		VerifyToken: verify_token,
	}

	err = u.Save(db)
	if err != nil {
		return fmt.Errorf("Sign up failed %v", err)
	}

	return nil
}