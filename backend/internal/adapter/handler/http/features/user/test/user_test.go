package user_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	usercreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/create"
	userdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/delete"
	userget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/get"
	userpassword "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/reset_password"
	userupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/update"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	usercreate_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	userdelete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	userget_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
	userpassword_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/reset_password"
	userupdate_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGuru(t *testing.T) {
	mockTxm := new(MockTxManager)
	mockHasher := new(MockPasswordHasher)
	
	svc := usercreate_service.NewCreateGuruService(mockTxm, mockHasher)
	handler := usercreate.NewCreateUserHandler(svc, httphelper.ImageStore{}, nil)

	t.Run("Success Create Guru", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
		ctx := middleware.WithActor(context.Background(), actor)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("username", "guru_test")
		_ = writer.WriteField("password", "password123")
		_ = writer.WriteField("nama_lengkap", "Guru Test")
		_ = writer.WriteField("jenis_kelamin", "LAKI-LAKI")
		_ = writer.WriteField("nip", "123456789012345678")
		_ = writer.WriteField("jabatan", "Guru")
		_ = writer.WriteField("bidang_studi", "Informatika")
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/guru", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		mockHasher.On("GenerateHash", "password123").Return("hashed_pass", nil).Once()
		
		mockTx := new(MockTx)
		mockTxm.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		
		mockUserRepo := new(MockUserRepository)
		mockTx.On("Pengguna").Return(mockUserRepo)
		mockUserRepo.On("UserExistByUsername", mock.Anything, "guru_test").Return(false, nil).Once()
		mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(user.ID(10), nil).Once()
		
		mockGuruRepo := new(MockProfilGuruRepository)
		mockTx.On("ProfilGuru").Return(mockGuruRepo)
		mockGuruRepo.On("ExistByNIP", mock.Anything, user.NIP("123456789012345678")).Return(false, nil).Once()
		mockGuruRepo.On("CreateProfilGuru", mock.Anything, mock.Anything).Return(user.ID(1), nil).Once()
		
		mockTx.On("Commit").Return(nil).Once()
		mockTx.On("Rollback").Return(nil)

		handler.CreateGuru(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized - No Actor", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("username", "guru_test")
		_ = writer.WriteField("password", "password123")
		_ = writer.WriteField("nama_lengkap", "Guru Test")
		_ = writer.WriteField("jenis_kelamin", "LAKI-LAKI")
		_ = writer.WriteField("nip", "123456789012345678")
		_ = writer.WriteField("jabatan", "Guru")
		_ = writer.WriteField("bidang_studi", "Informatika")
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/guru", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		handler.CreateGuru(w, req, nil)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetGuru(t *testing.T) {
	mockListRepo := new(MockGetGuruListRepo)
	mockProfilRepo := new(MockProfilGuruRepository)
	svc := userget_service.NewGetListGuruService(mockListRepo, mockProfilRepo)
	handler := userget.NewGetGuruHandler(svc)

	t.Run("Success Get List Guru", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/guru", nil)
		w := httptest.NewRecorder()

		mockListRepo.On("GetListGuru", mock.Anything, mock.Anything).Return([]query.GuruListItem{}, nil).Once()

		handler.ListGuru(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success Get Guru By ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/guru/1", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "id", Value: "1"}}

		mockProfilRepo.On("FindProfilGuruByID", mock.Anything, user.ID(1)).Return(user.DataGuru{
			IdPengguna: 1,
			Username:   "testguru",
		}, nil).Once()

		handler.GetGuruByID(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/guru/99", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "id", Value: "99"}}

		mockProfilRepo.On("FindProfilGuruByID", mock.Anything, user.ID(99)).Return(user.DataGuru{}, coreerror.ErrNotFound).Once()

		handler.GetGuruByID(w, req, params)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetSiswa(t *testing.T) {
	mockListRepo := new(MockGetSiswaListRepo)
	mockProfilRepo := new(MockProfilSiswaRepository)
	svc := userget_service.NewGetListSiswaService(mockListRepo, mockProfilRepo)
	handler := userget.NewGetSiswaHandler(svc)

	t.Run("Success Get List Siswa", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/siswa", nil)
		w := httptest.NewRecorder()

		mockListRepo.On("GetListSiswa", mock.Anything, mock.Anything).Return([]query.SiswaListItem{}, nil).Once()

		handler.ListSiswa(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success Get Siswa By ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/siswa/2", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "id", Value: "2"}}

		mockProfilRepo.On("FindProfilSiswaByID", mock.Anything, user.ID(2)).Return(user.DataSiswa{
			IdPengguna: 2,
			Username:   "testsiswa",
		}, nil).Once()

		handler.GetSiswaByID(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUpdateGuru(t *testing.T) {
	mockTxm := new(MockTxManager)
	mockSessions := new(MockSessionRepository)
	mockDeleteFile := new(MockDeleteFileRepo)
	mockUsers := new(MockUserRepository)
	
	svc := userupdate_service.NewUpdateUserService(mockTxm, mockSessions, mockDeleteFile, mockUsers)
	handler := userupdate.NewUpdateUserHandler(svc, httphelper.ImageStore{}, nil)

	t.Run("Success Update Guru", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
		ctx := middleware.WithActor(context.Background(), actor)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("nama_lengkap", "Guru Updated")
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPatch, "/guru/10", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "id", Value: "10"}}

		mockTx := new(MockTx)
		mockTxm.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		
		mockUserRepo := new(MockUserRepository)
		mockTx.On("Pengguna").Return(mockUserRepo)
		mockUserRepo.On("UpdateUser", mock.Anything, user.ID(10), mock.Anything).Return(nil).Once()
		
		mockTx.On("Commit").Return(nil).Once()
		mockTx.On("Rollback").Return(nil)
		
		mockSessions.On("RevokeSessionAllbyUser", mock.Anything, user.ID(10)).Return(nil).Once()

		handler.UpdateGuru(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockDeleteFile := new(MockDeleteFileRepo)
	svc := userdelete_service.NewDeleteUserService(mockUserRepo, mockDeleteFile)
	handler := userdelete.NewDeleteUserHandler(svc, nil)

	t.Run("Success Delete User", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
		ctx := middleware.WithActor(context.Background(), actor)

		req := httptest.NewRequest(http.MethodDelete, "/user/10", nil).WithContext(ctx)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "id", Value: "10"}}

		mockUserRepo.On("FindUserByID", mock.Anything, user.ID(10)).Return(user.Pengguna{ID: 10, Foto: "path/to/foto.jpg"}, nil).Once()
		mockDeleteFile.On("DeleteFile", mock.Anything, "path/to/foto.jpg").Return(nil).Once()
		mockUserRepo.On("DeleteUser", mock.Anything, user.ID(10)).Return(nil).Once()

		handler.DeleteUser(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestResetPassword(t *testing.T) {
	mockResetRepo := new(MockUserResetPasswordRepo)
	mockHasher := new(MockPasswordHasher)
	svc := userpassword_service.NewResetPasswordService(mockResetRepo, mockHasher)
	handler := userpassword.NewResetPasswordHandler(svc)

	t.Run("Success Reset Password", func(t *testing.T) {
		body := strings.NewReader(`{"password":"NewPassword123!"}`)
		req := httptest.NewRequest(http.MethodPut, "/user/10/reset-password", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idPengguna", Value: "10"}}

		mockHasher.On("GenerateHash", "NewPassword123!").Return("hashed_new_pass", nil).Once()
		mockResetRepo.On("ResetPassword", mock.Anything, user.ID(10), "hashed_new_pass").Return(nil).Once()

		handler.ResetPasswordHandler(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		body := strings.NewReader(`{"password":"short"}`)
		req := httptest.NewRequest(http.MethodPut, "/user/10/reset-password", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idPengguna", Value: "10"}}

		handler.ResetPasswordHandler(w, req, params)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "input too short")
	})
}
