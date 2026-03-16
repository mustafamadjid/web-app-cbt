package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	getjawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/get_jawaban"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGetJawabanOwnershipChecker struct {
	checkFn      func(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error)
	checkCalled  bool
	gotSiswaID   int
	gotAttemptID ujian.ID
}

func (f *fakeGetJawabanOwnershipChecker) CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error) {
	f.checkCalled = true
	f.gotSiswaID = idSiswa
	f.gotAttemptID = idAttempt
	if f.checkFn != nil {
		return f.checkFn(ctx, idSiswa, idAttempt)
	}
	return false, nil
}

type fakeGetJawabanRepo struct {
	getFn        func(ctx context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error)
	getCalled    bool
	gotAttemptID ujian.ID
}

func (f *fakeGetJawabanRepo) GetJawabanUjianByAttemptId(ctx context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error) {
	f.getCalled = true
	f.gotAttemptID = idAttempt
	if f.getFn != nil {
		return f.getFn(ctx, idAttempt)
	}
	return nil, nil
}

func (f *fakeGetJawabanRepo) SaveJawabanUjian(context.Context, ujian.ID, []ujian.JawabanUjian) error {
	return nil
}

var _ ujian_repo.JawabanUjianRepository = (*fakeGetJawabanRepo)(nil)

type getJawabanAPIResp struct {
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestGetJawabanUjianHandler(t *testing.T) {
	t.Parallel()

	waktuJawab := time.Date(2026, time.March, 16, 11, 45, 0, 0, time.UTC)
	repoErr := errors.New("repo error")

	tests := []struct {
		name          string
		method        string
		idAttempt     string
		withActor     bool
		checker       *fakeGetJawabanOwnershipChecker
		repo          *fakeGetJawabanRepo
		wantStatus    int
		wantCode      string
		wantMessage   string
		wantCheck     bool
		wantGet       bool
		assertPayload func(t *testing.T, data GetJawabanUjianResponse)
	}{
		{
			name:        "wrong method",
			method:      http.MethodPost,
			idAttempt:   "7",
			checker:     &fakeGetJawabanOwnershipChecker{},
			repo:        &fakeGetJawabanRepo{},
			wantStatus:  http.StatusMethodNotAllowed,
			wantCode:    "METHOD_NOT_ALLOWED",
			wantMessage: "method not allowed",
		},
		{
			name:        "missing actor",
			method:      http.MethodGet,
			idAttempt:   "7",
			checker:     &fakeGetJawabanOwnershipChecker{},
			repo:        &fakeGetJawabanRepo{},
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error : failed get actor from context",
		},
		{
			name:        "missing id attempt",
			method:      http.MethodGet,
			withActor:   true,
			checker:     &fakeGetJawabanOwnershipChecker{},
			repo:        &fakeGetJawabanRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid id attempt",
		},
		{
			name:        "invalid id attempt",
			method:      http.MethodGet,
			idAttempt:   "abc",
			withActor:   true,
			checker:     &fakeGetJawabanOwnershipChecker{},
			repo:        &fakeGetJawabanRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid id attempt",
		},
		{
			name:        "non positive id attempt",
			method:      http.MethodGet,
			idAttempt:   "0",
			withActor:   true,
			checker:     &fakeGetJawabanOwnershipChecker{},
			repo:        &fakeGetJawabanRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid id attempt",
		},
		{
			name:      "attempt not owned",
			method:    http.MethodGet,
			idAttempt: "7",
			withActor: true,
			checker: &fakeGetJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return false, nil
				},
			},
			repo:        &fakeGetJawabanRepo{},
			wantStatus:  http.StatusNotFound,
			wantCode:    "NOT_FOUND",
			wantMessage: "data not found",
			wantCheck:   true,
		},
		{
			name:      "service error becomes internal error",
			method:    http.MethodGet,
			idAttempt: "7",
			withActor: true,
			checker: &fakeGetJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo: &fakeGetJawabanRepo{
				getFn: func(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
					return nil, repoErr
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error",
			wantCheck:   true,
			wantGet:     true,
		},
		{
			name:      "success with data",
			method:    http.MethodGet,
			idAttempt: "7",
			withActor: true,
			checker: &fakeGetJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo: &fakeGetJawabanRepo{
				getFn: func(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
					essay := "jawaban essay"
					idPilihan := ujian.ID(13)
					return []ujian.JawabanUjian{
						{
							IdJawaban:    5,
							IdSoal:       11,
							IdPilihan:    &idPilihan,
							JawabanEssay: &essay,
							WaktuJawab:   &waktuJawab,
						},
					}, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantCheck:   true,
			wantGet:     true,
			assertPayload: func(t *testing.T, data GetJawabanUjianResponse) {
				t.Helper()
				assert.Equal(t, 7, data.IDAttempt)
				require.Len(t, data.Jawaban, 1)
				assert.Equal(t, 5, data.Jawaban[0].IDJawaban)
				assert.Equal(t, 11, data.Jawaban[0].IDSoal)
				require.NotNil(t, data.Jawaban[0].IDPilihan)
				assert.Equal(t, 13, *data.Jawaban[0].IDPilihan)
				require.NotNil(t, data.Jawaban[0].JawabanEssay)
				assert.Equal(t, "jawaban essay", *data.Jawaban[0].JawabanEssay)
				require.NotNil(t, data.Jawaban[0].WaktuJawab)
				assert.Equal(t, waktuJawab.Format(time.RFC3339), *data.Jawaban[0].WaktuJawab)
			},
		},
		{
			name:      "success with empty data",
			method:    http.MethodGet,
			idAttempt: "7",
			withActor: true,
			checker: &fakeGetJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo: &fakeGetJawabanRepo{
				getFn: func(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
					return []ujian.JawabanUjian{}, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantCheck:   true,
			wantGet:     true,
			assertPayload: func(t *testing.T, data GetJawabanUjianResponse) {
				t.Helper()
				assert.Equal(t, 7, data.IDAttempt)
				require.NotNil(t, data.Jawaban)
				assert.Len(t, data.Jawaban, 0)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			getter := getjawaban_service.NewGetJawabanUjianService(tc.repo)
			handler := NewGetJawabanUjianHandler(getjawaban_service.NewSiswaGetJawabanUjianService(tc.checker, getter))

			req := httptest.NewRequest(tc.method, "/siswa/ujian/jawaban/"+tc.idAttempt, nil)
			if tc.withActor {
				req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 9, Role: user.SISWA}))
			}

			rec := httptest.NewRecorder()
			var params httprouter.Params
			if tc.idAttempt != "" {
				params = httprouter.Params{{Key: "idAttempt", Value: tc.idAttempt}}
			}
			handler.GetJawabanUjian(rec, req, params)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantCheck, tc.checker.checkCalled)
			assert.Equal(t, tc.wantGet, tc.repo.getCalled)
			if tc.wantCheck {
				assert.Equal(t, 9, tc.checker.gotSiswaID)
				assert.Equal(t, ujian.ID(7), tc.checker.gotAttemptID)
			}
			if tc.wantGet {
				assert.Equal(t, ujian.ID(7), tc.repo.gotAttemptID)
			}

			var resp getJawabanAPIResp
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tc.wantCode == "" {
				assert.Equal(t, tc.wantMessage, resp.Message)
				assert.Nil(t, resp.Error)
				if tc.assertPayload != nil {
					tc.assertPayload(t, decodeGetJawabanResponse(t, resp))
				}
				return
			}

			require.NotNil(t, resp.Error)
			assert.Equal(t, tc.wantCode, resp.Error.Code)
			assert.Equal(t, tc.wantMessage, resp.Error.Message)
		})
	}
}

func decodeGetJawabanResponse(t *testing.T, resp getJawabanAPIResp) GetJawabanUjianResponse {
	t.Helper()

	var data GetJawabanUjianResponse
	require.NoError(t, json.Unmarshal(resp.Data, &data))
	return data
}
