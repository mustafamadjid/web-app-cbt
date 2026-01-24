package out

type PasswordHasher interface {
	Compare(hash string,plain string) bool
	Hash(plain string) (string,error)
}