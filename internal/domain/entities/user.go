package entities

type Mst_otp struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid()"`
	OTPCode     string `json:"otp_code"`
	DeviceID    string `gorm:"index" json:"device_id"`
	SignatureID string `gorm:"index" json:"signature_id"`
	BaseTableDTO
}

type Mst_users struct {
	ID            string           `gorm:"type:uuid;default:gen_random_uuid()"`
	Phonenumber   string           `gorm:"index;unique" json:"phone_number"`
	DeviceID      string           `json:"device_id"`
	PublicKeyPath string           `json:"public_key_path"`
	IsOnline      bool             `gorm:"default:false" json:"is_online"`
	IsDeleted     bool             `gorm:"default:false" json:"is_deleted"`
	UserDetail    Mst_users_detail `gorm:"foreignKey:UserID;reference:ID"`
	BaseTableDTO
}

type Mst_users_detail struct {
	ID       string `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID   string `gorm:"uniqueIndex" json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	ImgPath  string `json:"img_path"`
	Gender   string `json:"gender"`
	Age      int8   `json:"age"`
	BaseTableDTO
}
