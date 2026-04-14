package gradingujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	bank_soal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	querybank "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	gradingujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGradingJawabanRepo struct {
	getRet     []ujian.JawabanUjian
	getErr     error
	getCalled  bool
	gotAttempt ujian.ID
}

func (f *fakeGradingJawabanRepo) GetJawabanUjianByAttemptId(_ context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error) {
	f.getCalled = true
	f.gotAttempt = idAttempt
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.getRet, nil
}

func (*fakeGradingJawabanRepo) SaveJawabanUjian(context.Context, ujian.ID, []ujian.JawabanUjian) error {
	return nil
}

func (*fakeGradingJawabanRepo) ListHasilJawabanUjianByIdAttempt(context.Context, ujian.ID) ([]ujian.HasilJawabanUjian, error) {
	return nil, nil
}

type fakeGradingSoalRepo struct {
	getRet    []ujian.SoalUjianSiswa
	getErr    error
	getCalled bool
	gotBankID ujian.ID
}

func (f *fakeGradingSoalRepo) GetSoalUjianByBankSoal(_ context.Context, idBankSoal ujian.ID) ([]ujian.SoalUjianSiswa, error) {
	f.getCalled = true
	f.gotBankID = idBankSoal
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.getRet, nil
}

func (*fakeGradingSoalRepo) GetSoalUjianByBankSoalForSiswa(context.Context, ujian.ID) ([]ujian.SoalUjianSiswa, bool, error) {
	return nil, false, nil
}

type fakeGradingBankSoalRepo struct {
	idRet      ujian.ID
	idErr      error
	getCalled  bool
	gotAttempt ujian.ID
}

func (*fakeGradingBankSoalRepo) GetBankSoal(context.Context, querybank.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	return nil, nil
}

func (*fakeGradingBankSoalRepo) GetBankSoalUploaded(context.Context, querybank.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	return nil, nil
}

func (*fakeGradingBankSoalRepo) GetBankSoalByGuru(context.Context, bank_soal.ID) ([]bank_soal.BankSoal, error) {
	return nil, nil
}

func (*fakeGradingBankSoalRepo) GetBankSoalById(context.Context, bank_soal.ID) (bank_soal.BankSoal, error) {
	return bank_soal.BankSoal{}, nil
}

func (*fakeGradingBankSoalRepo) CreateBankSoal(context.Context, bank_soal.BankSoal) error {
	return nil
}

func (*fakeGradingBankSoalRepo) UpdateBankSoal(context.Context, bank_soal.ID, updatepatch.UpdateBankSoalPatch) error {
	return nil
}

func (*fakeGradingBankSoalRepo) DeleteBankSoal(context.Context, bank_soal.ID) error {
	return nil
}

func (f *fakeGradingBankSoalRepo) GetIdBankSoalByAttemptId(_ context.Context, idAttempt ujian.ID) (ujian.ID, error) {
	f.getCalled = true
	f.gotAttempt = idAttempt
	if f.idErr != nil {
		return 0, f.idErr
	}
	return f.idRet, nil
}

type fakeGradingUjianRepo struct {
	idRet      ujian.ID
	idErr      error
	getCalled  bool
	gotAttempt ujian.ID
}

func (f *fakeGradingUjianRepo) GetIdUjianByAttempt(_ context.Context, idAttempt ujian.ID) (ujian.ID, error) {
	f.getCalled = true
	f.gotAttempt = idAttempt
	if f.idErr != nil {
		return 0, f.idErr
	}
	return f.idRet, nil
}

func (*fakeGradingUjianRepo) CreateUjian(context.Context, ujian.PenjadwalanUjian) error {
	return nil
}

func (*fakeGradingUjianRepo) UpdateUjian(context.Context, ujian.ID, updatepatch.UpdatePenjadwalanUjian) error {
	return nil
}

func (*fakeGradingUjianRepo) DeleteUjian(context.Context, ujian.ID) error {
	return nil
}

type fakeGradingRepo struct {
	upsertNilaiErr error
	upsertBenarErr error
	upsertSalahErr error
	essayErr       error
	statErr        error

	upsertNilaiCalled bool
	gotNilai          float64
	gotHasil          ujian.HasilUjian

	upsertBenarCalled bool
	gotBenar          []ujian.StatistikSoal

	upsertSalahCalled bool
	gotSalah          []ujian.StatistikSoal
}

func (f *fakeGradingRepo) UpsertNilaiToHasilUjian(_ context.Context, totalNilai float64, hasilUjian ujian.HasilUjian) error {
	f.upsertNilaiCalled = true
	f.gotNilai = totalNilai
	f.gotHasil = hasilUjian
	return f.upsertNilaiErr
}

func (f *fakeGradingRepo) UpsertJawabanBenarToStatistikSoal(_ context.Context, soalBenar []ujian.StatistikSoal) error {
	f.upsertBenarCalled = true
	f.gotBenar = soalBenar
	return f.upsertBenarErr
}

func (f *fakeGradingRepo) UpsertJawabanSalahToStatistikSoal(_ context.Context, soalSalah []ujian.StatistikSoal) error {
	f.upsertSalahCalled = true
	f.gotSalah = soalSalah
	return f.upsertSalahErr
}

func (f *fakeGradingRepo) UpdateAndGradingEssayUjian(context.Context, []ujian.JawabanUjian, ujian.ID) error {
	return f.essayErr
}

func (f *fakeGradingRepo) UpsertToStatistikUjian(context.Context, ujian.ID) error {
	return f.statErr
}

func newGradingService(
	jawabanRepo *fakeGradingJawabanRepo,
	soalRepo *fakeGradingSoalRepo,
	bankRepo *fakeGradingBankSoalRepo,
	ujianRepo *fakeGradingUjianRepo,
	gradingRepo *fakeGradingRepo,
) *gradingujian_service.GradingUjianService {
	return gradingujian_service.NewGradingUjianService(jawabanRepo, soalRepo, bankRepo, ujianRepo, gradingRepo)
}

func sampleJawabanAndSoal() ([]ujian.JawabanUjian, []ujian.SoalUjianSiswa) {
	pilihanBenar := ujian.ID(101)
	pilihanSalah := ujian.ID(202)

	jawaban := []ujian.JawabanUjian{
		{IdSoal: 11, IdPilihan: &pilihanBenar},
		{IdSoal: 12, IdPilihan: &pilihanSalah},
	}

	soal := []ujian.SoalUjianSiswa{
		{
			IdSoal:    11,
			BobotSoal: 10,
			OpsiJawaban: []ujian.OpsiPilganUjian{
				{IdPilihanGanda: 101, IsBenar: true},
				{IdPilihanGanda: 102, IsBenar: false},
			},
		},
		{
			IdSoal:    12,
			BobotSoal: 20,
			OpsiJawaban: []ujian.OpsiPilganUjian{
				{IdPilihanGanda: 201, IsBenar: true},
				{IdPilihanGanda: 202, IsBenar: false},
			},
		},
	}

	return jawaban, soal
}

func TestGradingUjianService_TotalScore_BranchCoverage(t *testing.T) {
	t.Parallel()

	svc := newGradingService(&fakeGradingJawabanRepo{}, &fakeGradingSoalRepo{}, &fakeGradingBankSoalRepo{}, &fakeGradingUjianRepo{}, &fakeGradingRepo{})

	t.Run("branch 1 -> soal ujian kosong", func(t *testing.T) {
		total, benar, salah, err := svc.TotalScore([]ujian.JawabanUjian{}, nil, 99)

		assert.ErrorIs(t, err, coreerror.ErrArrayHasNoElement)
		assert.Zero(t, total)
		assert.Nil(t, benar)
		assert.Nil(t, salah)
	})

	t.Run("branch 2 -> hitung total nilai benar dan salah", func(t *testing.T) {
		jawaban, soal := sampleJawabanAndSoal()

		total, benar, salah, err := svc.TotalScore(jawaban, soal, 99)

		require.NoError(t, err)
		assert.Equal(t, 10.0, total)
		assert.Equal(t, []ujian.StatistikSoal{{IDSoal: 11, IDUjian: 99}}, benar)
		assert.Equal(t, []ujian.StatistikSoal{{IDSoal: 12, IDUjian: 99}}, salah)
	})

	t.Run("branch 3 -> jawaban tanpa pilihan diabaikan", func(t *testing.T) {
		_, soal := sampleJawabanAndSoal()
		jawaban := []ujian.JawabanUjian{{IdSoal: 11}}

		total, benar, salah, err := svc.TotalScore(jawaban, soal, 99)

		require.NoError(t, err)
		assert.Zero(t, total)
		assert.Empty(t, benar)
		assert.Empty(t, salah)
	})
}

func TestGradingUjianService_UpsertingToStatistikSoal_BranchCoverage(t *testing.T) {
	t.Parallel()

	t.Run("branch 1 -> upsert jawaban benar gagal", func(t *testing.T) {
		repoErr := errors.New("upsert benar gagal")
		repo := &fakeGradingRepo{upsertBenarErr: repoErr}
		svc := newGradingService(&fakeGradingJawabanRepo{}, &fakeGradingSoalRepo{}, &fakeGradingBankSoalRepo{}, &fakeGradingUjianRepo{}, repo)

		err := svc.UpsertingToStatistikSoal(context.Background(), []ujian.StatistikSoal{{IDSoal: 1}}, nil)

		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.upsertBenarCalled)
		assert.False(t, repo.upsertSalahCalled)
	})

	t.Run("branch 2 -> upsert jawaban salah gagal", func(t *testing.T) {
		repoErr := errors.New("upsert salah gagal")
		repo := &fakeGradingRepo{upsertSalahErr: repoErr}
		svc := newGradingService(&fakeGradingJawabanRepo{}, &fakeGradingSoalRepo{}, &fakeGradingBankSoalRepo{}, &fakeGradingUjianRepo{}, repo)

		err := svc.UpsertingToStatistikSoal(context.Background(), nil, []ujian.StatistikSoal{{IDSoal: 2}})

		assert.ErrorIs(t, err, repoErr)
		assert.False(t, repo.upsertBenarCalled)
		assert.True(t, repo.upsertSalahCalled)
	})

	t.Run("branch 3 -> statistik kosong tidak memanggil repo", func(t *testing.T) {
		repo := &fakeGradingRepo{}
		svc := newGradingService(&fakeGradingJawabanRepo{}, &fakeGradingSoalRepo{}, &fakeGradingBankSoalRepo{}, &fakeGradingUjianRepo{}, repo)

		err := svc.UpsertingToStatistikSoal(context.Background(), nil, nil)

		assert.NoError(t, err)
		assert.False(t, repo.upsertBenarCalled)
		assert.False(t, repo.upsertSalahCalled)
	})
}

func TestGradingUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	jawaban, soal := sampleJawabanAndSoal()
	jawabanRepoErr := errors.New("jawaban error")
	bankRepoErr := errors.New("bank error")
	soalRepoErr := errors.New("soal error")
	ujianRepoErr := errors.New("ujian error")
	upsertNilaiErr := errors.New("nilai error")
	upsertBenarErr := errors.New("benar error")
	upsertSalahErr := errors.New("salah error")

	tests := []struct {
		name         string
		idAttempt    int
		jawabanRepo  *fakeGradingJawabanRepo
		soalRepo     *fakeGradingSoalRepo
		bankRepo     *fakeGradingBankSoalRepo
		ujianRepo    *fakeGradingUjianRepo
		gradingRepo  *fakeGradingRepo
		wantErr      error
		assertResult func(t *testing.T, repo *fakeGradingRepo)
	}{
		{
			name:        "path 1 -> id attempt tidak valid",
			idAttempt:   0,
			jawabanRepo: &fakeGradingJawabanRepo{},
			soalRepo:    &fakeGradingSoalRepo{},
			bankRepo:    &fakeGradingBankSoalRepo{},
			ujianRepo:   &fakeGradingUjianRepo{},
			gradingRepo: &fakeGradingRepo{},
			wantErr:     coreerror.ErrMissingId,
		},
		{
			name:        "path 2 -> get jawaban ujian gagal",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getErr: jawabanRepoErr},
			soalRepo:    &fakeGradingSoalRepo{},
			bankRepo:    &fakeGradingBankSoalRepo{},
			ujianRepo:   &fakeGradingUjianRepo{},
			gradingRepo: &fakeGradingRepo{},
			wantErr:     jawabanRepoErr,
		},
		{
			name:        "path 3 -> get id bank soal gagal",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{},
			bankRepo:    &fakeGradingBankSoalRepo{idErr: bankRepoErr},
			ujianRepo:   &fakeGradingUjianRepo{},
			gradingRepo: &fakeGradingRepo{},
			wantErr:     bankRepoErr,
		},
		{
			name:        "path 4 -> get soal ujian gagal",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{getErr: soalRepoErr},
			bankRepo:    &fakeGradingBankSoalRepo{idRet: 31},
			ujianRepo:   &fakeGradingUjianRepo{},
			gradingRepo: &fakeGradingRepo{},
			wantErr:     soalRepoErr,
		},
		{
			name:        "path 5 -> get id ujian gagal",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{getRet: soal},
			bankRepo:    &fakeGradingBankSoalRepo{idRet: 31},
			ujianRepo:   &fakeGradingUjianRepo{idErr: ujianRepoErr},
			gradingRepo: &fakeGradingRepo{},
			wantErr:     ujianRepoErr,
		},
		{
			name:        "path 6 -> total score gagal karena soal kosong",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{getRet: []ujian.SoalUjianSiswa{}},
			bankRepo:    &fakeGradingBankSoalRepo{idRet: 31},
			ujianRepo:   &fakeGradingUjianRepo{idRet: 41},
			gradingRepo: &fakeGradingRepo{},
			wantErr:     coreerror.ErrArrayHasNoElement,
		},
		{
			name:        "path 7 -> upsert nilai hasil ujian gagal",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{getRet: soal},
			bankRepo:    &fakeGradingBankSoalRepo{idRet: 31},
			ujianRepo:   &fakeGradingUjianRepo{idRet: 41},
			gradingRepo: &fakeGradingRepo{upsertNilaiErr: upsertNilaiErr},
			wantErr:     upsertNilaiErr,
		},
		{
			name:        "path 8 -> upsert statistik benar gagal",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{getRet: soal},
			bankRepo:    &fakeGradingBankSoalRepo{idRet: 31},
			ujianRepo:   &fakeGradingUjianRepo{idRet: 41},
			gradingRepo: &fakeGradingRepo{upsertBenarErr: upsertBenarErr},
			wantErr:     upsertBenarErr,
		},
		{
			name:        "path 9 -> upsert statistik salah gagal",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{getRet: soal},
			bankRepo:    &fakeGradingBankSoalRepo{idRet: 31},
			ujianRepo:   &fakeGradingUjianRepo{idRet: 41},
			gradingRepo: &fakeGradingRepo{upsertSalahErr: upsertSalahErr},
			wantErr:     upsertSalahErr,
		},
		{
			name:        "path 10 -> berhasil grading ujian pilihan ganda",
			idAttempt:   21,
			jawabanRepo: &fakeGradingJawabanRepo{getRet: jawaban},
			soalRepo:    &fakeGradingSoalRepo{getRet: soal},
			bankRepo:    &fakeGradingBankSoalRepo{idRet: 31},
			ujianRepo:   &fakeGradingUjianRepo{idRet: 41},
			gradingRepo: &fakeGradingRepo{},
			assertResult: func(t *testing.T, repo *fakeGradingRepo) {
				t.Helper()
				assert.True(t, repo.upsertNilaiCalled)
				assert.Equal(t, 10.0, repo.gotNilai)
				assert.Equal(t, ujian.ID(21), repo.gotHasil.IdAttempt)
				assert.Equal(t, []ujian.StatistikSoal{{IDSoal: 11, IDUjian: 41}}, repo.gotBenar)
				assert.Equal(t, []ujian.StatistikSoal{{IDSoal: 12, IDUjian: 41}}, repo.gotSalah)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := newGradingService(tc.jawabanRepo, tc.soalRepo, tc.bankRepo, tc.ujianRepo, tc.gradingRepo)
			err := svc.GradingUjianPilgan(ctx, tc.idAttempt)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			if tc.assertResult != nil {
				tc.assertResult(t, tc.gradingRepo)
			}
		})
	}
}
