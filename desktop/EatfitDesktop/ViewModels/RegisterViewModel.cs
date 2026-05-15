using System;
using System.Threading.Tasks;
using EatfitDesktop.Helpers;
using EatfitDesktop.Models;
using EatfitDesktop.Services.Interfaces;

namespace EatfitDesktop.ViewModels
{
    public class RegisterViewModel : BaseViewModel
    {
        private readonly IAuthService _authService;

        private string _email = string.Empty;
        public string Email
        {
            get => _email;
            set => SetProperty(ref _email, value);
        }

        private string _login = string.Empty;
        public string Login
        {
            get => _login;
            set => SetProperty(ref _login, value);
        }

        public string Password { get; set; } = string.Empty;
        public string PasswordConfirmation { get; set; } = string.Empty;

        public RelayCommand RegisterCommand { get; }
        public RelayCommand LoginCommand { get; }

        public event EventHandler<RegisterResult>? RegistrationCompleted;
        public event EventHandler? LoginRequested;

        public RegisterViewModel(IAuthService authService)
        {
            _authService = authService;

            RegisterCommand = new RelayCommand(async _ => await RegisterAsync(), _ => CanRegister());
            LoginCommand = new RelayCommand(_ => LoginRequested?.Invoke(this, EventArgs.Empty));
        }

        private bool CanRegister() => !IsBusy && 
            !string.IsNullOrWhiteSpace(Email) && 
            !string.IsNullOrWhiteSpace(Login) &&
            !string.IsNullOrWhiteSpace(Password) &&
            !string.IsNullOrWhiteSpace(PasswordConfirmation);

        private async Task RegisterAsync()
        {
            if (!CanRegister()) return;

            IsBusy = true;
            ErrorMessage = string.Empty;

            try
            {
                var result = await _authService.RegisterAsync(new RegisterModel
                {
                    Email = Email,
                    Login = Login,
                    Password = Password,
                    PasswordConfirmation = PasswordConfirmation
                });

                if (result.Success)
                {
                    RegistrationCompleted?.Invoke(this, new RegisterResult { UserId = result.UserId, Email = Email, Login = Login });
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

    public class RegisterResult
    {
        public string UserId { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        public string Login { get; set; } = string.Empty;
    }
}