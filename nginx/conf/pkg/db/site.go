package db

type Site struct {
	Model
	Name      string `gorm:"type:varchar(64)" json:"name"`
	Domains   string `json:"domains"`
	Listeners string `json:"listeners"`
	CertId    uint   `json:"cert_id"`
	Locations string `json:"locations"`
	Http2     bool   `json:"http2"`
	Ipv6      bool   `json:"ipv6"`
	Cache     bool   `json:"cache"`
	Gzip      bool   `json:"gzip"`
	ForceSsl  bool   `json:"force_ssl"`
	Hsts      bool   `json:"hsts"`
}

func (t *Site) Count() (int64, error) {
	var count int64

	err := Db.Model(t).Count(&count).Error

	return count, err
}

func (t *Site) CertCount(ids []uint) (int64, error) {
	var count int64

	err := Db.Model(t).Where("cert_id IN ?", ids).Count(&count).Error

	return count, err
}

func (t *Site) Get(id uint) error {
	return Db.First(t, "id = ?", id).Error
}

func (t *Site) GetAll() ([]Site, error) {
	var sites []Site
	err := Db.Order("id desc").Find(&sites).Error
	if err != nil {
		return nil, err
	}
	return sites, err
}

func (t *Site) GetLocations() ([]Site, error) {
	var sites []Site
	err := Db.Select("locations").Order("id desc").Find(&sites).Error
	if err != nil {
		return nil, err
	}
	return sites, err
}

func (t *Site) Insert() error {
	return Db.Create(t).Error
}

func (t *Site) Update() error {
	return Db.Save(t).Error
}

func (t *Site) Updates(columns interface{}) error {
	return Db.Model(t).Updates(columns).Error
}

func (t *Site) Delete() error {
	return Db.Delete(t).Error
}

func (t *Site) DeleteAll(ids []uint) error {
	return Db.Delete(t, "id IN ?", ids).Error
}
