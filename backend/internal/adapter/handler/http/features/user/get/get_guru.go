package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
)

type GetGuruHandler struct {
	svc *user_service.GetGuruService
}

func NewGetGuruHandler(svc *user_service.GetGuruService) *GetGuruHandler {
	return &GetGuruHandler{svc: svc}
}

func (h *GetGuruHandler) ListGuru(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodGet {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListGuruFilters(req)
	if err != nil {
		logger.Info(req.Context(), "invalid guru filters", "op", "user.list_guru", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.ListGuru(req.Context(), filters)
	if err != nil {
		logger.Error(req.Context(), "failed listing guru", "op", "user.list_guru", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list guru")
		}
		return
	}

	responseData := make([]GuruResponseItem,0,len(items))
	for _, item := range items{
		responseData = append(responseData, GuruResponseItem{
			IdPengguna: item.IdPengguna,
			Username: item.Username,
			NoHp: item.NoHp,
			Email: item.Email,
			NamaLengkap: item.NamaLengkap,
			StatusAkun: item.StatusAkun,
			Nip: item.Nip,
			Jabatan: item.Jabatan,
			BidangStudi: item.BidangStudi,
			Foto: item.Foto,
			Role: item.Role,
			JenisKelamin: item.JenisKelamin,
		})
	}
	

	httpResponse.WriteOK(write, http.StatusOK, responseData, "Success")
}

func parseListGuruFilters(req *http.Request) (query.ListGuruFilter, error) {
	values := req.URL.Query()
	filters := query.ListGuruFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(filters.Search, "search"); err != nil {
		return query.ListGuruFilter{}, err
	}

	if statusRaw := strings.TrimSpace(values.Get("status")); statusRaw != "" {
		if err := validator.ValidateInputSafe(statusRaw, "status"); err != nil {
			return query.ListGuruFilter{}, err
		}
		status := user.StatusAkun(strings.ToUpper(statusRaw))
		filters.Status = &status
	}

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListGuruFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListGuruFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	filters.SortBy = strings.TrimSpace(values.Get("sort_by"))
	if err := validator.ValidateInputSafe(filters.SortBy, "sort_by"); err != nil {
		return query.ListGuruFilter{}, err
	}

	if sortDescRaw := strings.TrimSpace(values.Get("sort_desc")); sortDescRaw != "" {
		sortDesc, err := strconv.ParseBool(sortDescRaw)
		if err != nil {
			return query.ListGuruFilter{}, errors.New("sort_desc must be a boolean")
		}
		filters.SortDesc = sortDesc
	}

	if bidangRaw := strings.TrimSpace(values.Get("bidang")); bidangRaw != "" {
		if err := validator.ValidateInputSafe(bidangRaw, "bidang"); err != nil {
			return query.ListGuruFilter{}, err
		}
		filters.Bidang = &bidangRaw
	} else if bidangStudiRaw := strings.TrimSpace(values.Get("bidang_studi")); bidangStudiRaw != "" {
		if err := validator.ValidateInputSafe(bidangStudiRaw, "bidang_studi"); err != nil {
			return query.ListGuruFilter{}, err
		}
		filters.Bidang = &bidangStudiRaw
	}

	return filters, nil
}
