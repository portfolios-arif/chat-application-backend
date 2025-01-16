package responses

type CheckPhoneIsLoginResponsePayload struct {
	DeviceID string `json:"deviceId"`
}

type OTPResponsePayload struct {
	SignatureID string `json:"signatureId"`
}

type LoginResponsePayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponsePayload struct {
	AccessToken string `json:"accessToken"`
}
