using System;
using System.Windows;
using System.Windows.Controls;
using Eatfit.V1;
using EatfitDesktop.Services;

namespace EatfitDesktop.Views
{
    public partial class SettingsView : UserControl
    {
        private readonly GrpcProfileService _profileService;
        private readonly GrpcAuthService _authService;
        private readonly SessionService _sessionService;

        public event Action? BackRequested;

        public SettingsView(GrpcProfileService profileService, GrpcAuthService authService, SessionService sessionService)
        {
            InitializeComponent();
            _profileService = profileService;
            _authService = authService;
            _sessionService = sessionService;
        }

        public async void LoadProfile()
        {
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            var profile = await _profileService.GetProfileAsync(userId);
            if (profile == null) return;

            FirstNameBox.Text = profile.FirstName;
            LastNameBox.Text = profile.LastName;
            HeightBox.Text = profile.Height > 0 ? profile.Height.ToString() : "";

            // DatePicker
            if (!string.IsNullOrEmpty(profile.BirthDate) && DateTime.TryParse(profile.BirthDate, out var date))
            {
                BirthDatePicker.SelectedDate = date;
            }
            else
            {
                BirthDatePicker.SelectedDate = null;
            }

            foreach (ComboBoxItem item in GenderBox.Items)
            {
                if (item.Tag?.ToString() == profile.Gender)
                {
                    item.IsSelected = true;
                    break;
                }
            }
        }


        private async void SaveProfile_Click(object sender, RoutedEventArgs e)
        {
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            double? height = null;
            if (double.TryParse(HeightBox.Text, out double h) && h > 0)
                height = h;

            string? gender = null;
            if (GenderBox.SelectedItem is ComboBoxItem item && item.Tag != null)
                gender = item.Tag.ToString();

            string birthDate = BirthDatePicker.SelectedDate?.ToString("yyyy-MM-dd") ?? "";

            var success = await _profileService.UpdateProfileAsync(
                userId,
                FirstNameBox.Text,
                LastNameBox.Text,
                height,
                birthDate,
                gender ?? ""
            );

            if (success)
            {
                ProfileErrorText.Text = "Profile saved!";
                ProfileErrorText.Foreground = new System.Windows.Media.SolidColorBrush(
                    System.Windows.Media.Color.FromRgb(0x11, 0x99, 0x8E));
            }
            else
            {
                ProfileErrorText.Text = "Failed to save profile";
            }
        }

        private async void ChangeLogin_Click(object sender, RoutedEventArgs e)
        {
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            var newLogin = NewLoginBox.Text;
            var password = LoginPasswordBox.Password;

            if (string.IsNullOrWhiteSpace(newLogin) || string.IsNullOrWhiteSpace(password))
            {
                LoginErrorText.Text = "All fields are required";
                return;
            }

            var result = await _authService.ChangeLoginAsync(userId, newLogin, password);

            if (result.Success)
            {
                LoginErrorText.Text = "Login changed successfully!";
                LoginErrorText.Foreground = new System.Windows.Media.SolidColorBrush(
                    System.Windows.Media.Color.FromRgb(0x11, 0x99, 0x8E));
                NewLoginBox.Clear();
                LoginPasswordBox.Clear();
            }
            else
            {
                LoginErrorText.Text = result.Message;
            }
        }

        private async void ChangePassword_Click(object sender, RoutedEventArgs e)
        {
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            var currentPassword = CurrentPasswordBox.Password;
            var newPassword = NewPasswordBox.Password;
            var confirmPassword = ConfirmPasswordBox.Password;

            if (string.IsNullOrWhiteSpace(currentPassword) ||
                string.IsNullOrWhiteSpace(newPassword) ||
                string.IsNullOrWhiteSpace(confirmPassword))
            {
                PasswordErrorText.Text = "All fields are required";
                return;
            }

            var result = await _authService.ChangePasswordAsync(userId, currentPassword, newPassword, confirmPassword);

            if (result.Success)
            {
                PasswordErrorText.Text = "Password changed successfully!";
                PasswordErrorText.Foreground = new System.Windows.Media.SolidColorBrush(
                    System.Windows.Media.Color.FromRgb(0x11, 0x99, 0x8E));
                CurrentPasswordBox.Clear();
                NewPasswordBox.Clear();
                ConfirmPasswordBox.Clear();
            }
            else
            {
                PasswordErrorText.Text = result.Message;
            }
        }

        private void Back_Click(object sender, RoutedEventArgs e)
        {
            BackRequested?.Invoke();
        }
    }
}