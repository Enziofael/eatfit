using System;
using System.Threading.Tasks;
using EatfitDesktop.Helpers;
using EatfitDesktop.Services.Interfaces;

namespace EatfitDesktop.ViewModels
{
    public class ForgotPasswordViewModel : BaseViewModel
    {
        private readonly IAuthService _authService;

        private string _loginIdentifier = string.Empty;
        public string LoginIdentifier
        {
            get => _loginIdentifier;
            set => SetProperty(ref _loginIdentifier, value);
        }

        public RelayCommand ResetPasswordCommand { get; }
        public RelayCommand LoginCommand { get; }

        public event EventHandler<ResetPasswordData>? ResetCodeSent;
        public event EventHandler? LoginRequested;

        public ForgotPasswordViewModel(IAuthService authService)
        {
            _authService = authService;

            ResetPasswordCommand = new RelayCommand(async _ => await SendResetCodeAsync(), _ => CanSend());
            LoginCommand = new RelayCommand(_ => LoginRequested?.Invoke(this, EventArgs.Empty));
        }

        private bool CanSend() => !IsBusy && !string.IsNullOrWhiteSpace(LoginIdentifier);

        private async Task SendResetCodeAsync()
        {
            if (!CanSend()) return;

            IsBusy = true;
            ErrorMessage = string.Empty;

            try
            {
                var result = await _authService.ForgotPasswordAsync(LoginIdentifier);

                if (result.Success)
                {
                    // Передаём reset_token и login_identifier
                    ResetCodeSent?.Invoke(this, new ResetPasswordData
                    {
                        ResetToken = result.UserId, // Токен в поле UserId
                        LoginIdentifier = LoginIdentifier
                    });
                }
                else
                {
                    ErrorMessage = result.Message;
                }
            }
            catch (Exception ex)
            {
                ErrorMessage = $"Error: {ex.Message}";
            }
            finally
            {
                IsBusy = false;
            }
        }
    }

    public class ResetPasswordData
    {
        public string ResetToken { get; set; } = string.Empty;
        public string LoginIdentifier { get; set; } = string.Empty;
    }
}