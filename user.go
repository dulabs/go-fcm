package fcm

type User struct {
	LocalId       string `json:"localId,omitempty"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	PhotoUrl      string `json:"photoUrl,omitempty"`
	ValidSince    string `json:"validSince,omitempty"`
	Disabled      string `json:"disabled,omitempty"`
	LastLoginAt   string `json:"lastLoginAt,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`

	ProviderUserInfo []ProviderInfo `json:"providerUserInfo,omitempty"`
}

type ProviderInfo struct {
	ProviderId  string `json:"providerId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	PhotoUrl    string `json:"photoUrl,omitempty"`
	Email       string `json:"email,omitempty"`
}

type LookupUserResponse struct {
	Kind  string `json:"kind,omitempty"`
	Users []User `json:"users,omitempty"`
}
