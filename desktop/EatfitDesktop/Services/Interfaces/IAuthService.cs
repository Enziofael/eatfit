using System.Threading.Tasks;
using EatfitDesktop.Models;

namespace EatfitDesktop.Services.Interfaces
{
    public interface IAuthService
    {
        Task<AuthResult> RegisterAsync(RegisterModel model);
        Task<AuthResult> VerifyEmailAsync(string userId, string code);
        Task<AuthResult> LoginAsync(LoginModel model);
        Task<AuthResult> ForgotPasswordAsync(string loginIdentifier);
        Task<AuthResult> ResetPasswordAsync(string userId, string code, string newPassword);
        Task<AuthResult> DeleteAccountAsync(string userId);
        Task<AuthResult> RefreshTokenAsync(string refreshToken);
        Task LogoutAsync(string refreshToken);
    }
}