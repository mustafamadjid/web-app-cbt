package siswa_ujian_test

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	jawabanujianrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/jawaban_ujian"
	ujianlistrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/list"
	ujiansiswarepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/ujian_siswa"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	jawabanget "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/get_jawaban"
	jawabanhasil "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/hasil_jawaban"
	jawabansave "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban"
	listselesai "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list_soal_ujian_selesai"
	listsoal "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/soal_ujian"
	siswalist "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list"
	waktuselesai "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/waktu_selesai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSiswaUjianServices_ListAndSoalAndJawaban_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	listRepo := ujiansiswarepo.NewUjianSiswaRepo(scope.Pool(), nil)
	listSoalRepo := ujianlistrepo.NewListSoalUjianRepo(scope.Pool(), nil)
	jawabanRepo := jawabanujianrepo.NewJawabanUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	now := time.Date(2099, time.March, 2, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_IN_PROGRESS, &now, nil, testutil.Ptr(now.Add(time.Hour)))

	tingkatKelasID := int(exam.Kelas.ID)
	items, err := siswalist.NewListUjianSiswaService(listRepo).ListUjianSiswa(scope.Context(), int(exam.Siswa.ID), query.ListUjianFilter{Limit: 10, TingkatKelasID: &tingkatKelasID})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ujian.ID(exam.Ujian.ID), items[0].IdUjian)

	waktuSelesai, err := waktuselesai.NewGetWaktuSelesaiService(listRepo).GetWaktuSelesai(scope.Context(), int(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.True(t, waktuSelesai.Equal(exam.Jadwal.WaktuSelesai))

	soalSiswa, err := listsoal.NewListSoalUjianSiswaService(listSoalRepo).ListSoalUjianSiswa(scope.Context(), ujian.ID(exam.Jadwal.ID))
	require.NoError(t, err)
	require.Len(t, soalSiswa, 2)

	jawabanText := "jawaban essay integration"
	require.NoError(t, jawabansave.NewJawabanUjianService(jawabanRepo).SaveJawabanUjian(scope.Context(), ujian.ID(attempt.ID), []ujian.JawabanUjian{
		{IdSoal: ujian.ID(exam.SoalPilgan.ID), IdPilihan: testutil.Ptr(ujian.ID(exam.OpsiBenar.ID)), WaktuJawab: testutil.Ptr(now.Add(10 * time.Minute))},
		{IdSoal: ujian.ID(exam.SoalEssay.ID), JawabanEssay: &jawabanText, WaktuJawab: testutil.Ptr(now.Add(20 * time.Minute))},
	}))

	gotJawaban, err := jawabanget.NewGetJawabanUjianService(jawabanRepo).GetJawabanUjianByAttemptId(scope.Context(), ujian.ID(attempt.ID))
	require.NoError(t, err)
	require.Len(t, gotJawaban, 2)

	ownedJawaban, err := jawabanhasil.NewHasilJawabanUjianService(jawabanRepo).ListHasilJawabanUjianByAttempt(scope.Context(), int(attempt.ID))
	assert.Error(t, err)
	assert.Empty(t, ownedJawaban)

	gradedBy := user.ID(exam.Guru.ID)
	nilai := 88.0
	passed := true
	essayGraded := false
	fixtures.CreateHasilUjian(attempt.ID, &gradedBy, &nilai, &passed, &essayGraded, testutil.Ptr(now.Add(time.Hour)), exam.Jadwal.ID)

	hasil, err := jawabanhasil.NewHasilJawabanUjianService(jawabanRepo).ListHasilJawabanUjianByAttempt(scope.Context(), int(attempt.ID))
	require.NoError(t, err)
	require.Len(t, hasil, 2)
	assert.Equal(t, jawabanText, *hasil[1].JawabanSiswa.JawabanEssay)
}

func TestSiswaUjianServices_ListSelesai_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	listRepo := ujianlistrepo.NewListUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	now := time.Date(2099, time.March, 2, 9, 0, 0, 0, time.UTC)
	fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, &now, testutil.Ptr(now.Add(45*time.Minute)), testutil.Ptr(now.Add(time.Hour)))

	items, err := listselesai.NewListUjianSelesaiSiswaService(listRepo).ListUjianSelesaiSiswa(scope.Context(), int(exam.Siswa.ID))
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ujian.ID(exam.Ujian.ID), items[0].IdUjian)
}
