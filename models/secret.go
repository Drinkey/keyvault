package models

// type Secret struct {
// 	ID        int       `json:"id"`
// 	NameSpace Namespace `json:"namespace"`
// 	Key       string    `json:"key"`
// 	Value     string    `json:"value"`
// }
type Secret struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	NamespaceID uint      `json:"namespace_id" gorm:"foreignKey:NamespaceID"`
	Namespace   Namespace `json:"namespace"`
}

func (s Secret) IsEmpty() bool {
	return s.ID == 0
}

func CreateSecret(key, value, namespace string) (err error) {
	ns, err := GetNamespace(namespace)
	if err != nil {
		return
	}
	err = db.Create(&Secret{
		Key:         key,
		Value:       value,
		NamespaceID: ns.ID,
	}).Error
	if err != nil {
		return
	}
	return nil
}

func GetSecret(key, namespace string) (s Secret, err error) {
	ns, err := GetNamespace(namespace)
	err = db.Where("key = ?", key).Where("namespace_id = ?", ns.ID).First(&s).Error
	if err != nil {
		return Secret{}, err
	}
	s.Namespace = ns
	return
}
