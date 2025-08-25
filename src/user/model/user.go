// st-novel-go/src/user/model/user.go
package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type PlanType string

const (
	PlanFree PlanType = "免费版"
	PlanPro  PlanType = "专业版"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Avatar    string         `gorm:"type:mediumtext" json:"avatar"`
	Bio       string         `gorm:"type:text" json:"bio"`
	Phone     string         `gorm:"type:varchar(50)" json:"phone"`
	Region    string         `gorm:"type:varchar(100)" json:"region"`
	Timezone  string         `gorm:"type:varchar(100)" json:"timezone"`
	Plan      PlanType       `gorm:"type:varchar(50);default:'免费版'" json:"plan"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	// Set default values for new users
	if u.Region == "" {
		u.Region = "China"
	}
	if u.Timezone == "" {
		u.Timezone = "Asia/Shanghai (UTC+8)"
	}
	if u.Plan == "" {
		u.Plan = PlanFree
	}
	return
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
