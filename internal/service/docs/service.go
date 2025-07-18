package docs

import (
	"github.com/NarthurN/FileServerService/internal/repository"
	def "github.com/NarthurN/FileServerService/internal/service"
)

var _ def.FileServerService = (*service)(nil)

type service struct {
	repo repository.FileServerRepository
}

func NewService(repo repository.FileServerRepository) *service {
	return &service{
		repo: repo,
	}
}
