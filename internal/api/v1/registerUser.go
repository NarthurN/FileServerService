package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

func (a *api) RegisterUser(ctx context.Context, req *fileserverV1.RegisterRequest) (fileserverV1.RegisterUserRes, error) {
	log.Printf("🔄 API: Регистрация пользователя %s", req.Login)

	user, err := a.service.RegisterUser(ctx, req.Token, req.Login, req.Pswd)
	if err != nil {
		log.Printf("🚨 API: Ошибка регистрации пользователя %s: %v", req.Login, err)
		if errors.Is(err, model.ErrInvalidAdminToken) {
			return &fileserverV1.BadRequestError{
				Error: fileserverV1.BadRequestErrorError{
					Code: 400,
					Text: fmt.Sprintf("🚨 Токен %s: %v", req.Token, err),
				},
			}, nil
		}
		if errors.Is(err, model.ErrLoginAlreadyExists) {
			return &fileserverV1.BadRequestError{
				Error: fileserverV1.BadRequestErrorError{
					Code: 400,
					Text: fmt.Sprintf("🚨 Логин %s: %v", req.Login, err),
				},
			}, nil
		}

		return &fileserverV1.InternalServerError{
			Error: fileserverV1.InternalServerErrorError{
				Code: 500,
				Text: fmt.Sprintf("🚨 Ошибка регистрации пользователя %s: %v", req.Login, err),
			},
		}, nil
	}
	log.Printf("🎉 API: Пользователь %s успешно зарегистрирован", req.Login)
	return &fileserverV1.RegisterResponse{
		Response: fileserverV1.RegisterResponseResponse{
			Login: user.Login,
		},
	}, nil
}
