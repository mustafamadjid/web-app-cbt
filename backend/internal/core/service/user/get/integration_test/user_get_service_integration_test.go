package integration_test

import (
	"testing"

	profilgururepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_guru"
	profilsiswarepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	usersvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGetService_ListAndFind(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)

	guruUser := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guruUser.ID)
	siswaUser := fixtures.CreateUser(user.SISWA)
	kelas := fixtures.CreateKelas(9100)
	namaKelas := fixtures.CreateNamaKelas(kelas.ID, "XI IT Integration")
	fixtures.CreateSiswaProfile(siswaUser.ID, namaKelas.ID)

	guruRepo := profilgururepo.NewProfilgGuruRepo(tx, nil)
	siswaRepo := profilsiswarepo.NewProfilSiswaRepo(tx, nil)
	guruSvc := usersvc.NewGetListGuruService(guruRepo, guruRepo)
	siswaSvc := usersvc.NewGetListSiswaService(siswaRepo, siswaRepo)

	gurus, err := guruSvc.ListGuru(ctx, query.ListGuruFilter{Search: guruUser.Username, Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, gurus)

	guruProfile, err := guruSvc.FindProfilGuruByID(ctx, guruUser.ID)
	require.NoError(t, err)
	assert.Equal(t, guruUser.Username, guruProfile.Username)

	siswas, err := siswaSvc.ListSiswa(ctx, query.ListSiswaFilter{Search: siswaUser.Username, Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, siswas)

	siswaProfile, err := siswaSvc.FindProfilSiswaByID(ctx, siswaUser.ID)
	require.NoError(t, err)
	assert.Equal(t, siswaUser.Username, siswaProfile.Username)
}
