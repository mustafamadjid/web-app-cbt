package user_service_test

// ===== fake hasher =====

type FakeHasher struct {
	Hash      string
	Err       error
	Called    bool
	LastPlain string
}

func (fakeHash *FakeHasher) ComparePaswordAndHashed(hash string, plain string) bool {
	return hash == plain
}

func (fakeHash *FakeHasher) GenerateHash(plain string) (string, error) {
	fakeHash.Called = true
	fakeHash.LastPlain = plain
	if fakeHash.Err != nil {
		return "", fakeHash.Err
	}
	return fakeHash.Hash, nil
}
