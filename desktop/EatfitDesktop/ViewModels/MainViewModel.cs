using System;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Windows.Input;
using EatfitDesktop.Helpers;

namespace EatfitDesktop.ViewModels
{
    public class MainViewModel : INotifyPropertyChanged
    {
        private string _welcomeMessage = "Welcome to Eatfit!";
        public string WelcomeMessage
        {
            get => _welcomeMessage;
            set
            {
                _welcomeMessage = value;
                OnPropertyChanged();
            }
        }

        public ICommand LogoutCommand { get; }

        public event EventHandler? LogoutRequested;

        public MainViewModel()
        {
            LogoutCommand = new RelayCommand(_ => Logout());
        }

        private void Logout()
        {
            LogoutRequested?.Invoke(this, EventArgs.Empty);
        }

        public event PropertyChangedEventHandler? PropertyChanged;
        protected void OnPropertyChanged([CallerMemberName] string? propertyName = null)
            => PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
    }
}