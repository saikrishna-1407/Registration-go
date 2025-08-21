package forms

import "time"

type UserForm struct {
	FullName        string    `json:"full_name" binding:"required,max=50"`
	Email           string    `json:"email" binding:"required,email"`
	Password        string    `json:"password" binding:"required"`
	ConfirmPassword string    `json:"confirmPassword" binding:"required"`
	Username        string    `json:"username" binding:"required,max=20"`
	DateOfBirth     time.Time `json:"dateOfBirth" binding:"required" time_format:"2006-01-02"`
	PhoneNumber     string    `json:"phone_number" binding:"required"`
	Gender          string    `json:"gender"`
	Country         string    `json:"country" binding:"required"`
	State           string    `json:"state"`
	PinCode         string    `json:"pin_code"`
	ReferralCode    string    `json:"referral_code"`
	TermsAccepted   bool      `json:"terms_accepted" binding:"required"`
}
