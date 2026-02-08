package app

import kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"

type KelasModule struct {
	Service *kelas_service.GetKelasService
}

func BuildKelasModule(infra *InfraModule) *KelasModule