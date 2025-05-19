package db

import (
	"time"
)

type Cert struct {
	Model
	Name          string    `gorm:"type:varchar(64)" json:"name"`
	Type          uint      `json:"type"`
	DnsChallenge  bool      `json:"dns_challenge"`
	DnsProvider   string    `gorm:"type:varchar(64)" json:"dns_provider"`
	DnsCredential string    `json:"dns_credential"`
	Domains       string    `json:"domains"`
	Email         string    `json:"email"`
	Crt           string    `json:"crt"`
	Key           string    `json:"key"`
	Expires       time.Time `json:"expires"`
}

func (t *Cert) Count() (int64, error) {
	var count int64

	err := Db.Model(t).Count(&count).Error

	return count, err
}

func (t *Cert) GetAll() ([]Cert, error) {
	var certs []Cert
	err := Db.Order("id desc").Find(&certs).Error
	if err != nil {
		return nil, err
	}
	return certs, err
}

func (t *Cert) GetFields(fields []string) ([]Cert, error) {
	var certs []Cert
	err := Db.Select(fields).Order("id desc").Find(&certs).Error
	if err != nil {
		return nil, err
	}
	return certs, err
}

func (t *Cert) Insert() error {
	return Db.Create(t).Error
}

func (t *Cert) Update() error {
	return Db.Save(t).Error
}

func (t *Cert) DeleteAll(ids []uint) error {
	return Db.Delete(t, "id IN ?", ids).Error
}
