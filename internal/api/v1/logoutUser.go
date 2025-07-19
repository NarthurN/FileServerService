package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

func (a *api) LogoutUser(ctx context.Context, params fileserverV1.LogoutUserParams) (fileserverV1.LogoutUserRes, error) {
	log.Printf("🔄 API: Выход пользователя %s", params.Token)

	err := a.service.LogoutUser(ctx, params.Token)
	if err != nil {
		log.Printf("🚨 API: Ошибка выхода пользователя %s: %v", params.Token, err)
		if errors.Is(err, model.ErrInvalidToken) {
			return &fileserverV1.UnauthorizedError{
				Error: fileserverV1.UnauthorizedErrorError{
					Code: 401,
					Text: fmt.Sprintf("🚨 Токен %s: %v", params.Token, err),
				},
			}, nil
		}

		return &fileserverV1.InternalServerError{
			Error: fileserverV1.InternalServerErrorError{
				Code: 500,
				Text: fmt.Sprintf("🚨 Токен %s: %v", params.Token, err),
			},
		}, nil
	}

	log.Printf("🎉 API: Пользователь %s успешно вышел", params.Token)
	return &fileserverV1.LogoutResponse{
		Response: fileserverV1.LogoutResponseResponse{
			"token": true,
		},
	}, nil
}
