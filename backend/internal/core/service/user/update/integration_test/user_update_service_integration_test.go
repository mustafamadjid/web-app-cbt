package integration_test

import (
	"testing"
	"time"

	postgres "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	sessionrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/session"
	userrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	deletefilesvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/delete_file_system"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	usersvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserUpdateService_UpdateGuruAndSiswa(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	txm := postgres.NewTxManager(scope.Pool(), nil)
	userRepo := userrepo.NewUserRepo(scope.Pool(), nil)
	sessionRepo := sessionrepo.NewSessionRepo(scope.Pool(), nil)
	deleteFileSvc := deletefilesvc.NewDeleteFileService(t.TempDir())
	svc := usersvc.NewUpdateUserService(txm, sessionRepo, deleteFileSvc, userRepo)

	hasher := bcrypt.NewHasher(4)
	guruEmail := user.Email("guru-update@example.com")
	guruPass, err := hasher.GenerateHash("password-guru")
	require.NoError(t, err)
	guruID, err := userRepo.CreateUser(scope.Context(), user.Pengguna{
		Username:       "guru_update_01",
		Email:          &guruEmail,
		PasswordHashed: guruPass,
		NamaLengkap:    "Guru Update",
		JenisKelamin:   "LAKI_LAKI",
		Role:           user.GURU,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)
	fixtures.CreateGuruProfile(guruID)

	// create active session so update can revoke it
	sessionGuruID, err := sessionRepo.CreateSession(scope.Context(), guruID, user.GURU, time.Now().Add(time.Hour))
	require.NoError(t, err)

	updatedEmail := "guru-update-new@example.com"
	updatedName := "Guru Update Baru"
	require.NoError(t, svc.UpdateGuru(scope.Context(), usersvc.UpdateGuruCmd{
		IdPengguna:  guruID,
		Email:       &updatedEmail,
		NamaLengkap: &updatedName,
	}, user.Actor{Role: user.ADMIN}))

	storedGuru, err := userRepo.FindUserByID(scope.Context(), guruID)
	require.NoError(t, err)
	assert.Equal(t, updatedName, storedGuru.NamaLengkap)
	require.NotNil(t, storedGuru.Email)
	assert.Equal(t, updatedEmail, string(*storedGuru.Email))
	sessionGuru, err := sessionRepo.GetSession(scope.Context(), sessionGuruID)
	require.NoError(t, err)
	assert.True(t, sessionGuru.Revoked)

	siswaEmail := user.Email("siswa-update@example.com")
	siswaPass, err := hasher.GenerateHash("password-siswa")
	require.NoError(t, err)
	kelas := fixtures.CreateKelas(9200)
	namaKelas := fixtures.CreateNamaKelas(kelas.ID, "XII Update")
	siswaID, err := userRepo.CreateUser(scope.Context(), user.Pengguna{
		Username:       "siswa_update_01",
		Email:          &siswaEmail,
		PasswordHashed: siswaPass,
		NamaLengkap:    "Siswa Update",
		JenisKelamin:   "PEREMPUAN",
		Role:           user.SISWA,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)
	fixtures.CreateSiswaProfile(siswaID, namaKelas.ID)
	sessionSiswaID, err := sessionRepo.CreateSession(scope.Context(), siswaID, user.SISWA, time.Now().Add(time.Hour))
	require.NoError(t, err)

	newSiswaName := "Siswa Update Baru"
	newAbsen := 2
	require.NoError(t, svc.UpdateSiswa(scope.Context(), usersvc.UpdateSiswaCmd{
		IdPengguna:  siswaID,
		NamaLengkap: &newSiswaName,
		NoAbsen:     &newAbsen,
	}, user.Actor{Role: user.ADMIN}))

	storedSiswa, err := userRepo.FindUserByID(scope.Context(), siswaID)
	require.NoError(t, err)
	assert.Equal(t, newSiswaName, storedSiswa.NamaLengkap)
	sessionSiswa, err := sessionRepo.GetSession(scope.Context(), sessionSiswaID)
	require.NoError(t, err)
	assert.True(t, sessionSiswa.Revoked)
}
