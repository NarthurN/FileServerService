package docs

import (
	"github.com/NarthurN/FileServerService/internal/service"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

var _ fileserverV1.Handler = (*api)(nil)

type api struct {
	fileserverV1.UnimplementedHandler

	service service.FileServerService
}

func NewAPI(service service.FileServerService) *api {
	return &api{
		service: service,
	}
}
