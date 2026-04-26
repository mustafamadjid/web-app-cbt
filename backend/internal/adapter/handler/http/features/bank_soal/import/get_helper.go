package importhandler

import (
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

func toJobResponse(j importsoal.ImportSoalJob) jobResponse {
	return jobResponse{
		IDJob:      j.IDJob,
		IDBankSoal: j.IDBankSoal,
		Status:     string(j.Status),
		ErrorMsg:   j.ErrorMsg,
		WarningMsg: j.WarningMsg,
		TotalSoal:  j.TotalSoal,
		CreatedAt:  httphelper.FormatRFC3339(j.CreatedAt),
		UpdatedAt:  httphelper.FormatRFC3339(j.UpdatedAt),
	}
}

func toJobResponses(items []importsoal.ImportSoalJob) []jobResponse {
	response := make([]jobResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toJobResponse(item))
	}

	return response
}
