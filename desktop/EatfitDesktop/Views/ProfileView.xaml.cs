using System;
using System.Windows;
using System.Windows.Controls;

namespace EatfitDesktop.Views
{
    public partial class ProfileView : UserControl
    {
        public event Action? LogoutRequested;
        public event Action? SettingsRequested;

        public ProfileView()
        {
            InitializeComponent();
        }

        private void Logout_Click(object sender, RoutedEventArgs e)
        {
            var result = MessageBox.Show("Are you sure you want to logout?", "Logout",
                MessageBoxButton.YesNo, MessageBoxImage.Question);

            if (result == MessageBoxResult.Yes)
            {
                LogoutRequested?.Invoke();
            }
        }

        private void Settings_Click(object sender, RoutedEventArgs e)
        {
            SettingsRequested?.Invoke();
        }
    }
}