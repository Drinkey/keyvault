package models

type Namespace struct {
	ID        uint   `json:"namespace_id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"unique;not null"`
	MasterKey string
	Nonce     string
}

func CreateNamespace(name, key, nonce string) (err error) {
	err = db.Create(&Namespace{
		Name:      name,
		MasterKey: key,
		Nonce:     nonce,
	}).Error
	if err != nil {
		return
	}
	return nil
}

func GetNamespace(name string) (n Namespace, err error) {
	err = db.Where("name = ?", name).First(&n).Error
	if err != nil {
		return Namespace{}, err
	}
	return n, nil
}

func ListNamespace() (n []Namespace, err error) {
	err = db.Find(&n).Error
	if err != nil {
		return nil, err
	}
	return
}

func (ns Namespace) IsEmpty() bool {
	return ns.ID == 0
}
