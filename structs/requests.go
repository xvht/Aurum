package structs

type AuthRequest struct {
	Username   string `json:"username"`
	LicenseKey string `json:"license_key"`
	HWID       string `json:"hwid"`
	Nonce      string `json:"nonce"`
	Timestamp  int64  `json:"timestamp"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` // Do *NOT* use this for authentication, always hash it
}

type AdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` // Do *NOT* use this for authentication, always hash it

	TargetUser string `json:"target_user"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` // Do *NOT* use this for authentication, always hash it

	Type     string `json:"type"` // username or password
	NewValue string `json:"new_value"`
}

type InviteRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` // Do *NOT* use this for authentication, always hash it

	InviteCode string `json:"invite_code"`
}
