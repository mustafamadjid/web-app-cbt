package app

import (
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/update"
	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"

	getsvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
	updatesvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/update"
)

type ProfilSekolahModule struct {
	GetHandler    *httpget.GetProfilSekolahHandler
	UpdateHandler *httpupdate.UpdateProfilSekolahHandler
}

func BuildProfilSekolahModule(cfg Config, infra *InfraModule) *ProfilSekolahModule {
	store := httpx.ImageStore{
		Dir:      cfg.ImageStore.Dir,
		BaseURL:  cfg.ImageStore.BaseURL,
		Route:    cfg.ImageStore.Route,
		MaxBytes: cfg.ImageStore.MaxBytes,
	}

	getSvc := getsvc.NewGetProfilSekolahService(infra.profilSekolah)
	updateSvc := updatesvc.NewUpdateProfilSekolahService(infra.profilSekolah)

	return &ProfilSekolahModule{
		GetHandler:    httpget.NewGetProfilSekolahHandler(getSvc),
		UpdateHandler: httpupdate.NewUpdateProfilSekolahHandler(updateSvc, store),
	}
}
