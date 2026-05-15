using System;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;

namespace EatfitDesktop.Views
{
    public partial class MainView : UserControl
    {
        public event Action? LogoutRequested;

        public MainView()
        {
            InitializeComponent();
            ShowFeed();
        }

        private void FeedBtn_Click(object sender, RoutedEventArgs e) => ShowFeed();
        private void DiaryBtn_Click(object sender, RoutedEventArgs e) => ShowDiary();
        private void MealsBtn_Click(object sender, RoutedEventArgs e) => ShowMeals();
        private void MessagesBtn_Click(object sender, RoutedEventArgs e) => ShowMessages();
        private void ProfileBtn_Click(object sender, RoutedEventArgs e) => ShowProfile();

        private void ShowFeed()
        {
            ContentArea.Content = CreatePlaceholder("Feed", "Publications from other users will be here");
        }

        private void ShowDiary()
        {
            ContentArea.Content = CreatePlaceholder("Diary", "Training & meal tracking will be here");
        }

        private void ShowMeals()
        {
            ContentArea.Content = CreatePlaceholder("Meals", "Create and browse meals here");
        }

        private void ShowMessages()
        {
            ContentArea.Content = CreatePlaceholder("Messages", "Chats will be here");
        }

        private void ShowProfile()
        {
            var stackPanel = new StackPanel
            {
                HorizontalAlignment = HorizontalAlignment.Center,
                VerticalAlignment = VerticalAlignment.Center
            };

            stackPanel.Children.Add(new TextBlock
            {
                Text = "Profile",
                FontSize = 32,
                FontWeight = FontWeights.Bold,
                Foreground = new SolidColorBrush(Color.FromRgb(0x66, 0x7E, 0xEA)),
                HorizontalAlignment = HorizontalAlignment.Center,
                Margin = new Thickness(0, 0, 0, 8)
            });

            stackPanel.Children.Add(new TextBlock
            {
                Text = "Your profile will be here",
                FontSize = 16,
                Foreground = new SolidColorBrush(Color.FromRgb(0x66, 0x66, 0x66)),
                HorizontalAlignment = HorizontalAlignment.Center,
                Margin = new Thickness(0, 0, 0, 24)
            });

            var logoutBtn = new Button
            {
                Content = "Logout",
                Width = 200,
                Height = 40,
                Style = Application.Current.Resources["PrimaryButton"] as Style
            };
            logoutBtn.Click += (s, e) => LogoutRequested?.Invoke();

            stackPanel.Children.Add(logoutBtn);

            ContentArea.Content = stackPanel;
        }

        private UIElement CreatePlaceholder(string title, string subtitle)
        {
            var stackPanel = new StackPanel
            {
                HorizontalAlignment = HorizontalAlignment.Center,
                VerticalAlignment = VerticalAlignment.Center
            };

            stackPanel.Children.Add(new TextBlock
            {
                Text = title,
                FontSize = 32,
                FontWeight = FontWeights.Bold,
                Foreground = new SolidColorBrush(Color.FromRgb(0x66, 0x7E, 0xEA)),
                HorizontalAlignment = HorizontalAlignment.Center,
                Margin = new Thickness(0, 0, 0, 8)
            });

            stackPanel.Children.Add(new TextBlock
            {
                Text = subtitle,
                FontSize = 16,
                Foreground = new SolidColorBrush(Color.FromRgb(0x66, 0x66, 0x66)),
                HorizontalAlignment = HorizontalAlignment.Center
            });

            return stackPanel;
        }
    }
}