package user_test

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	userout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	"github.com/stretchr/testify/mock"
)

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) Begin(ctx context.Context) (txout.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(txout.Tx), args.Error(1)
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Pengguna() userout.UserRepository {
	return m.Called().Get(0).(userout.UserRepository)
}

func (m *MockTx) ProfilGuru() userout.ProfilGuruRepository {
	return m.Called().Get(0).(userout.ProfilGuruRepository)
}

func (m *MockTx) ProfilSiswa() userout.ProfilSiswaRepository {
	return m.Called().Get(0).(userout.ProfilSiswaRepository)
}

func (m *MockTx) Commit() error {
	return m.Called().Error(0)
}

func (m *MockTx) Rollback() error {
	return m.Called().Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.Pengguna), args.Error(1)
}

func (m *MockUserRepository) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error) {
	args := m.Called(ctx, pengguna)
	return args.Get(0).(user.ID), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna updatepatch.Pengguna) error {
	return m.Called(ctx, idPengguna, pengguna).Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id user.ID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepository) DeleteUsers(ctx context.Context, ids []user.ID) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.Pengguna), args.Error(1)
}

func (m *MockUserRepository) UserExistByEmail(ctx context.Context, email user.Email) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UserExistByNoHp(ctx context.Context, noHp string) (bool, error) {
	args := m.Called(ctx, noHp)
	return args.Bool(0), args.Error(1)
}

type MockProfilGuruRepository struct {
	mock.Mock
}

func (m *MockProfilGuruRepository) CreateProfilGuru(ctx context.Context, p user.ProfilGuru) (user.ID, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(user.ID), args.Error(1)
}

func (m *MockProfilGuruRepository) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	args := m.Called(ctx, nip)
	return args.Bool(0), args.Error(1)
}

func (m *MockProfilGuruRepository) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.DataGuru), args.Error(1)
}

func (m *MockProfilGuruRepository) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, p updatepatch.ProfilGuru) error {
	return m.Called(ctx, idPengguna, p).Error(0)
}

type MockProfilSiswaRepository struct {
	mock.Mock
}

func (m *MockProfilSiswaRepository) CreateProfilSiswa(ctx context.Context, p user.ProfilSiswa) (user.ID, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(user.ID), args.Error(1)
}

func (m *MockProfilSiswaRepository) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	args := m.Called(ctx, nisn)
	return args.Bool(0), args.Error(1)
}

func (m *MockProfilSiswaRepository) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.DataSiswa), args.Error(1)
}

func (m *MockProfilSiswaRepository) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, p updatepatch.ProfilSiswa) error {
	return m.Called(ctx, idPengguna, p).Error(0)
}

type MockGetSiswaListRepo struct {
	mock.Mock
}

func (m *MockGetSiswaListRepo) GetListSiswa(ctx context.Context, f query.ListSiswaFilter) ([]query.SiswaListItem, error) {
	args := m.Called(ctx, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]query.SiswaListItem), args.Error(1)
}

type MockUserResetPasswordRepo struct {
	mock.Mock
}

func (m *MockUserResetPasswordRepo) ResetPassword(ctx context.Context, id user.ID, pass string) error {
	return m.Called(ctx, id, pass).Error(0)
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) GenerateHash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) ComparePaswordAndHashed(hash string, plain string) bool {
	args := m.Called(hash, plain)
	return args.Bool(0)
}

type MockAktivitasUserRepository struct {
	mock.Mock
}

func (m *MockAktivitasUserRepository) CreateAktivitasUser(ctx context.Context, a aktivitas_user.AktivitasUser) error {
	return m.Called(ctx, a).Error(0)
}

func (m *MockAktivitasUserRepository) GetAktivitasUser(ctx context.Context) ([]aktivitas_user.AktivitasUser, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]aktivitas_user.AktivitasUser), args.Error(1)
}

type MockGetGuruListRepo struct {
	mock.Mock
}

func (m *MockGetGuruListRepo) GetListGuru(ctx context.Context, q query.ListGuruFilter) ([]query.GuruListItem, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]query.GuruListItem), args.Error(1)
}

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) GetSession(ctx context.Context, sessionID string) (session.Session, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(session.Session), args.Error(1)
}

func (m *MockSessionRepository) GetSessionByUserId(ctx context.Context, userId user.ID) (session.Session, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(session.Session), args.Error(1)
}

func (m *MockSessionRepository) GetAllActiveSession(ctx context.Context) ([]session.SessionWithUser, error) {
	args := m.Called(ctx)
	return args.Get(0).([]session.SessionWithUser), args.Error(1)
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, userID user.ID, role user.Role, expiresAt time.Time) (string, error) {
	args := m.Called(ctx, userID, role, expiresAt)
	return args.String(0), args.Error(1)
}

func (m *MockSessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	return m.Called(ctx, sessionID).Error(0)
}

func (m *MockSessionRepository) RevokeSessionAllbyUser(ctx context.Context, userID user.ID) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockSessionRepository) RevokeExpiredSessions(ctx context.Context, userID user.ID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) HasActiveSession(ctx context.Context, userID user.ID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

type MockDeleteFileRepo struct {
	mock.Mock
}

func (m *MockDeleteFileRepo) DeleteFile(ctx context.Context, path string) error {
	return m.Called(ctx, path).Error(0)
}
