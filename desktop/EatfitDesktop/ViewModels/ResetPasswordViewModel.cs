using System;
using System.Threading.Tasks;
using EatfitDesktop.Helpers;
using EatfitDesktop.Services.Interfaces;

namespace EatfitDesktop.ViewModels
{
    public class ResetPasswordViewModel : BaseViewModel
    {
        private readonly IAuthService _authService;
        private string _userId = string.Empty;

        private string _verificationCode = string.Empty;
        public string VerificationCode
        {
            get => _verificationCode;
            set => SetProperty(ref _verificationCode, value);
        }

        public string NewPassword { get; set; } = string.Empty;
        public string PasswordConfirmation { get; set; } = string.Empty;

        public RelayCommand ResetPasswordCommand { get; }
        public RelayCommand LoginCommand { get; }

        public event EventHandler? PasswordResetSuccess;
        public event EventHandler? LoginRequested;

        public ResetPasswordViewModel(IAuthService authService)
        {
            _authService = authService;

            ResetPasswordCommand = new RelayCommand(async _ => await ResetPasswordAsync(), _ => CanReset());
            LoginCommand = new RelayCommand(_ => LoginRequested?.Invoke(this, EventArgs.Empty));
        }

        public void SetUserId(string userId) => _userId = userId;

        private bool CanReset() => !IsBusy &&
            !string.IsNullOrWhiteSpace(VerificationCode) &&
            !string.IsNullOrWhiteSpace(NewPassword) &&
            NewPassword == PasswordConfirmation;

        private string _resetToken = string.Empty;

        public void SetResetToken(string resetToken)
        {
            _resetToken = resetToken;
        }

        private async Task ResetPasswordAsync()
        {
            if (!CanReset()) return;

            IsBusy = true;
            ErrorMessage = string.Empty;

            try
            {
                var result = await _authService.ResetPasswordAsync(_resetToken, VerificationCode, NewPassword);

                if (result.Success)
                {
                    PasswordResetSuccess?.Invoke(this, EventArgs.Empty);
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
}