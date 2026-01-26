package login
type LoginCmd struct {
	Username    string
	Password string
}

type LoginRes struct {
	AccessToken  string
	RefreshToken string
}