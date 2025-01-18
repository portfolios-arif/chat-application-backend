package requests

import "mime/multipart"

type CheckEmailRequestPayload struct {
	Email   string `json:"email" validate:"required,email"`
	IsLogin bool   `json:"isLogin"`
}

type OTPRequestPayload struct {
	DeviceID string `json:"deviceId" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type ValidateOTPRequestPayload struct {
	OTPCode     interface{} `json:"otp_code" validate:"required,len=6"`
	SignatureID string      `json:"signatureId" validate:"required"`
}

type RegisterRequestPayload struct {
	Email     string                `form:"email" binding:"required,email"`
	DeviceID  string                `form:"deviceId" binding:"required"`
	Username  string                `form:"username" binding:"required,min=6,max=15,ascii"`
	Fullname  string                `form:"fullName" binding:"required,min=5,max=20"`
	Gender    string                `form:"gender" binding:"omitempty,alpha"`
	Age       int8                  `form:"age" binding:"omitempty,numeric"`
	PublicKey *multipart.FileHeader `form:"publicKey" binding:"required"`
	ImgFile   *multipart.FileHeader `form:"imgFile" binding:"omitempty"`
}

type LoginRequestPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type RefreshTokenRequestPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
