using System.Windows.Controls;
using EatfitDesktop.ViewModels;

namespace EatfitDesktop.Views
{
    public partial class ResetPasswordView : UserControl
    {
        public ResetPasswordView()
        {
            InitializeComponent();

            NewPasswordBox.PasswordChanged += (s, e) =>
            {
                if (DataContext is ResetPasswordViewModel vm)
                    vm.NewPassword = NewPasswordBox.Password;
            };

            ConfirmPasswordBox.PasswordChanged += (s, e) =>
            {
                if (DataContext is ResetPasswordViewModel vm)
                    vm.PasswordConfirmation = ConfirmPasswordBox.Password;
            };
        }
    }
}