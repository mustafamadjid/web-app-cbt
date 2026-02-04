package httpx

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
)

type AktivitasUserHandler struct {
	svc *aktivitas_user_service.AktivitasUserService
}

type AktivitasUserResponse struct {
	IdAktivitas aktivitas_user.AktivitasID `json:"id_aktivitas"`
	IdPengguna  user.ID                    `json:"id_pengguna"`
	Action      aktivitas_user.Action      `json:"action"`
	Description string                     `json:"description"`
	IpAddress   string                     `json:"ip_address"`
	CreatedAt   string                     `json:"created_at"`
}

func NewAktivitasUserHandler(svc *aktivitas_user_service.AktivitasUserService) *AktivitasUserHandler {
	return &AktivitasUserHandler{svc: svc}
}

func (h *AktivitasUserHandler) GetAktivitasUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	data, err := h.svc.GetAktivitasUserService(r.Context())
	if err != nil {
		logger.Error(r.Context(), "failed getting aktivitas user", "op", "aktivitas_user.get", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get aktivitas user")
		return
	}

	response := make([]AktivitasUserResponse, 0, len(data))
	for _, item := range data {
		response = append(response, AktivitasUserResponse{
			IdAktivitas: item.IdAktivitas,
			IdPengguna:  item.IdPengguna,
			Action:      item.Action,
			Description: item.Description,
			IpAddress:   item.IpAddress,
			CreatedAt:   formatAktivitasTimestamp(item.CreatedAt),
		})
	}

	httpResponse.WriteOK(w, http.StatusOK, response, "success")
}

func formatAktivitasTimestamp(t time.Time) string {
	days := map[time.Weekday]string{
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
		time.Sunday:    "Minggu",
	}

	months := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}

	dayName := days[t.Weekday()]
	monthName := months[t.Month()]
	if dayName == "" || monthName == "" {
		return t.Format("2006-01-02 15:04")
	}

	return dayName + " " + t.Format("02") + " " + monthName + " " + t.Format("15.04")
}
