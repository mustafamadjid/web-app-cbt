package bcrypt

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func NewHasher(cost int) *Hasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &Hasher{cost: cost}
}

func (h *Hasher) GenerateHash(plain string)(string,error){
	b, error := bcrypt.GenerateFromPassword([]byte(plain),h.cost)
	if error != nil {
		return "", error
	}
	return string(b), nil
}

func (h *Hasher) ComparePaswordAndHashed(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}