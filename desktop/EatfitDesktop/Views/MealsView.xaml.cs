using System;
using System.Collections.Generic;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using Eatfit.V1;
using EatfitDesktop.Services;

namespace EatfitDesktop.Views
{
    public partial class MealsView : UserControl
    {
        private readonly GrpcMealService _mealService;
        private readonly SessionService _sessionService;
        private List<MealData> _meals = new();
        private bool _isLoaded;

        public event Action<string>? MealSelected;
        public event Action<string>? MealEditRequested;
        public event Action? AddMealRequested;

        public MealsView(GrpcMealService mealService, SessionService sessionService)
        {
            InitializeComponent();
            _mealService = mealService;
            _sessionService = sessionService;

            // Отложенная загрузка после полной инициализации UI
            Loaded += (s, e) =>
            {
                if (!_isLoaded)
                {
                    _isLoaded = true;
                    LoadMeals();
                }
            };
        }

        public async void LoadMeals()
        {
            var userId = _sessionService?.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            var sortBy = ((ComboBoxItem)SortBox.SelectedItem)?.Tag?.ToString() ?? "created_at";
            _meals = await _mealService.ListMealsAsync(userId, sortBy);
            MealsList.ItemsSource = _meals;
        }

        private async void Search_Click(object sender, RoutedEventArgs e)
        {
            var userId = _sessionService?.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            var query = SearchBox.Text;
            var sortBy = ((ComboBoxItem)SortBox.SelectedItem)?.Tag?.ToString() ?? "created_at";

            if (string.IsNullOrWhiteSpace(query))
            {
                _meals = await _mealService.ListMealsAsync(userId, sortBy);
            }
            else
            {
                _meals = await _mealService.SearchMealsAsync(userId, query, sortBy);
            }
            MealsList.ItemsSource = _meals;
        }

        private void SortBox_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            if (_isLoaded)
            {
                LoadMeals();
            }
        }

        private void Meal_Click(object sender, MouseButtonEventArgs e)
        {
            if ((sender as FrameworkElement)?.DataContext is MealData meal)
            {
                MealSelected?.Invoke(meal.MealId);
            }
        }

        private void Edit_Click(object sender, RoutedEventArgs e)
        {
            if ((sender as Button)?.DataContext is MealData meal)
            {
                MealEditRequested?.Invoke(meal.MealId);
            }
        }

        private async void Delete_Click(object sender, RoutedEventArgs e)
        {
            if ((sender as Button)?.DataContext is MealData meal)
            {
                var result = MessageBox.Show($"Delete \"{meal.Name}\"?", "Confirm",
                    MessageBoxButton.YesNo, MessageBoxImage.Question);
                if (result == MessageBoxResult.Yes)
                {
                    var success = await _mealService.DeleteMealAsync(meal.MealId);
                    if (success) LoadMeals();
                }
            }
        }

        private void SearchBox_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.Key == Key.Enter)
            {
                Search_Click(sender, e);
            }
        }


        private void AddMeal_Click(object sender, RoutedEventArgs e)
        {
            AddMealRequested?.Invoke();
        }
    }
}