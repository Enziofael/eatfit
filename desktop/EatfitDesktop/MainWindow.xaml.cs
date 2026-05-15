using System;
using System.Windows;
using EatfitDesktop.Services;
using EatfitDesktop.ViewModels;

namespace EatfitDesktop.Views
{
    public partial class MainWindow : Window
    {
        private readonly GrpcAuthService _authService;
        private readonly SessionService _sessionService;

        public MainWindow()
        {
            InitializeComponent();

            _authService = new GrpcAuthService();
            _sessionService = new SessionService();

            if (_sessionService.IsLoggedIn)
            {
                ShowMainView();
            }
            else
            {
                ShowLoginView();
            }
        }

        private void ShowLoginView()
        {
            var vm = new LoginViewModel(_authService, _sessionService);
            vm.LoginSuccessful += (s, e) => ShowMainView();
            vm.RegisterRequested += (s, e) => ShowRegisterView();
            vm.ForgotPasswordRequested += (s, e) => ShowForgotPasswordView();

            ContentArea.Content = new LoginView { DataContext = vm };
        }

        private void ShowRegisterView()
        {
            var vm = new RegisterViewModel(_authService);
            vm.RegistrationCompleted += (s, result) => ShowVerifyEmailView(result.UserId, result.Email, result.Login);
            vm.LoginRequested += (s, e) => ShowLoginView();

            ContentArea.Content = new RegisterView { DataContext = vm };
        }

        private void ShowVerifyEmailView(string userId, string email, string login)
        {
            var vm = new VerifyEmailViewModel(_authService);
            vm.SetUserData(userId, email, login);
            vm.VerificationCompleted += (s, result) =>
            {
                MessageBox.Show("Email verified! Please login.", "Success",
                    MessageBoxButton.OK, MessageBoxImage.Information);
                ShowLoginView();
            };
            vm.BackRequested += (s, e) => ShowRegisterView();

            ContentArea.Content = new VerifyEmailView { DataContext = vm };
        }

        private string? _resetToken;

        private void ShowForgotPasswordView()
        {
            var vm = new ForgotPasswordViewModel(_authService);
            vm.ResetCodeSent += (s, data) =>
            {
                _resetToken = data.ResetToken;

                if (!string.IsNullOrEmpty(_resetToken))
                {
                    ShowResetPasswordView(_resetToken);
                }
                else
                {
                    MessageBox.Show("If account exists, reset code sent to email.",
                        "Information", MessageBoxButton.OK, MessageBoxImage.Information);
                }
            };
            vm.LoginRequested += (s, e) => ShowLoginView();

            ContentArea.Content = new ForgotPasswordView { DataContext = vm };
        }

        private void ShowResetPasswordView(string resetToken)
        {
            var vm = new ResetPasswordViewModel(_authService);
            vm.SetResetToken(resetToken);
            vm.PasswordResetSuccess += (s, e) =>
            {
                MessageBox.Show("Password reset successfully!", "Success",
                    MessageBoxButton.OK, MessageBoxImage.Information);
                ShowLoginView();
            };
            vm.LoginRequested += (s, e) => ShowLoginView();

            ContentArea.Content = new ResetPasswordView { DataContext = vm };
        }

        private void ShowMainView()
        {
            var vm = new MainViewModel();
            vm.LogoutRequested += async (s, e) =>
            {
                if (!string.IsNullOrEmpty(_sessionService.RefreshToken))
                {
                    await _authService.LogoutAsync(_sessionService.RefreshToken);
                }

                _sessionService.ClearSession();
                ShowLoginView();
            };

            ContentArea.Content = new MainView { DataContext = vm };
        }
    }
}