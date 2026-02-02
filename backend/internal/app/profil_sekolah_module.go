package app

import (
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/update"

	getsvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
	updatesvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/update"
)

type ProfilSekolahModule struct {
	GetHandler    *httpget.GetProfilSekolahHandler
	UpdateHandler *httpupdate.UpdateProfilSekolahHandler
}

func BuildProfilSekolahModule(infra *InfraModule) *ProfilSekolahModule {
	getSvc := getsvc.NewGetProfilSekolahService(infra.profilSekolah)
	updateSvc := updatesvc.NewUpdateProfilSekolahService(infra.profilSekolah)

	return &ProfilSekolahModule{
		GetHandler:    httpget.NewGetProfilSekolahHandler(getSvc),
		UpdateHandler: httpupdate.NewUpdateProfilSekolahHandler(updateSvc),
	}
}
