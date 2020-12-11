package models

type Certificate struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"unique;not null"`
	SignRequest string `gorm:"not null"`
	Certificate string
	Token       string `gorm:"unique;not null"`
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
