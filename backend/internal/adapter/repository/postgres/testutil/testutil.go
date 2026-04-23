package testutil

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgcontract "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	banksoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

var uniqueCounter atomic.Uint64

type CommittedFixtureScope struct {
	t        *testing.T
	ctx      context.Context
	pool     *pgxpool.Pool
	cleanups []func()
}

type Fixtures struct {
	t           *testing.T
	ctx         context.Context
	q           pgcontract.Executor
	addCleanup  func(fn func())
}

type CreatedUser struct {
	ID          user.ID
	Username    string
	Email       string
	NoHP        string
	Role        user.Role
	NamaLengkap string
}

type CreatedKelas struct {
	ID           int64
	TingkatKelas int
}

type CreatedNamaKelas struct {
	ID          int64
	IDKelas     int64
	NamaKelas   string
}

type CreatedMapel struct {
	ID        int64
	IDKelas   int64
	KodeMapel string
	NamaMapel string
}

type CreatedBankSoal struct {
	ID         int64
	IDMapel    int64
	IDKelas    int64
	IDPengguna int64
	Nama       string
}

type CreatedBankSoalVersion struct {
	ID         int64
	IDBankSoal int64
	VersionNo  int
	Status     string
}

type CreatedSoal struct {
	ID        int64
	IDVersion int64
	TipeSoal  string
	NoUrut    int
}

type CreatedOpsi struct {
	ID      int64
	IDSoal  int64
	Isi     string
	IsBenar bool
}

type CreatedSesi struct {
	ID       int64
	KodeSesi string
	NamaSesi string
}

type CreatedRuangUjian struct {
	ID         int64
	KodeRuang  string
	NamaRuang  string
}

type CreatedPengumuman struct {
	ID       int64
	Judul    string
	TanggalRilis string
	TanggalSelesai string
}

type CreatedUjian struct {
	ID         int64
	IDBankSoal int64
	IDKelas    int64
	IDNamaKelas *int64
	IDGuru     int64
}

type CreatedJadwalUjian struct {
	ID          int64
	IDUjian     int64
	IDSesi      int64
	IDRuangan   int64
	IDPengawas  int64
	Token       string
	WaktuMulai  time.Time
	WaktuSelesai time.Time
}

type CreatedPeserta struct {
	ID            int64
	IDJadwalUjian int64
	IDSiswa       int64
}

type CreatedAttempt struct {
	ID             int64
	IDPesertaUjian int64
	Status         ujian.StatusAttempt
}

type CreatedSession struct {
	ID         string
	IDPengguna int64
	Role       user.Role
}

type CreatedImportJob struct {
	ID         int64
	IDBankSoal int64
	IDPengguna int64
	Status     importsoal.JobStatus
}

type CreatedGradingJob struct {
	ID        int64
	IDAttempt int64
	Status    ujian.JobStatus
}

type ExamFixture struct {
	Kelas       CreatedKelas
	NamaKelas   CreatedNamaKelas
	Guru        CreatedUser
	Pengawas    CreatedUser
	Siswa       CreatedUser
	Mapel       CreatedMapel
	BankSoal    CreatedBankSoal
	Version     CreatedBankSoalVersion
	SoalPilgan  CreatedSoal
	SoalEssay   CreatedSoal
	OpsiBenar   CreatedOpsi
	OpsiSalah   CreatedOpsi
	Sesi        CreatedSesi
	Ruang       CreatedRuangUjian
	Ujian       CreatedUjian
	Jadwal      CreatedJadwalUjian
	Peserta     CreatedPeserta
}

func BeginRollbackTx(t *testing.T) (context.Context, pgx.Tx) {
	t.Helper()

	ctx := context.Background()
	pool := OpenPool(t, ctx)
	tx, err := pool.Begin(ctx)
	if err != nil {
		pool.Close()
		t.Fatalf("begin transaction: %v", err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(context.Background()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			t.Errorf("rollback transaction: %v", err)
		}
		pool.Close()
	})

	return ctx, tx
}

func NewCommittedFixtureScope(t *testing.T) *CommittedFixtureScope {
	t.Helper()

	ctx := context.Background()
	scope := &CommittedFixtureScope{
		t:    t,
		ctx:  ctx,
		pool: OpenPool(t, ctx),
	}
	t.Cleanup(scope.cleanup)
	return scope
}

func (s *CommittedFixtureScope) Context() context.Context {
	return s.ctx
}

func (s *CommittedFixtureScope) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *CommittedFixtureScope) Fixtures() *Fixtures {
	return &Fixtures{
		t:   s.t,
		ctx: s.ctx,
		q:   s.pool,
		addCleanup: func(fn func()) {
			s.cleanups = append(s.cleanups, fn)
		},
	}
}

func (s *CommittedFixtureScope) AddCleanupQuery(query string, args ...any) {
	s.cleanups = append(s.cleanups, func() {
		if _, err := s.pool.Exec(context.Background(), query, args...); err != nil {
			s.t.Errorf("cleanup failed for %q: %v", query, err)
		}
	})
}

func (s *CommittedFixtureScope) cleanup() {
	for i := len(s.cleanups) - 1; i >= 0; i-- {
		s.cleanups[i]()
	}
	s.pool.Close()
}

func NewFixtures(t *testing.T, ctx context.Context, q pgcontract.Executor) *Fixtures {
	t.Helper()
	return &Fixtures{t: t, ctx: ctx, q: q}
}

func OpenPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	dbURL := strings.TrimSpace(os.Getenv("POSTGRES_TEST_DBURL"))
	if dbURL == "" {
		t.Skip("POSTGRES_TEST_DBURL is not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("create pgx pool: %v", err)
	}

	return pool
}

func UniqueSuffix(prefix string) string {
	value := uniqueCounter.Add(1)
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), value)
}

func uniqueNumeric() int64 {
	return time.Now().UnixNano()%1_000_000_000_000 + int64(uniqueCounter.Add(1))
}

func Ptr[T any](value T) *T {
	return &value
}

func (f *Fixtures) LookupSeedUserID(username string) user.ID {
	f.t.Helper()

	var id user.ID
	err := f.q.QueryRow(f.ctx, `
		SELECT id_pengguna
		FROM pengguna
		WHERE username = $1
	`, username).Scan(&id)
	if err != nil {
		f.t.Fatalf("lookup seed pengguna %q: %v", username, err)
	}
	return id
}

func (f *Fixtures) LookupKelasIDByTingkat(tingkat int) int64 {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		SELECT id_kelas
		FROM kelas
		WHERE tingkat_kelas = $1
	`, tingkat).Scan(&id)
	if err != nil {
		f.t.Fatalf("lookup kelas by tingkat %d: %v", tingkat, err)
	}
	return id
}

func (f *Fixtures) LookupNamaKelasID(idKelas int64, namaKelas string) int64 {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		SELECT id_nama_kelas
		FROM nama_kelas
		WHERE id_kelas = $1 AND nama_kelas = $2
	`, idKelas, namaKelas).Scan(&id)
	if err != nil {
		f.t.Fatalf("lookup nama_kelas %q: %v", namaKelas, err)
	}
	return id
}

func (f *Fixtures) CreateUser(role user.Role) CreatedUser {
	f.t.Helper()

	n := uniqueNumeric()
	username := fmt.Sprintf("u%018d", n%1_000_000_000_000_000_000)
	email := fmt.Sprintf("u%018d@example.com", n%1_000_000_000_000_000_000)
	noHP := fmt.Sprintf("08%010d", n%10_000_000_000)
	nama := "Integration " + username

	var id user.ID
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO pengguna (
			foto,
			nama_lengkap,
			jenis_kelamin,
			username,
			password,
			email,
			no_hp,
			id_role,
			status_akun
		)
		VALUES (
			NULL,
			$1,
			1,
			$2,
			$3,
			$4,
			$5,
			(SELECT id_role FROM role WHERE nama_role = $6),
			'AKTIF'
		)
		RETURNING id_pengguna
	`, nama, username, "hashed-"+username, email, noHP, string(role)).Scan(&id)
	if err != nil {
		f.t.Fatalf("create user %q: %v", username, err)
	}

	f.registerCleanup(`DELETE FROM pengguna WHERE id_pengguna = $1`, id)

	return CreatedUser{
		ID:          id,
		Username:    username,
		Email:       email,
		NoHP:        noHP,
		Role:        role,
		NamaLengkap: nama,
	}
}

func (f *Fixtures) CreateGuruProfile(idPengguna user.ID) int64 {
	f.t.Helper()

	nip := fmt.Sprintf("19800101%010d", uniqueNumeric()%10_000_000_000)
	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO profil_guru (id_pengguna, nip, jabatan, bidang_studi)
		VALUES ($1, $2, $3, $4)
		RETURNING id_guru
	`, idPengguna, nip, "Guru", "Matematika").Scan(&id)
	if err != nil {
		f.t.Fatalf("create profil guru: %v", err)
	}

	f.registerCleanup(`DELETE FROM profil_guru WHERE id_guru = $1`, id)
	return id
}

func (f *Fixtures) CreateKelas(tingkat int) CreatedKelas {
	f.t.Helper()

	tingkat = 1000 + int(uniqueNumeric()%20000)

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO kelas (tingkat_kelas)
		VALUES ($1)
		RETURNING id_kelas
	`, tingkat).Scan(&id)
	if err != nil {
		f.t.Fatalf("create kelas %d: %v", tingkat, err)
	}

	f.registerCleanup(`DELETE FROM kelas WHERE id_kelas = $1`, id)
	return CreatedKelas{ID: id, TingkatKelas: tingkat}
}

func (f *Fixtures) CreateNamaKelas(idKelas int64, nama string) CreatedNamaKelas {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO nama_kelas (id_kelas, nama_kelas)
		VALUES ($1, $2)
		RETURNING id_nama_kelas
	`, idKelas, nama).Scan(&id)
	if err != nil {
		f.t.Fatalf("create nama kelas %q: %v", nama, err)
	}

	f.registerCleanup(`DELETE FROM nama_kelas WHERE id_nama_kelas = $1`, id)
	return CreatedNamaKelas{ID: id, IDKelas: idKelas, NamaKelas: nama}
}

func (f *Fixtures) CreateSiswaProfile(idPengguna user.ID, idNamaKelas int64) int64 {
	f.t.Helper()

	nisn := fmt.Sprintf("%010d", uniqueNumeric()%10_000_000_000)
	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO profil_siswa (
			id_pengguna,
			id_nama_kelas,
			nisn,
			no_absen,
			angkatan,
			tempat_lahir,
			tanggal_lahir
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_siswa
	`, idPengguna, idNamaKelas, nisn, 1, 2024, "Bandung", time.Date(2007, time.January, 2, 0, 0, 0, 0, time.UTC)).Scan(&id)
	if err != nil {
		f.t.Fatalf("create profil siswa: %v", err)
	}

	f.registerCleanup(`DELETE FROM profil_siswa WHERE id_siswa = $1`, id)
	return id
}

func (f *Fixtures) CreateMapel(idKelas int64) CreatedMapel {
	f.t.Helper()

	suffix := uniqueNumeric()
	kode := fmt.Sprintf("MP%08d", suffix%100_000_000)
	nama := "Mapel " + fmt.Sprint(suffix)
	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO mata_pelajaran (id_kelas, kode_mapel, nama_mapel, deskripsi)
		VALUES ($1, $2, $3, $4)
		RETURNING id_mapel
	`, idKelas, kode, nama, "deskripsi "+nama).Scan(&id)
	if err != nil {
		f.t.Fatalf("create mapel %q: %v", nama, err)
	}

	f.registerCleanup(`DELETE FROM mata_pelajaran WHERE id_mapel = $1`, id)
	return CreatedMapel{ID: id, IDKelas: idKelas, KodeMapel: kode, NamaMapel: nama}
}

func (f *Fixtures) CreateBankSoal(idMapel, idKelas int64, idPengguna user.ID) CreatedBankSoal {
	f.t.Helper()

	nama := "Bank " + UniqueSuffix("soal")
	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO bank_soal (id_mapel, id_kelas, id_pengguna, nama_bank_soal, deskripsi, materi)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_bank_soal
	`, idMapel, idKelas, idPengguna, nama, "deskripsi "+nama, "materi "+nama).Scan(&id)
	if err != nil {
		f.t.Fatalf("create bank soal %q: %v", nama, err)
	}

	f.registerCleanup(`DELETE FROM bank_soal WHERE id_bank_soal = $1`, id)
	return CreatedBankSoal{
		ID:         id,
		IDMapel:    idMapel,
		IDKelas:    idKelas,
		IDPengguna: int64(idPengguna),
		Nama:       nama,
	}
}

func (f *Fixtures) CreateBankSoalVersion(idBankSoal int64, createdBy user.ID, status string, versionNo int) CreatedBankSoalVersion {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO bank_soal_version (id_bank_soal, version_no, status, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id_bank_soal_version
	`, idBankSoal, versionNo, status, createdBy).Scan(&id)
	if err != nil {
		f.t.Fatalf("create bank soal version: %v", err)
	}

	f.registerCleanup(`DELETE FROM bank_soal_version WHERE id_bank_soal_version = $1`, id)
	return CreatedBankSoalVersion{ID: id, IDBankSoal: idBankSoal, VersionNo: versionNo, Status: status}
}

func (f *Fixtures) SetActiveBankSoalVersion(idBankSoal, idVersion int64) {
	f.t.Helper()

	if _, err := f.q.Exec(f.ctx, `
		UPDATE bank_soal
		SET id_bank_soal_version_aktif = $1
		WHERE id_bank_soal = $2
	`, idVersion, idBankSoal); err != nil {
		f.t.Fatalf("set active bank soal version: %v", err)
	}
}

func (f *Fixtures) CreateSoal(idVersion int64, tipe, pertanyaan string, bobot float64, noUrut int) CreatedSoal {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO isi_soal (id_bank_soal_version, tipe_soal, pertanyaan, gambar, bobot_soal, no_urut_soal, created_at, updated_at)
		VALUES ($1, $2, $3, NULL, $4, $5, now(), now())
		RETURNING id_soal
	`, idVersion, tipe, pertanyaan, bobot, noUrut).Scan(&id)
	if err != nil {
		f.t.Fatalf("create soal %q: %v", pertanyaan, err)
	}

	f.registerCleanup(`DELETE FROM isi_soal WHERE id_soal = $1`, id)
	return CreatedSoal{ID: id, IDVersion: idVersion, TipeSoal: tipe, NoUrut: noUrut}
}

func (f *Fixtures) CreateOpsi(idSoal int64, isi string, isBenar bool) CreatedOpsi {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO opsi_pilihan_ganda (id_soal, isi_pilihan, is_benar)
		VALUES ($1, $2, $3)
		RETURNING id_pilihan_ganda
	`, idSoal, isi, isBenar).Scan(&id)
	if err != nil {
		f.t.Fatalf("create opsi %q: %v", isi, err)
	}

	f.registerCleanup(`DELETE FROM opsi_pilihan_ganda WHERE id_pilihan_ganda = $1`, id)
	return CreatedOpsi{ID: id, IDSoal: idSoal, Isi: isi, IsBenar: isBenar}
}

func (f *Fixtures) CreateSesi() CreatedSesi {
	f.t.Helper()

	suffix := uniqueNumeric()
	kode := fmt.Sprintf("SESI%08d", suffix%100_000_000)
	nama := "Sesi " + fmt.Sprint(suffix)

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO sesi_ujian (kode_sesi, nama_sesi)
		VALUES ($1, $2)
		RETURNING id_sesi
	`, kode, nama).Scan(&id)
	if err != nil {
		f.t.Fatalf("create sesi: %v", err)
	}

	f.registerCleanup(`DELETE FROM sesi_ujian WHERE id_sesi = $1`, id)
	return CreatedSesi{ID: id, KodeSesi: kode, NamaSesi: nama}
}

func (f *Fixtures) CreateRuangUjian() CreatedRuangUjian {
	f.t.Helper()

	suffix := uniqueNumeric()
	kode := fmt.Sprintf("RUANG%08d", suffix%100_000_000)
	nama := "Ruang " + fmt.Sprint(suffix)

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO ruang_ujian (nama_ruangan, kode_ruang)
		VALUES ($1, $2)
		RETURNING id_ruangan
	`, nama, kode).Scan(&id)
	if err != nil {
		f.t.Fatalf("create ruang ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM ruang_ujian WHERE id_ruangan = $1`, id)
	return CreatedRuangUjian{ID: id, KodeRuang: kode, NamaRuang: nama}
}

func (f *Fixtures) CreatePengumuman(idPengguna user.ID, tanggalRilis, tanggalSelesai time.Time) CreatedPengumuman {
	f.t.Helper()

	judul := "Pengumuman " + UniqueSuffix("judul")
	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO pengumuman (
			id_pengguna,
			judul_pengumuman,
			isi_pengumuman,
			tanggal_rilis_pengumuman,
			tanggal_selesai_pengumuman,
			dokumen_pengumuman
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_pengumuman
	`, idPengguna, judul, "isi "+judul, tanggalRilis.Format("2006-01-02"), tanggalSelesai.Format("2006-01-02"), "").Scan(&id)
	if err != nil {
		f.t.Fatalf("create pengumuman: %v", err)
	}

	f.registerCleanup(`DELETE FROM pengumuman WHERE id_pengumuman = $1`, id)
	return CreatedPengumuman{
		ID:              id,
		Judul:           judul,
		TanggalRilis:    tanggalRilis.Format("2006-01-02"),
		TanggalSelesai:  tanggalSelesai.Format("2006-01-02"),
	}
}

func (f *Fixtures) CreateUjian(
	idBankSoal int64,
	idKelas int64,
	idNamaKelas *int64,
	idGuru user.ID,
	nama string,
) CreatedUjian {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO ujian (
			id_bank_soal,
			id_kelas,
			id_nama_kelas,
			id_guru,
			nama_ujian,
			deskripsi_ujian,
			acak_soal
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_ujian
	`, idBankSoal, idKelas, idNamaKelas, idGuru, nama, "deskripsi "+nama, false).Scan(&id)
	if err != nil {
		f.t.Fatalf("create ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM ujian WHERE id_ujian = $1`, id)
	return CreatedUjian{
		ID:          id,
		IDBankSoal:  idBankSoal,
		IDKelas:     idKelas,
		IDNamaKelas: idNamaKelas,
		IDGuru:      int64(idGuru),
	}
}

func (f *Fixtures) CreateJadwalUjian(
	idUjian int64,
	idSesi int64,
	idRuangan int64,
	idPengawas user.ID,
	tanggalUjian time.Time,
	waktuMulai time.Time,
	waktuSelesai time.Time,
	status ujian.StatusUjian,
) CreatedJadwalUjian {
	f.t.Helper()

	token := "TKN" + fmt.Sprint(uniqueCounter.Add(1))
	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO jadwal_ujian (
			id_ujian,
			id_sesi,
			id_ruangan,
			id_pengawas,
			tanggal_ujian,
			waktu_mulai,
			waktu_selesai,
			token,
			status_ujian
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id_jadwal_ujian
	`, idUjian, idSesi, idRuangan, idPengawas, tanggalUjian, waktuMulai, waktuSelesai, token, status).Scan(&id)
	if err != nil {
		f.t.Fatalf("create jadwal ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM jadwal_ujian WHERE id_jadwal_ujian = $1`, id)
	return CreatedJadwalUjian{
		ID:           id,
		IDUjian:      idUjian,
		IDSesi:       idSesi,
		IDRuangan:    idRuangan,
		IDPengawas:   int64(idPengawas),
		Token:        token,
		WaktuMulai:   waktuMulai,
		WaktuSelesai: waktuSelesai,
	}
}

func (f *Fixtures) CreatePeserta(idJadwal int64, idSiswa user.ID) CreatedPeserta {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO peserta_ujian (id_jadwal_ujian, id_siswa)
		VALUES ($1, $2)
		RETURNING id_peserta_ujian
	`, idJadwal, idSiswa).Scan(&id)
	if err != nil {
		f.t.Fatalf("create peserta ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM peserta_ujian WHERE id_peserta_ujian = $1`, id)
	return CreatedPeserta{ID: id, IDJadwalUjian: idJadwal, IDSiswa: int64(idSiswa)}
}

func (f *Fixtures) CreateAttempt(idPeserta int64, status ujian.StatusAttempt, waktuMulai, waktuSubmit, deadline *time.Time) CreatedAttempt {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO attempt_ujian (id_peserta_ujian, status_attempt, waktu_mulai, waktu_submit, deadline_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_attempt
	`, idPeserta, status, waktuMulai, waktuSubmit, deadline).Scan(&id)
	if err != nil {
		f.t.Fatalf("create attempt ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM attempt_ujian WHERE id_attempt = $1`, id)
	return CreatedAttempt{ID: id, IDPesertaUjian: idPeserta, Status: status}
}

func (f *Fixtures) CreateJawaban(idAttempt, idSoal int64, idPilihan *int64, jawabanEssay *string, waktuJawab *time.Time, essayIsBenar *bool) int64 {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO jawaban_ujian_siswa (id_attempt, id_soal, id_pilihan, jawaban_essay, waktu_jawab, essay_is_benar)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_jawaban
	`, idAttempt, idSoal, idPilihan, jawabanEssay, waktuJawab, essayIsBenar).Scan(&id)
	if err != nil {
		f.t.Fatalf("create jawaban ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM jawaban_ujian_siswa WHERE id_jawaban = $1`, id)
	return id
}

func (f *Fixtures) CreateHasilUjian(idAttempt int64, gradedBy *user.ID, nilai *float64, passed *bool, essayGraded *bool, gradedAt *time.Time, idJadwal int64) int64 {
	f.t.Helper()

	if essayGraded == nil {
		essayGraded = Ptr(false)
	}

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO hasil_ujian (id_attempt, graded_by, nilai_akhir, passed, essay_graded, graded_at, id_jadwal_ujian)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_hasil
	`, idAttempt, gradedBy, nilai, passed, essayGraded, gradedAt, idJadwal).Scan(&id)
	if err != nil {
		f.t.Fatalf("create hasil ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM hasil_ujian WHERE id_hasil = $1`, id)
	return id
}

func (f *Fixtures) CreateStatistikUjian(idJadwal int64, tertinggi, terendah, rata float64, total int) int64 {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO statistik_ujian (id_jadwal_ujian, nilai_tertinggi, nilai_terendah, nilai_rata_rata, total_peserta_ujian)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_statistik_ujian
	`, idJadwal, tertinggi, terendah, rata, total).Scan(&id)
	if err != nil {
		f.t.Fatalf("create statistik ujian: %v", err)
	}

	f.registerCleanup(`DELETE FROM statistik_ujian WHERE id_statistik_ujian = $1`, id)
	return id
}

func (f *Fixtures) CreateStatistikSoal(idSoal, idUjian int64, benar, salah int) int64 {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO statistik_soal (id_soal, id_ujian, jumlah_jawaban_benar, jumlah_jawaban_salah)
		VALUES ($1, $2, $3, $4)
		RETURNING id_statistik_soal
	`, idSoal, idUjian, benar, salah).Scan(&id)
	if err != nil {
		f.t.Fatalf("create statistik soal: %v", err)
	}

	f.registerCleanup(`DELETE FROM statistik_soal WHERE id_statistik_soal = $1`, id)
	return id
}

func (f *Fixtures) CreateImportJob(idBankSoal int64, idPengguna user.ID, status importsoal.JobStatus, filePath string) CreatedImportJob {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO import_soal_job (id_bank_soal, id_pengguna, status, file_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id_job
	`, idBankSoal, idPengguna, status, filePath).Scan(&id)
	if err != nil {
		f.t.Fatalf("create import soal job: %v", err)
	}

	f.registerCleanup(`DELETE FROM import_soal_job WHERE id_job = $1`, id)
	return CreatedImportJob{ID: id, IDBankSoal: idBankSoal, IDPengguna: int64(idPengguna), Status: status}
}

func (f *Fixtures) CreateSession(idPengguna user.ID, role user.Role, expiresAt time.Time, revokedAt *time.Time) CreatedSession {
	f.t.Helper()

	var id string
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO sessions (id_pengguna, role, expires_at, revoked_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, idPengguna, role, expiresAt, revokedAt).Scan(&id)
	if err != nil {
		f.t.Fatalf("create session: %v", err)
	}

	f.registerCleanup(`DELETE FROM sessions WHERE id = $1`, id)
	return CreatedSession{ID: id, IDPengguna: int64(idPengguna), Role: role}
}

func (f *Fixtures) CreateGradingJob(idAttempt int64, status ujian.JobStatus) CreatedGradingJob {
	f.t.Helper()

	var id int64
	err := f.q.QueryRow(f.ctx, `
		INSERT INTO grading_jobs (id_attempt, status)
		VALUES ($1, $2)
		RETURNING id_grading_jobs
	`, idAttempt, status).Scan(&id)
	if err != nil {
		f.t.Fatalf("create grading job: %v", err)
	}

	f.registerCleanup(`DELETE FROM grading_jobs WHERE id_grading_jobs = $1`, id)
	return CreatedGradingJob{ID: id, IDAttempt: idAttempt, Status: status}
}

func (f *Fixtures) CreateExamFixture() ExamFixture {
	f.t.Helper()

	kelas := f.CreateKelas(90 + int(uniqueCounter.Add(1)%10))
	namaKelas := f.CreateNamaKelas(kelas.ID, "Kelas "+UniqueSuffix("it"))
	guru := f.CreateUser(user.GURU)
	f.CreateGuruProfile(guru.ID)
	pengawas := f.CreateUser(user.GURU)
	f.CreateGuruProfile(pengawas.ID)
	siswa := f.CreateUser(user.SISWA)
	f.CreateSiswaProfile(siswa.ID, namaKelas.ID)
	mapel := f.CreateMapel(kelas.ID)
	bank := f.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	version := f.CreateBankSoalVersion(bank.ID, guru.ID, "published", 1)
	f.SetActiveBankSoalVersion(bank.ID, version.ID)
	soalPilgan := f.CreateSoal(version.ID, "pilihan_ganda", "Pertanyaan PG "+UniqueSuffix("pg"), 10, 1)
	opsiBenar := f.CreateOpsi(soalPilgan.ID, "Pilihan benar", true)
	opsiSalah := f.CreateOpsi(soalPilgan.ID, "Pilihan salah", false)
	soalEssay := f.CreateSoal(version.ID, "essay", "Pertanyaan essay "+UniqueSuffix("essay"), 15, 2)
	sesi := f.CreateSesi()
	ruang := f.CreateRuangUjian()
	now := time.Date(2099, time.January, 1, 8, 0, 0, 0, time.UTC)
	uj := f.CreateUjian(bank.ID, kelas.ID, Ptr(namaKelas.ID), guru.ID, "Ujian "+UniqueSuffix("it"))
	jadwal := f.CreateJadwalUjian(
		uj.ID,
		sesi.ID,
		ruang.ID,
		pengawas.ID,
		now,
		now.Add(30*time.Minute),
		now.Add(90*time.Minute),
		ujian.BELUM_MULAI,
	)
	peserta := f.CreatePeserta(jadwal.ID, siswa.ID)

	return ExamFixture{
		Kelas:      kelas,
		NamaKelas:  namaKelas,
		Guru:       guru,
		Pengawas:   pengawas,
		Siswa:      siswa,
		Mapel:      mapel,
		BankSoal:   bank,
		Version:    version,
		SoalPilgan: soalPilgan,
		SoalEssay:  soalEssay,
		OpsiBenar:  opsiBenar,
		OpsiSalah:  opsiSalah,
		Sesi:       sesi,
		Ruang:      ruang,
		Ujian:      uj,
		Jadwal:     jadwal,
		Peserta:    peserta,
	}
}

func (f *Fixtures) CountRows(table string, predicate string, args ...any) int {
	f.t.Helper()

	query := "SELECT COUNT(*) FROM " + table
	if strings.TrimSpace(predicate) != "" {
		query += " WHERE " + predicate
	}

	var count int
	if err := f.q.QueryRow(f.ctx, query, args...).Scan(&count); err != nil {
		f.t.Fatalf("count rows %s: %v", table, err)
	}
	return count
}

func (f *Fixtures) registerCleanup(query string, args ...any) {
	if f.addCleanup == nil {
		return
	}

	f.addCleanup(func() {
		if _, err := f.q.Exec(context.Background(), query, args...); err != nil {
			f.t.Errorf("cleanup failed for %q: %v", query, err)
		}
	})
}

func ToBankSoalID(value int64) banksoal.ID {
	return banksoal.ID(value)
}

func ToUjianID(value int64) ujian.ID {
	return ujian.ID(value)
}
