package out

type PasswordHasher interface {
	ComparePaswordAndHashed(hash string,plain string) bool
	GenerateHash(plain string) (string,error)
}