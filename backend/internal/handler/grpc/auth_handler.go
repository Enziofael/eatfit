// internal/handler/grpc/auth_handler.go
package grpc

import (
	"context"
	"log"

	pb "github.com/Enziofael/eatfit/backend/api/gen/eatfit/v1"
	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/Enziofael/eatfit/backend/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthHandler - обработчик gRPC запросов аутентификации
type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

// NewAuthHandler создаёт новый экземпляр AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register обрабатывает запрос на регистрацию
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Конвертируем proto запрос в доменную модель
	registerReq := &model.RegisterAccountRequest{
		Email:                req.GetEmail(),
		Login:                req.GetLogin(),
		Password:             req.GetPassword(),
		PasswordConfirmation: req.GetPasswordConfirmation(),
	}

	// Вызываем сервис
	account, _, err := h.authService.Register(ctx, registerReq)
	if err != nil {
		log.Printf("ERROR Register: %v", err) // ← ДОБАВЛЕНО
		return nil, handleAuthError(err)
	}

	// Формируем ответ
	return &pb.RegisterResponse{
		UserId:  account.ID.String(),
		Message: "Registration successful. Please check your email for verification code.",
		Success: true,
	}, nil
}

// VerifyEmail обрабатывает запрос на подтверждение email
func (h *AuthHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	accountID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
	}

	err = h.authService.VerifyEmail(ctx, accountID, req.GetVerificationCode())
	if err != nil {
		log.Printf("ERROR VerifyEmail: %v", err) // ← ДОБАВЛЕНО
		return nil, handleAuthError(err)
	}

	return &pb.VerifyEmailResponse{
		Success: true,
		Message: "Email verified successfully. You can now login.",
	}, nil
}

// Login обрабатывает запрос на вход в систему
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Определяем, что ввёл пользователь: email или login
	var identifier string
	switch method := req.LoginMethod.(type) {
	case *pb.LoginRequest_Email:
		identifier = method.Email
	case *pb.LoginRequest_Login:
		identifier = method.Login
	default:
		return nil, status.Errorf(codes.InvalidArgument, "email or login is required")
	}

	loginReq := &model.LoginRequest{
		LoginIdentifier: identifier,
		Password:        req.GetPassword(),
		DeviceInfo:      req.GetDeviceInfo(),
	}

	result, err := h.authService.Login(ctx, loginReq)
	if err != nil {
		log.Printf("ERROR Login: %v", err) // ← ДОБАВЛЕНО
		return nil, handleAuthError(err)
	}

	return &pb.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    result.TokenType,
		User: &pb.UserProfile{
			UserId:        result.Account.ID.String(),
			Email:         result.Account.Email,
			Login:         result.Account.Login,
			EmailVerified: result.Account.EmailVerified,
			CreatedAt:     timestamppb.New(result.Account.CreatedAt),
		},
	}, nil
}

// RefreshToken обрабатывает запрос на обновление токенов
func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	result, err := h.authService.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		log.Printf("ERROR RefreshToken: %v", err) // ← ДОБАВЛЕНО
		return nil, handleAuthError(err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// Logout обрабатывает запрос на выход из системы
func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.authService.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		log.Printf("ERROR Logout: %v", err) // ← ДОБАВЛЕНО
		return nil, handleAuthError(err)
	}

	return &pb.LogoutResponse{
		Success: true,
		Message: "Logged out successfully.",
	}, nil
}

// ResendVerification обрабатывает запрос на повторную отправку кода
func (h *AuthHandler) ResendVerification(ctx context.Context, req *pb.ResendVerificationRequest) (*pb.ResendVerificationResponse, error) {
	accountID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
	}

	err = h.authService.ResendVerification(ctx, accountID)
	if err != nil {
		log.Printf("ERROR ResendVerification: %v", err) // ← ДОБАВЛЕНО
		return nil, handleAuthError(err)
	}

	return &pb.ResendVerificationResponse{
		Success: true,
		Message: "Verification code has been resent to your email.",
	}, nil
}

// handleAuthError конвертирует бизнес-ошибки в gRPC статус коды
func handleAuthError(err error) error {
	if err == nil {
		return nil
	}

	// Логируем ошибку
	log.Printf("Auth error: %v", err)

	errMsg := err.Error()

	switch {
	case contains(errMsg, "email already exists"),
		contains(errMsg, "login already exists"):
		return status.Errorf(codes.AlreadyExists, "%s", errMsg)

	case contains(errMsg, "invalid email"),
		contains(errMsg, "invalid login"),
		contains(errMsg, "invalid password"),
		contains(errMsg, "password and confirmation"),
		contains(errMsg, "invalid user ID"),
		contains(errMsg, "invalid verification code"):
		return status.Errorf(codes.InvalidArgument, "%s", errMsg)

	case contains(errMsg, "email not verified"):
		return status.Errorf(codes.FailedPrecondition, "%s", errMsg)

	case contains(errMsg, "not found"),
		contains(errMsg, "account not found"):
		return status.Errorf(codes.NotFound, "%s", errMsg)

	case contains(errMsg, "expired"):
		return status.Errorf(codes.DeadlineExceeded, "%s", errMsg)

	case contains(errMsg, "too many"):
		return status.Errorf(codes.ResourceExhausted, "%s", errMsg)

	default:
		// Логируем неизвестную ошибку
		log.Printf("Unknown error: %v", err)
		return status.Errorf(codes.Internal, "internal server error: %v", err)
	}
}

// contains проверяет, содержит ли строка подстроку
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ForgotPassword обрабатывает запрос на восстановление пароля
func (h *AuthHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	resetToken, err := h.authService.ForgotPassword(ctx, req.GetLoginIdentifier())
	if err != nil {
		return nil, handleAuthError(err)
	}

	return &pb.ForgotPasswordResponse{
		Success:    true,
		Message:    "If the account exists, a reset code has been sent to your email.",
		ResetToken: resetToken,
	}, nil
}

// ResetPassword обрабатывает запрос на сброс пароля
func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	err := h.authService.ResetPassword(
		ctx,
		req.GetResetToken(),
		req.GetVerificationCode(),
		req.GetNewPassword(),
		req.GetPasswordConfirmation(),
	)
	if err != nil {
		return nil, handleAuthError(err)
	}

	return &pb.ResetPasswordResponse{
		Success: true,
		Message: "Password has been reset successfully. You can now login with your new password.",
	}, nil
}

func (h *AuthHandler) ChangeLogin(ctx context.Context, req *pb.ChangeLoginRequest) (*pb.ChangeLoginResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	err = h.authService.ChangeLogin(ctx, userID, req.GetNewLogin(), req.GetPassword())
	if err != nil {
		return nil, handleAuthError(err)
	}

	return &pb.ChangeLoginResponse{
		Success:  true,
		Message:  "Login changed successfully",
		NewLogin: req.GetNewLogin(),
	}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	err = h.authService.ChangePassword(ctx, userID, req.GetCurrentPassword(), req.GetNewPassword(), req.GetPasswordConfirmation())
	if err != nil {
		return nil, handleAuthError(err)
	}

	return &pb.ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}
