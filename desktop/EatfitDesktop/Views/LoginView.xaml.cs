using System.Windows.Controls;
using EatfitDesktop.ViewModels;

namespace EatfitDesktop.Views
{
    public partial class LoginView : UserControl
    {
        public LoginView()
        {
            InitializeComponent();

            // Привязываем пароль через код, так как PasswordBox не поддерживает MVVM напрямую
            PasswordBox.PasswordChanged += (s, e) =>
            {
                if (DataContext is LoginViewModel vm)
                    vm.Password = PasswordBox.Password;
            };
        }
    }
}