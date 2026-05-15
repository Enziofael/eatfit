using EatfitDesktop.Helpers;
using EatfitDesktop.Models;
using EatfitDesktop.Services;
using EatfitDesktop.Services.Interfaces;
using System;
using System.Threading.Tasks;
using System.Windows;

namespace EatfitDesktop.ViewModels
{
    public class LoginViewModel : BaseViewModel
    {
        private readonly IAuthService _authService;
        private readonly SessionService _sessionService;

        private string _loginIdentifier = string.Empty;
        public string LoginIdentifier
        {
            get => _loginIdentifier;
            set => SetProperty(ref _loginIdentifier, value);
        }

        public string Password { get; set; } = string.Empty;

        public RelayCommand LoginCommand { get; }
        public RelayCommand RegisterCommand { get; }
        public RelayCommand ForgotPasswordCommand { get; }

        public event EventHandler? LoginSuccessful;
        public event EventHandler? RegisterRequested;
        public event EventHandler? ForgotPasswordRequested;

        public LoginViewModel(IAuthService authService, SessionService sessionService)
        {
            _authService = authService;
            _sessionService = sessionService;

            LoginCommand = new RelayCommand(async _ => await LoginAsync(), _ => CanLogin());
            RegisterCommand = new RelayCommand(_ => RegisterRequested?.Invoke(this, EventArgs.Empty));
            ForgotPasswordCommand = new RelayCommand(_ => ForgotPasswordRequested?.Invoke(this, EventArgs.Empty));
        }

        private bool CanLogin() => !IsBusy && !string.IsNullOrWhiteSpace(LoginIdentifier) && !string.IsNullOrWhiteSpace(Password);

        private async Task LoginAsync()
        {
            if (!CanLogin()) return;

            IsBusy = true;
            ErrorMessage = string.Empty;

            try
            {
                var result = await _authService.LoginAsync(new LoginModel
                {
                    LoginIdentifier = LoginIdentifier,
                    Password = Password
                });

                if (result.Success)
                {
                    _sessionService.SaveSession(result.AccessToken, result.RefreshToken, result.UserId, result.Email, result.Login);
                    LoginSuccessful?.Invoke(this, EventArgs.Empty);
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