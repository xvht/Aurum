package enums

const (
	KeyCheckFailure string = "Key check failed"   // Key check failed
	KeyDeactivated  string = "Key is deactivated" // Key is deactivated

	InvalidHWID string = "Invalid HWID" // HWID is invalid
	ValidHWID   string = "Valid HWID"   // HWID is valid

	InvalidKey string = "Invalid key" // Key is invalid
	ValidKey   string = "Valid key"   // Key is valid

	KeyGenOK    string = "Key generated"                           // Key generated
	KeyGenLimit string = "Wait 7 days before generating a new key" // Key generation limit reached
)
