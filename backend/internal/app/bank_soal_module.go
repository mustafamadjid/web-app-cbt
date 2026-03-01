package app

import (
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/create"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/update"
	bank_soal_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
	bank_soal_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete"
	bank_soal_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/get"
	bank_soal_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update"
)

type BankSoalModule struct {
	GetService    *bank_soal_get_service.GetBankSoalService
	CreateService *bank_soal_create_service.CreateBankSoalService
	UpdateService *bank_soal_update_service.UpdateBankSoalService
	DeleteService *bank_soal_delete_service.DeleteBankSoalService

	GetHandler    *httpget.GetBankSoalHandler
	CreateHandler *httpcreate.CreateBankSoalHandler
	UpdateHandler *httpupdate.UpdateBankSoalHandler
	DeleteHandler *httpdelete.DeleteBankSoalHandler
}

func BuildBankSoalModule(infra *InfraModule) *BankSoalModule {
	getSvc := bank_soal_get_service.NewGetBankSoalService(infra.bankSoalRepo)
	createSvc := bank_soal_create_service.NewCreateBankSoalService(infra.bankSoalRepo)
	updateSvc := bank_soal_update_service.NewUpdateBankSoalService(infra.bankSoalRepo)
	deleteSvc := bank_soal_delete_service.NewDeleteBankSoalService(infra.bankSoalRepo)

	getHandler := httpget.NewGetBankSoalHandler(getSvc)
	createHandler := httpcreate.NewCreateBankSoalHandler(createSvc)
	updateHandler := httpupdate.NewUpdateBankSoalHandler(updateSvc)
	deleteHandler := httpdelete.NewDeleteBankSoalHandler(deleteSvc)

	return &BankSoalModule{
		GetService:    getSvc,
		CreateService: createSvc,
		UpdateService: updateSvc,
		DeleteService: deleteSvc,
		GetHandler:    getHandler,
		CreateHandler: createHandler,
		UpdateHandler: updateHandler,
		DeleteHandler: deleteHandler,
	}
}
