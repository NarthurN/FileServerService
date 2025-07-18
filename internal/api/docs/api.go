package docs

import (
	"github.com/NarthurN/FileServerService/internal/service"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

var _ fileserverV1.Handler = (*api)(nil)

type api struct {
	fileserverV1.UnimplementedHandler
	
	docsService service.FileServerService
}

func NewAPI(docsService service.FileServerService) *api {
	return &api{
		docsService: docsService,
	}
}
