package v1

import (
	"context"
	"fmt"
	"log"

	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

func (a *api) LoginUser(ctx context.Context, req *fileserverV1.LoginRequest) (fileserverV1.LoginUserRes, error) {
	log.Printf("🔄 API: Аутентификация пользователя %s", req.Login)

	token, err := a.service.AuthenticateUser(ctx, req.Login, req.Pswd)
	if err != nil {
		log.Printf("🚨 API: Ошибка аутентификации пользователя %s: %v", req.Login, err)
		return &fileserverV1.BadRequestError{
			Error: fileserverV1.BadRequestErrorError{
				Code: 400,
				Text: fmt.Sprintf("🚨 Логин %s: %v", req.Login, err),
			},
		}, nil
	}

	log.Printf("🎉 API: Пользователь %s успешно аутентифицирован", req.Login)
	return &fileserverV1.LoginResponse{
		Response: fileserverV1.LoginResponseResponse{
			Token: token,
		},
	}, nil
}
