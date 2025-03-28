package entities

type OTPData struct {
	Value      string `json:"value"`
	OTP        string `json:"otp"`
	IsVerified bool   `json:"is_verified"`
	Expire     int64  `json:"expire"`
}
