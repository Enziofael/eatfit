using System;
using System.Net.Http;
using System.Threading.Tasks;
using Eatfit.V1;
using EatfitDesktop.Models;
using EatfitDesktop.Services.Interfaces;
using Grpc.Core;
using Grpc.Net.Client;

namespace EatfitDesktop.Services
{
    public class GrpcAuthService : IAuthService, IDisposable
    {
        private readonly GrpcChannel _channel;
        private readonly Eatfit.V1.AuthService.AuthServiceClient _client;
        private bool _disposed;

        public GrpcAuthService(string serverUrl = "http://localhost:50051")
        {
            var handler = new HttpClientHandler
            {
                ServerCertificateCustomValidationCallback =
                    HttpClientHandler.DangerousAcceptAnyServerCertificateValidator
            };

            _channel = GrpcChannel.ForAddress(serverUrl, new GrpcChannelOptions
            {
                HttpHandler = handler
            });

            _client = new Eatfit.V1.AuthService.AuthServiceClient(_channel);
        }

        public async Task<AuthResult> RegisterAsync(RegisterModel model)
        {
            try
            {
                var request = new RegisterRequest
                {
                    Email = model.Email,
                    Login = model.Login,
                    Password = model.Password,
                    PasswordConfirmation = model.PasswordConfirmation
                };

                var response = await _client.RegisterAsync(request);

                return new AuthResult
                {
                    Success = response.Success,
                    Message = response.Message,
                    UserId = response.UserId
                };
            }
            catch (RpcException ex)
            {
                return new AuthResult { Success = false, Message = ex.Status.Detail };
            }
        }

        public async Task<AuthResult> VerifyEmailAsync(string userId, string code)
        {
            try
            {
                var request = new VerifyEmailRequest
                {
                    UserId = userId,
                    VerificationCode = code
                };

                var response = await _client.VerifyEmailAsync(request);

                return new AuthResult
                {
                    Success = response.Success,
                    Message = response.Message
                };
            }
            catch (RpcException ex)
            {
                return new AuthResult { Success = false, Message = ex.Status.Detail };
            }
        }

        public async Task<AuthResult> LoginAsync(LoginModel model)
        {
            try
            {
                var request = new LoginRequest
                {
                    Password = model.Password,
                    DeviceInfo = $"WPF Desktop - {Environment.MachineName}"
                };

                if (model.LoginIdentifier.Contains('@'))
                    request.Email = model.LoginIdentifier;
                else
                    request.Login = model.LoginIdentifier;

                var response = await _client.LoginAsync(request);

                return new AuthResult
                {
                    Success = true,
                    Message = "Login successful",
                    AccessToken = response.AccessToken,
                    RefreshToken = response.RefreshToken,
                    UserId = response.User.UserId,
                    Email = response.User.Email,
                    Login = response.User.Login
                };
            }
            catch (RpcException ex)
            {
                return new AuthResult { Success = false, Message = ex.Status.Detail };
            }
        }

        public async Task<AuthResult> ForgotPasswordAsync(string loginIdentifier)
        {
            try
            {
                var request = new ForgotPasswordRequest
                {
                    LoginIdentifier = loginIdentifier
                };

                var response = await _client.ForgotPasswordAsync(request);

                return new AuthResult
                {
                    Success = response.Success,
                    Message = response.Message,
                    UserId = response.ResetToken // Используем как временное хранилище токена
                };
            }
            catch (RpcException ex)
            {
                return new AuthResult { Success = false, Message = ex.Status.Detail };
            }
        }

        public async Task<AuthResult> ResetPasswordAsync(string resetToken, string code, string newPassword)
        {
            try
            {
                var request = new ResetPasswordRequest
                {
                    ResetToken = resetToken,
                    VerificationCode = code,
                    NewPassword = newPassword,
                    PasswordConfirmation = newPassword // Если подтверждение не нужно отдельно
                };

                var response = await _client.ResetPasswordAsync(request);

                return new AuthResult
                {
                    Success = response.Success,
                    Message = response.Message
                };
            }
            catch (RpcException ex)
            {
                return new AuthResult { Success = false, Message = ex.Status.Detail };
            }
        }

        public async Task<AuthResult> DeleteAccountAsync(string userId)
        {
            await Task.CompletedTask;
            return new AuthResult { Success = true };
        }

        public async Task<AuthResult> RefreshTokenAsync(string refreshToken)
        {
            try
            {
                var request = new RefreshTokenRequest { RefreshToken = refreshToken };
                var response = await _client.RefreshTokenAsync(request);

                return new AuthResult
                {
                    Success = true,
                    AccessToken = response.AccessToken,
                    RefreshToken = response.RefreshToken
                };
            }
            catch
            {
                return new AuthResult { Success = false, Message = "Session expired" };
            }
        }

        public async Task LogoutAsync(string refreshToken)
        {
            try
            {
                await _client.LogoutAsync(new LogoutRequest { RefreshToken = refreshToken });
            }
            catch
            {
                // Игнорируем ошибки при выходе
            }
        }

        public async Task<AuthResult> ChangeLoginAsync(string userId, string newLogin, string password)
        {
            try
            {
                var request = new ChangeLoginRequest
                {
                    UserId = userId,
                    NewLogin = newLogin,
                    Password = password
                };
                var response = await _client.ChangeLoginAsync(request);
                return new AuthResult { Success = response.Success, Message = response.Message };
            }
            catch (RpcException ex)
            {
                return new AuthResult { Success = false, Message = ex.Status.Detail };
            }
        }

        public async Task<AuthResult> ChangePasswordAsync(string userId, string currentPassword,
            string newPassword, string confirmPassword)
        {
            try
            {
                var request = new ChangePasswordRequest
                {
                    UserId = userId,
                    CurrentPassword = currentPassword,
                    NewPassword = newPassword,
                    PasswordConfirmation = confirmPassword
                };
                var response = await _client.ChangePasswordAsync(request);
                return new AuthResult { Success = response.Success, Message = response.Message };
            }
            catch (RpcException ex)
            {
                return new AuthResult { Success = false, Message = ex.Status.Detail };
            }
        }

        public void Dispose()
        {
            if (_disposed) return;
            _disposed = true;
            _channel?.Dispose();
            GC.SuppressFinalize(this);
        }
    }
}