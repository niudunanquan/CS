package db

type Upstream struct {
	Model
	Name   string `gorm:"type:varchar(64)" json:"name"`
	Config string `json:"config"`
}

func (t *Upstream) Count() (int64, error) {
	var count int64

	err := Db.Model(t).Count(&count).Error

	return count, err
}

func (t *Upstream) Get(id uint) error {
	return Db.First(t, "id = ?", id).Error
}

func (t *Upstream) GetAll() ([]Upstream, error) {
	var upstreams []Upstream
	err := Db.Order("id desc").Find(&upstreams).Error
	if err != nil {
		return nil, err
	}
	return upstreams, err
}

func (t *Upstream) GetFields(fields []string) ([]Upstream, error) {
	var upstreams []Upstream
	err := Db.Select(fields).Order("id desc").Find(&upstreams).Error
	if err != nil {
		return nil, err
	}
	return upstreams, err
}

func (t *Upstream) Insert() error {
	return Db.Create(t).Error
}

func (t *Upstream) Update() error {
	return Db.Save(t).Error
}

func (t *Upstream) Updates(columns interface{}) error {
	return Db.Model(t).Updates(columns).Error
}

func (t *Upstream) Delete() error {
	return Db.Delete(t).Error
}

func (t *Upstream) DeleteAll(ids []uint) error {
	return Db.Delete(t, "id IN ?", ids).Error
}
