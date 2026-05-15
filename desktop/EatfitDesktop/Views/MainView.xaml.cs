using Eatfit.V1;
using EatfitDesktop.Services;
using EatfitDesktop.ViewModels;
using System;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;

namespace EatfitDesktop.Views
{
    public partial class MainView : UserControl
    {
        public event Action? LogoutRequested;
        private readonly GrpcMealService _mealService = new();
        private readonly SessionService _sessionService;
        private readonly GrpcDiaryService _diaryService = new();

        public MainView(SessionService sessionService)
        {
            InitializeComponent();
            _sessionService = sessionService;
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
            var profileService = new GrpcProfileService();
            var view = new DiaryView(_diaryService, _mealService, profileService, _sessionService);
            view.MealViewRequested += (mealId) => ShowMealDetail(mealId);
            view.LoadDiary();
            ContentArea.Content = view;
        }

        private void ShowMeals()
        {
            var view = new MealsView(_mealService, _sessionService);
            view.LoadMeals();
            view.AddMealRequested += () => ShowMealEditor(null);
            view.MealEditRequested += (mealId) => ShowMealEditor(mealId);
            view.MealSelected += (mealId) => ShowMealDetail(mealId);
            ContentArea.Content = view;
        }

        private void ShowMealEditor(string? mealId)
        {
            var view = new MealEditorView(_mealService, _sessionService);
            view.LoadMeal(mealId);
            view.Saved += () =>
            {
                ShowMeals();
            };
            view.Cancelled += () =>
            {
                ShowMeals();
            };
            ContentArea.Content = view;
        }

        private void ShowMessages()
        {
            ContentArea.Content = CreatePlaceholder("Messages", "Chats will be here");
        }

        private void ShowProfile()
        {
            var profileService = new GrpcProfileService();

            var vm = new ProfileViewModel(profileService, _sessionService);
            var view = new ProfileView { DataContext = vm };

            view.LogoutRequested += () => LogoutRequested?.Invoke();
            view.SettingsRequested += () => ShowSettings();

            _ = vm.LoadProfileAsync();

            ContentArea.Content = view;
        }

        private void ShowSettings()
        {
            var profileService = new GrpcProfileService();
            var authService = new GrpcAuthService();

            var view = new SettingsView(profileService, authService, _sessionService);
            view.BackRequested += () => ShowProfile();
            view.LoadProfile();

            ContentArea.Content = view;
        }

        private void ShowMealDetail(string mealId)
        {
            var view = new MealDetailView(_mealService);
            view.LoadMeal(mealId);
            view.BackRequested += () => ShowMeals();
            ContentArea.Content = view;
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