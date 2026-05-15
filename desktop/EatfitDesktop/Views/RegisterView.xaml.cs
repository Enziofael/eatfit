using System.Windows.Controls;
using EatfitDesktop.ViewModels;

namespace EatfitDesktop.Views
{
    public partial class RegisterView : UserControl
    {
        public RegisterView()
        {
            InitializeComponent();

            PasswordBox.PasswordChanged += (s, e) =>
            {
                if (DataContext is RegisterViewModel vm)
                    vm.Password = PasswordBox.Password;
            };

            ConfirmPasswordBox.PasswordChanged += (s, e) =>
            {
                if (DataContext is RegisterViewModel vm)
                    vm.PasswordConfirmation = ConfirmPasswordBox.Password;
            };
        }
    }
}