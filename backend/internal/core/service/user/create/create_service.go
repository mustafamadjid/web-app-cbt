package user_service

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
)

type CreateTx struct {
	txm    txout.TxManager
	hasher out.PasswordHasher
}

const (
	opCreateGuru  = "user.create_guru"
	opCreateSiswa = "user.create_siswa"
)

func NewCreateGuruService(txm txout.TxManager, hasher out.PasswordHasher) *CreateTx {
	return &CreateTx{txm: txm, hasher: hasher}
}

func NewCreateSiswaService(txm txout.TxManager, hasher out.PasswordHasher) *CreateTx {
	return &CreateTx{txm: txm, hasher: hasher}
}
