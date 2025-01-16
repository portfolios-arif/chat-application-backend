package requests

type CheckPhoneRequestPayload struct {
	Phonenumber string `json:"phonenumber" validate:"required,min=7,max=15,numeric"`
	IsLogin     bool   `json:"isLogin" validate:"required,boolean"`
}

type OTPRequestPayload struct {
	DeviceID    string `json:"deviceId" validate:"required,alphanum"`
	Phonenumber string `json:"phonenumber" validate:"required,min=7,max=15,numeric"`
}

type ValidateOTPRequestPayload struct {
	OTPCode     interface{} `json:"otp_code" validate:"required,len=4"`
	SignatureID string      `json:"signatureId" validate:"required"`
}

type RegisterRequestPayload struct {
	Phonenumber string      `json:"phonenumber" validate:"required,min=7,max=15,numeric"`
	DeviceID    string      `json:"deviceId" validate:"required,alphanum"`
	PublicKey   string      `json:"publicKey" validate:"required,base64"`
	Username    string      `json:"username" validate:"required,min=6,max=15,ascii"`
	FullName    string      `json:"fullName" validate:"required,min=5,max=20"`
	Gender      string      `json:"gender" validate:"required,alpha"`
	Age         interface{} `json:"age" validate:"required,numberic"`
}

type LoginRequestPayload struct {
	Phonenumber string `json:"phonenumber" validate:"required,min=7,max=15,numeric"`
}
