using System;
using System.Threading.Tasks;
using EatfitDesktop.Helpers;
using EatfitDesktop.Services.Interfaces;

namespace EatfitDesktop.ViewModels
{
    public class VerifyEmailViewModel : BaseViewModel
    {
        private readonly IAuthService _authService;
        private string _userId = string.Empty;
        private string _email = string.Empty;
        private string _login = string.Empty;

        private string _verificationCode = string.Empty;
        public string VerificationCode
        {
            get => _verificationCode;
            set => SetProperty(ref _verificationCode, value);
        }

        public RelayCommand VerifyCommand { get; }
        public RelayCommand BackCommand { get; }

        public event EventHandler<VerificationResult>? VerificationCompleted;
        public event EventHandler? BackRequested;

        public VerifyEmailViewModel(IAuthService authService)
        {
            _authService = authService;

            VerifyCommand = new RelayCommand(async _ => await VerifyAsync(), _ => CanVerify());
            BackCommand = new RelayCommand(async _ => await BackAsync());
        }

        public void SetUserData(string userId, string email, string login)
        {
            _userId = userId;
            _email = email;
            _login = login;
        }

        private bool CanVerify() => !IsBusy && !string.IsNullOrWhiteSpace(VerificationCode) && VerificationCode.Length == 6;

        private async Task VerifyAsync()
        {
            if (!CanVerify()) return;

            IsBusy = true;
            ErrorMessage = string.Empty;

            try
            {
                var result = await _authService.VerifyEmailAsync(_userId, VerificationCode);

                if (result.Success)
                {
                    VerificationCompleted?.Invoke(this, new VerificationResult
                    {
                        Success = true,
                        Email = _email,
                        Login = _login
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

        private async Task BackAsync()
        {
            // Удаляем неподтверждённый аккаунт
            await _authService.DeleteAccountAsync(_userId);
            BackRequested?.Invoke(this, EventArgs.Empty);
        }
    }

    public class VerificationResult
    {
        public bool Success { get; set; }
        public string Email { get; set; } = string.Empty;
        public string Login { get; set; } = string.Empty;
    }
}