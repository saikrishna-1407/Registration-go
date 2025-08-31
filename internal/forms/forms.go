package forms

import "time"

type UserForm struct {
	FullName        string    `json:"fullName" binding:"required,max=50"`
	Email           string    `json:"email" binding:"required,email"`
	Password        string    `json:"password" binding:"required,min=6"`
	ConfirmPassword string    `json:"confirmPassword" binding:"required"`
	Username        string    `json:"username" binding:"required,max=20"`
	DateOfBirth     time.Time `json:"dateOfBirth" binding:"required" time_format:"2006-01-02" time_utc:"true"`
	PhoneNumber     string    `json:"phoneNumber" binding:"required"`
	Gender          string    `json:"gender"`
	Country         string    `json:"country" binding:"required"`
	State           string    `json:"state"`
	PinCode         string    `json:"pinCode"`
	ReferralCode    string    `json:"referralCode"`
	TermsAccepted   bool      `json:"termsAccepted" binding:"required"`
}
