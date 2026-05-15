using System;
using System.Collections.Generic;
using System.Linq;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;
using Eatfit.V1;
using EatfitDesktop.Services;

namespace EatfitDesktop.Views
{
    public partial class MealEditorView : UserControl
    {
        private readonly GrpcMealService _mealService;
        private readonly SessionService _sessionService;
        private string? _mealId;
        private List<MealData> _allMeals = new(); // для поиска при добавлении компонента

        public event Action? Saved;
        public event Action? Cancelled;

        public MealEditorView(GrpcMealService mealService, SessionService sessionService)
        {
            InitializeComponent();
            _mealService = mealService;
            _sessionService = sessionService;
        }

        public async void LoadMeal(string? mealId)
        {
            _mealId = mealId;

            // Загружаем все блюда пользователя для поиска компонентов
            _allMeals = await _mealService.ListMealsAsync(_sessionService.UserId!);

            if (mealId != null)
            {
                TitleText.Text = "Edit Meal";
                var meal = await _mealService.GetMealAsync(mealId);
                if (meal != null)
                {
                    NameBox.Text = meal.Name;
                    DescriptionBox.Text = meal.Description;
                    RecipeBox.Text = meal.Recipe;
                    CaloriesBox.Text = meal.Calories.ToString("F1");
                    ProteinsBox.Text = meal.Proteins.ToString("F1");
                    FatsBox.Text = meal.Fats.ToString("F1");
                    CarbsBox.Text = meal.Carbs.ToString("F1");
                    WaterBox.Text = meal.Water.ToString("F1");

                    if (meal.Components.Count > 0)
                    {
                        foreach (var comp in meal.Components)
                        {
                            AddComponentRow(comp.ComponentMealId, comp.ComponentName, comp.Amount);
                        }
                    }
                }
            }
        }

        private async void Save_Click(object sender, RoutedEventArgs e)
        {
            var name = NameBox.Text.Trim();
            if (string.IsNullOrEmpty(name))
            {
                ErrorText.Text = "Name is required";
                return;
            }

            var calories = ParseDouble(CaloriesBox.Text);
            var proteins = ParseDouble(ProteinsBox.Text);
            var fats = ParseDouble(FatsBox.Text);
            var carbs = ParseDouble(CarbsBox.Text);
            var water = ParseDouble(WaterBox.Text);

            // Собираем компоненты из UI
            var components = new List<MealComponentInput>();
            foreach (StackPanel panel in ComponentsList.Items)
            {
                var componentMealId = panel.Tag as string;
                if (string.IsNullOrEmpty(componentMealId)) continue;

                var amountBox = panel.Children.OfType<TextBox>().FirstOrDefault();
                var amount = amountBox != null && double.TryParse(amountBox.Text, out var a) ? a : 100;

                components.Add(new MealComponentInput
                {
                    ComponentMealId = componentMealId,
                    Amount = amount
                });
            }

            bool success;

            if (_mealId == null)
            {
                success = await _mealService.CreateMealAsync(
                    _sessionService.UserId!, name, DescriptionBox.Text, RecipeBox.Text,
                    "", calories, proteins, fats, carbs, water, components);
            }
            else
            {
                success = await _mealService.UpdateMealAsync(
                    _mealId, name, DescriptionBox.Text, RecipeBox.Text,
                    "", calories, proteins, fats, carbs, water, components);
            }

            if (success)
            {
                Saved?.Invoke();
            }
            else
            {
                ErrorText.Text = "Failed to save meal";
            }
        }

        private void Cancel_Click(object sender, RoutedEventArgs e)
        {
            Cancelled?.Invoke();
        }

        private void AddComponent_Click(object sender, RoutedEventArgs e)
        {
            ShowComponentPicker();
        }

        private void ShowComponentPicker()
        {
            // Собираем ID уже добавленных компонентов
            var addedIds = new HashSet<string>();
            foreach (StackPanel panel in ComponentsList.Items)
            {
                var id = panel.Tag as string;
                if (!string.IsNullOrEmpty(id))
                    addedIds.Add(id);
            }

            // Фильтруем: убираем текущее блюдо и уже добавленные
            var available = _allMeals
                .Where(m => m.MealId != _mealId && !addedIds.Contains(m.MealId))
                .ToList();

            if (available.Count == 0)
            {
                MessageBox.Show("No other meals available to add as component.", "Info",
                    MessageBoxButton.OK, MessageBoxImage.Information);
                return;
            }

            var dialog = new Window
            {
                Title = "Select Component",
                Width = 400,
                Height = 450,
                WindowStartupLocation = WindowStartupLocation.CenterOwner,
                Owner = Window.GetWindow(this),
                ResizeMode = ResizeMode.NoResize
            };

            var grid = new Grid { Margin = new Thickness(15) };
            grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });
            grid.RowDefinitions.Add(new RowDefinition { Height = new GridLength(1, GridUnitType.Star) });
            grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });

            var searchBox = new TextBox { Margin = new Thickness(0, 0, 0, 10), Height = 30 };
            searchBox.TextChanged += (s, e) =>
            {
                var query = searchBox.Text.ToLower();
                var filtered = available
                    .Where(m => m.Name.ToLower().Contains(query))
                    .ToList();
                ((ListBox)grid.Children[1]).ItemsSource = filtered;
            };
            Grid.SetRow(searchBox, 0);
            grid.Children.Add(searchBox);

            var listBox = new ListBox();
            listBox.ItemsSource = available;
            listBox.DisplayMemberPath = "Name";
            listBox.SelectionChanged += (s, e) =>
            {
                var okBtn = (Button)((StackPanel)grid.Children[2]).Children[1];
                okBtn.IsEnabled = listBox.SelectedItem != null;
            };
            Grid.SetRow(listBox, 1);
            grid.Children.Add(listBox);

            var buttons = new StackPanel
            {
                Orientation = Orientation.Horizontal,
                HorizontalAlignment = HorizontalAlignment.Right,
                Margin = new Thickness(0, 10, 0, 0)
            };
            var cancelBtn = new Button { Content = "Cancel", Width = 80, Margin = new Thickness(0, 0, 10, 0) };
            cancelBtn.Click += (_, _) => dialog.Close();
            var okBtn = new Button
            {
                Content = "Add",
                Width = 80,
                IsEnabled = false,
                Style = Application.Current.Resources["PrimaryButton"] as Style
            };
            okBtn.Click += (_, _) =>
            {
                if (listBox.SelectedItem is MealData meal)
                {
                    AddComponentRow(meal.MealId, meal.Name, 100);
                }
                dialog.Close();
            };
            buttons.Children.Add(cancelBtn);
            buttons.Children.Add(okBtn);
            Grid.SetRow(buttons, 2);
            grid.Children.Add(buttons);

            dialog.Content = grid;
            searchBox.Focus();
            dialog.ShowDialog();
        }

        private void AddComponentRow(string componentMealId, string name, double amount)
        {
            var panel = new StackPanel { Orientation = Orientation.Horizontal, Margin = new Thickness(0, 5, 0, 0) };

            var nameText = new TextBlock
            {
                Text = name,
                Width = 180,
                VerticalAlignment = VerticalAlignment.Center,
                Margin = new Thickness(0, 0, 10, 0),
                TextTrimming = TextTrimming.CharacterEllipsis
            };

            var amountBox = new TextBox
            {
                Text = amount.ToString(),
                Width = 60,
                Margin = new Thickness(0, 0, 5, 0)
            };

            var unitText = new TextBlock
            {
                Text = "g",
                VerticalAlignment = VerticalAlignment.Center,
                Margin = new Thickness(0, 0, 10, 0)
            };

            var deleteBtn = new Button
            {
                Content = "✕",
                Width = 25,
                Height = 25,
                Foreground = new SolidColorBrush(Color.FromRgb(0xe7, 0x4c, 0x3c)),
                Background = Brushes.Transparent,
                BorderThickness = new Thickness(0),
                Cursor = System.Windows.Input.Cursors.Hand
            };
            deleteBtn.Click += (_, _) => ComponentsList.Items.Remove(panel);

            panel.Tag = componentMealId; // храним ID в Tag
            panel.Children.Add(nameText);
            panel.Children.Add(amountBox);
            panel.Children.Add(unitText);
            panel.Children.Add(deleteBtn);

            ComponentsList.Items.Add(panel);
        }

        private static double ParseDouble(string s) => double.TryParse(s, out var d) ? d : 0;
    }
}