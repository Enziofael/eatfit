using System.Windows;
using System.Windows.Controls;
using Eatfit.V1;
using EatfitDesktop.Services;

namespace EatfitDesktop.Views
{
    public partial class MealDetailView : UserControl
    {
        private readonly GrpcMealService _mealService;

        public event Action? BackRequested;

        public MealDetailView(GrpcMealService mealService)
        {
            InitializeComponent();
            _mealService = mealService;
        }

        public async void LoadMeal(string mealId)
        {
            var meal = await _mealService.GetMealAsync(mealId);
            if (meal == null) return;

            NameText.Text = meal.Name;
            KcalText.Text = $"{meal.Calories:F0} kcal per 100g";
            CaloriesText.Text = meal.Calories.ToString("F0");
            ProteinsText.Text = meal.Proteins.ToString("F0");
            FatsText.Text = meal.Fats.ToString("F0");
            CarbsText.Text = meal.Carbs.ToString("F0");
            WaterText.Text = meal.Water.ToString("F0");

            // Состав
            if (meal.Components.Count > 0)
            {
                ComponentsPanel.Visibility = Visibility.Visible;
                foreach (var comp in meal.Components)
                {
                    var item = new Border
                    {
                        Padding = new Thickness(0, 5, 0, 5),
                        BorderBrush = System.Windows.Media.Brushes.LightGray,
                        BorderThickness = new Thickness(0, 0, 0, 1)
                    };
                    var stack = new StackPanel();
                    stack.Children.Add(new TextBlock
                    {
                        Text = $"{comp.ComponentName} — {comp.Amount:F0}g",
                        FontWeight = FontWeights.Medium
                    });
                    stack.Children.Add(new TextBlock
                    {
                        Text = $"Cal: {comp.Calories:F0} | P: {comp.Proteins:F0} | F: {comp.Fats:F0} | C: {comp.Carbs:F0}",
                        FontSize = 11,
                        Foreground = new System.Windows.Media.SolidColorBrush(
                            System.Windows.Media.Color.FromRgb(0x99, 0x99, 0x99))
                    });
                    item.Child = stack;
                    ComponentsList.Items.Add(item);
                }
            }

            // Описание
            if (!string.IsNullOrWhiteSpace(meal.Description))
            {
                DescriptionPanel.Visibility = Visibility.Visible;
                DescriptionText.Text = meal.Description;
            }

            // Рецепт
            if (!string.IsNullOrWhiteSpace(meal.Recipe))
            {
                RecipePanel.Visibility = Visibility.Visible;
                RecipeText.Text = meal.Recipe;
            }
        }

        private void Back_Click(object sender, RoutedEventArgs e)
        {
            BackRequested?.Invoke();
        }
    }
}