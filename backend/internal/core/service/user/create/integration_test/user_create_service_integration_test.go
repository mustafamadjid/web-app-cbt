package integration_test

import (
	"time"
	"testing"

	postgres "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	userrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	usersvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCreateService_CreateGuruAndSiswa(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	txm := postgres.NewTxManager(scope.Pool(), nil)
	hasher := bcrypt.NewHasher(4)
	guruSvc := usersvc.NewCreateGuruService(txm, hasher)
	siswaSvc := usersvc.NewCreateSiswaService(txm, hasher)

	actor := user.Actor{Role: user.ADMIN}
	emailGuru := "guru.integration@example.com"
	noHpGuru := "081234567890"
	guruRes, err := guruSvc.CreateGuru(scope.Context(), usersvc.CreateGuruCmd{
		Username:     "guru_int_01",
		Email:        &emailGuru,
		Password:     "password-guru",
		NamaLengkap:  "Guru Integration",
		JenisKelamin: "LAKI_LAKI",
		NoHp:         &noHpGuru,
		Nip:          "123456789012345678",
		Jabatan:      "Wali Kelas",
		BidangStudi:  "Matematika",
	}, actor)
	require.NoError(t, err)

	emailSiswa := "siswa.integration@example.com"
	noHpSiswa := "081234567891"
	kelas := fixtures.CreateKelas(9001)
	namaKelas := fixtures.CreateNamaKelas(kelas.ID, "X IT Integration")
	siswaRes, err := siswaSvc.CreateSiswa(scope.Context(), usersvc.CreateSiswaCmd{
		Username:      "siswa_int_01",
		Email:         &emailSiswa,
		Password:      "password-siswa",
		NamaLengkap:   "Siswa Integration",
		JenisKelamin:  "PEREMPUAN",
		NoHp:          &noHpSiswa,
		IdNamaKelas:   user.ID(namaKelas.ID),
		Nisn:          "1234567890",
		NoAbsen:       1,
		Angkatan:      2024,
		TempatLahir:   "Bandung",
		TanggalLahir:  time.Date(2007, time.January, 2, 0, 0, 0, 0, time.UTC),
	}, actor)
	require.NoError(t, err)

	userRepo := userrepo.NewUserRepo(scope.Pool(), nil)
	guruUser, err := userRepo.FindUserByID(scope.Context(), guruRes.IdPengguna)
	require.NoError(t, err)
	assert.Equal(t, "guru_int_01", guruUser.Username)
	assert.True(t, hasher.ComparePaswordAndHashed(guruUser.PasswordHashed, "password-guru"))

	siswaUser, err := userRepo.FindUserByID(scope.Context(), siswaRes.IdPengguna)
	require.NoError(t, err)
	assert.Equal(t, "siswa_int_01", siswaUser.Username)
	assert.True(t, hasher.ComparePaswordAndHashed(siswaUser.PasswordHashed, "password-siswa"))
}
