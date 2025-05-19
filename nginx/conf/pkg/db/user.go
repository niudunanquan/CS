package db

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model
	Username  string `gorm:"type:varchar(32);unique" json:"username"`
	Password  string `gorm:"type:varchar(512)" json:"password"`
	Email     string `gorm:"type:varchar(64)" json:"email"`
	Role      int    `json:"role"`
	OtpUrl    string `gorm:"type:varchar(256)" json:"otp_url"`
	EnableOtp bool   `json:"enable_otp"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.OtpUrl, err = u.GetOtpUrl()
	if err != nil {
		return
	}
	err = u.HashPwd(u.Password)
	return
}

func (t *User) HashPwd(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err == nil {
		t.Password = string(hash)
	}
	return err
}

func (t *User) VerifyPwd(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.Password), []byte(pwd))
	return err == nil
}

func (t *User) GetOtpUrl() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OpenResty Manager",
		AccountName: t.Username,
	})
	if err != nil {
		return "", err
	}
	return key.URL(), err
}

func (t *User) VerifyOtp(passcode string) bool {
	key, err := otp.NewKeyFromURL(t.OtpUrl)
	if err != nil {
		return false
	}
	return totp.Validate(passcode, key.Secret())
}

func (t *User) Get(id uint) error {
	return Db.First(t, "id = ?", id).Error
}

func (t *User) GetAll() ([]User, error) {
	var users []User
	err := Db.Select("id", "username", "email", "role", "otp_url", "enable_otp", "updated_at").Order("id desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, err
}

func (t *User) Take(where ...interface{}) error {
	return Db.Take(t, where...).Error
}

func (t *User) GetByUsername() error {
	return t.Take("username = ?", t.Username)
}

func (t *User) Count() (int64, error) {
	var count int64

	if err := Db.Model(t).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (t *User) FindPageByParams(page, size int, where ...interface{}) (*Pagination, error) {
	p := Pagination{Size: size, Page: page}

	if err := Db.Model(t).Count(&p.Total).Error; err != nil {
		return nil, err
	}

	var users []User
	if err := Db.Model(t).Offset(p.GetOffset()).Limit(p.GetSize()).Order(p.GetSort()).Find(&users, where...).Error; err != nil {
		return nil, err
	}
	p.Data = &users

	return &p, nil
}

func (t *User) Insert() error {
	return Db.Create(t).Error
}

func (t *User) Update() error {
	return Db.Save(t).Error
}

func (t *User) Updates(values interface{}) error {
	return Db.Model(t).Updates(values).Error
}

func (t *User) Delete(ids []uint) error {
	return Db.Delete(t, "id IN ?", ids).Error
}
