package models

type Certificate struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"unique;not null"`
	SignRequest string `json:"req" gorm:"not null"`
	Certificate string `json:"certificate"`
	Token       string `json:"token" gorm:"unique;not null"`
}

func (c Certificate) IsEmpty() bool {
	return c.ID == 0
}

func CreateCertificateRequest(name, req, token string) (err error) {
	err = db.Create(&Certificate{
		Name:        name,
		SignRequest: req,
		Token:       token,
	}).Error
	if err != nil {
		return
	}
	return nil
}

func GetCertificate(name string) (cert Certificate, err error) {
	err = db.Where("name = ?", name).First(&cert).Error
	if err != nil {
		return Certificate{}, err
	}
	return
}

func UpdateSignedCertificateByName(name, signed string) (cert Certificate, err error) {
	err = db.Model(&cert).Where("name = ?", name).Update("certificate", signed).Error
	if err != nil {
		return
	}
	return GetCertificate(name)
}
