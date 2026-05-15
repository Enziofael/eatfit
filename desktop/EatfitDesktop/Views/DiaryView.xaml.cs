using System;
using System.Collections.Generic;
using System.Linq;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Shapes;
using Eatfit.V1;
using EatfitDesktop.Services;

namespace EatfitDesktop.Views
{
    public partial class DiaryView : UserControl
    {
        private readonly GrpcDiaryService _diaryService;
        private readonly GrpcMealService _mealService;
        private readonly GrpcProfileService _profileService;
        private readonly SessionService _sessionService;
        private List<ConsumptionGroup> _groups = new();
        private List<MealData> _allMeals = new();
        private DateTime _currentMonth = DateTime.Today;

        public DiaryView(GrpcDiaryService diaryService, GrpcMealService mealService,
            GrpcProfileService profileService, SessionService sessionService)
        {
            InitializeComponent();
            _diaryService = diaryService;
            _mealService = mealService;
            _profileService = profileService;
            _sessionService = sessionService;
            DayPicker.SelectedDate = DateTime.Today;
        }

        public async void LoadDiary()
        {
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            _groups = await _diaryService.ListConsumptionsAsync(userId);
            _allMeals = await _mealService.ListMealsAsync(userId);
            RenderConsumptions();
            await UpdateProgress();
            await DrawWeightChart();
        }

        private void RenderConsumptions()
        {
            ConsumptionList.Items.Clear();
            if (_groups.Count == 0)
            {
                ConsumptionList.Items.Add(new TextBlock
                {
                    Text = "No consumption records yet",
                    Foreground = new SolidColorBrush(Color.FromRgb(0x99, 0x99, 0x99)),
                    Margin = new Thickness(0, 20, 0, 0),
                    HorizontalAlignment = HorizontalAlignment.Center
                });
                return;
            }

            foreach (var group in _groups)
            {
                var groupPanel = new StackPanel { Margin = new Thickness(0, 0, 0, 15) };
                groupPanel.Children.Add(new TextBlock
                {
                    Text = group.Date,
                    FontWeight = FontWeights.SemiBold,
                    FontSize = 14,
                    Margin = new Thickness(0, 0, 0, 8)
                });

                foreach (var record in group.Records)
                {
                    groupPanel.Children.Add(CreateConsumptionCard(record));
                }
                ConsumptionList.Items.Add(groupPanel);
            }
        }

        private Border CreateConsumptionCard(ConsumptionRecord record)
        {
            var border = new Border
            {
                Background = Brushes.White,
                CornerRadius = new CornerRadius(8),
                Padding = new Thickness(12),
                Margin = new Thickness(0, 0, 0, 6)
            };

            var grid = new Grid();
            grid.ColumnDefinitions.Add(new ColumnDefinition());
            grid.ColumnDefinitions.Add(new ColumnDefinition { Width = GridLength.Auto });
            grid.ColumnDefinitions.Add(new ColumnDefinition { Width = GridLength.Auto });

            var info = new StackPanel();
            info.Children.Add(new TextBlock { Text = $"{record.MealName} — {record.Amount:F0}g", FontWeight = FontWeights.Medium });
            info.Children.Add(new TextBlock
            {
                Text = $"Cal: {record.Calories:F0} | P: {record.Proteins:F0} | F: {record.Fats:F0} | C: {record.Carbs:F0}",
                FontSize = 11,
                Foreground = new SolidColorBrush(Color.FromRgb(0x99, 0x99, 0x99))
            });
            Grid.SetColumn(info, 0);
            grid.Children.Add(info);

            if (!string.IsNullOrEmpty(record.MealId))
            {
                var viewBtn = new Button { Content = "👁", Width = 28, Height = 28, Tag = record.MealId };
                viewBtn.Click += (s, e) => MealViewRequested?.Invoke(record.MealId);
                Grid.SetColumn(viewBtn, 1);
                grid.Children.Add(viewBtn);
            }

            var deleteBtn = new Button
            {
                Content = "✕",
                Width = 28,
                Height = 28,
                Foreground = new SolidColorBrush(Color.FromRgb(0xe7, 0x4c, 0x3c))
            };
            deleteBtn.Click += async (s, e) =>
            {
                if (MessageBox.Show("Delete this record?", "Confirm", MessageBoxButton.YesNo) == MessageBoxResult.Yes)
                {
                    await _diaryService.DeleteConsumptionAsync(record.RecordId);
                    LoadDiary();
                }
            };
            Grid.SetColumn(deleteBtn, 2);
            grid.Children.Add(deleteBtn);

            border.Child = grid;
            return border;
        }

        private async Task UpdateProgress()
        {
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            var selectedDate = DayPicker.SelectedDate ?? DateTime.Today;
            var dateStr = selectedDate.ToString("yyyy-MM-dd");

            // Суммируем потребление за выбранный день
            var dayRecords = _groups
                .Where(g => g.Date == dateStr)
                .SelectMany(g => g.Records).ToList();

            double cal = dayRecords.Sum(r => r.Calories);
            double prot = dayRecords.Sum(r => r.Proteins);
            double fat = dayRecords.Sum(r => r.Fats);
            double carb = dayRecords.Sum(r => r.Carbs);
            double water = dayRecords.Sum(r => r.Water);

            // Загружаем нормы из профиля
            var profile = await _profileService.GetProfileAsync(userId);
            double calNorm = profile?.Norms?.Calories ?? 2000;
            double protNorm = profile?.Norms?.Proteins ?? 150;
            double fatNorm = profile?.Norms?.Fats ?? 65;
            double carbNorm = profile?.Norms?.Carbs ?? 300;
            double waterNorm = profile?.Norms?.Water ?? 2500;

            UpdateBar(CalBar, CalText, cal, calNorm, "kcal");
            UpdateBar(ProtBar, ProtText, prot, protNorm, "g");
            UpdateBar(FatBar, FatText, fat, fatNorm, "g");
            UpdateBar(CarbBar, CarbText, carb, carbNorm, "g");
            UpdateBar(WaterBar, WaterText, water, waterNorm, "ml");
        }

        private static void UpdateBar(Border bar, TextBlock text, double value, double max, string unit)
        {
            var pct = max > 0 ? Math.Min(value / max, 1.5) : 0;
            bar.Width = pct * 250;
            text.Text = $"{value:F0} / {max:F0} {unit}";
        }

        private async Task DrawWeightChart()
        {
            WeightCanvas.Children.Clear();
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            // Загружаем историю веса из БД
            var history = await _profileService.GetWeightHistoryAsync(userId);
            var monthData = history
                .Where(e => e.RecordedAt.ToDateTime().Year == _currentMonth.Year &&
                            e.RecordedAt.ToDateTime().Month == _currentMonth.Month)
                .OrderBy(e => e.RecordedAt)
                .Select(e => (Date: e.RecordedAt.ToDateTime(), e.Weight))
                .ToList();

            MonthText.Text = _currentMonth.ToString("MMMM yyyy");

            if (monthData.Count < 2)
            {
                WeightCanvas.Children.Add(new TextBlock
                {
                    Text = "Not enough data",
                    Foreground = new SolidColorBrush(Color.FromRgb(0x99, 0x99, 0x99)),
                    Margin = new Thickness(80, 70, 0, 0)
                });
                return;
            }

            var w = WeightCanvas.ActualWidth > 0 ? WeightCanvas.ActualWidth : 280;
            var h = WeightCanvas.ActualHeight > 0 ? WeightCanvas.ActualHeight : 180;
            var pad = 20;

            var min = monthData.Min(p => p.Weight);
            var max = monthData.Max(p => p.Weight);
            var range = max - min > 0 ? max - min : 1;

            // Оси
            WeightCanvas.Children.Add(new Line { X1 = pad, Y1 = pad, X2 = pad, Y2 = h - pad, Stroke = new SolidColorBrush(Color.FromRgb(0xcc, 0xcc, 0xcc)) });
            WeightCanvas.Children.Add(new Line { X1 = pad, Y1 = h - pad, X2 = w - pad, Y2 = h - pad, Stroke = new SolidColorBrush(Color.FromRgb(0xcc, 0xcc, 0xcc)) });

            // Линия графика
            var polyline = new Polyline { Stroke = new SolidColorBrush(Color.FromRgb(0x66, 0x7E, 0xEA)), StrokeThickness = 2 };
            for (int i = 0; i < monthData.Count; i++)
            {
                var x = pad + (i * (w - pad * 2) / (monthData.Count - 1));
                var y = h - pad - ((monthData[i].Weight - min) / range * (h - pad * 2));
                polyline.Points.Add(new Point(x, y));

                WeightCanvas.Children.Add(new Ellipse
                {
                    Width = 6,
                    Height = 6,
                    Fill = new SolidColorBrush(Color.FromRgb(0x66, 0x7E, 0xEA)),
                    Stroke = Brushes.White,
                    StrokeThickness = 2,
                    Margin = new Thickness(x - 3, y - 3, 0, 0)
                });

                // Подпись значения
                WeightCanvas.Children.Add(new TextBlock
                {
                    Text = monthData[i].Weight.ToString("F1"),
                    FontSize = 9,
                    Foreground = new SolidColorBrush(Color.FromRgb(0x66, 0x66, 0x66)),
                    Margin = new Thickness(x - 15, y - 18, 0, 0)
                });
            }
            WeightCanvas.Children.Add(polyline);
        }

        public event Action<string>? MealViewRequested;
        private void DayPicker_SelectedDateChanged(object sender, SelectionChangedEventArgs e) => _ = UpdateProgress();
        private void PrevMonth_Click(object sender, RoutedEventArgs e) { _currentMonth = _currentMonth.AddMonths(-1); _ = DrawWeightChart(); }
        private void NextMonth_Click(object sender, RoutedEventArgs e) { _currentMonth = _currentMonth.AddMonths(1); _ = DrawWeightChart(); }
        private void SearchBox_KeyDown(object sender, KeyEventArgs e) { if (e.Key == Key.Enter) Search_Click(sender, e); }

        private void Search_Click(object sender, RoutedEventArgs e)
        {
            var query = SearchBox.Text.Trim();
            if (string.IsNullOrEmpty(query)) return;
            var results = _allMeals.Where(m => m.Name.Contains(query, StringComparison.OrdinalIgnoreCase)).ToList();
            if (results.Count == 0)
            {
                MessageBox.Show("No meals found", "Search", MessageBoxButton.OK, MessageBoxImage.Information);
                return;
            }
            ShowMealPicker(results);
        }

        private void AddConsumption_Click(object sender, RoutedEventArgs e) => ShowMealPicker(_allMeals);

        private void ShowMealPicker(List<MealData> meals)
        {
            var dialog = new Window
            {
                Title = "Select Meal",
                Width = 450,
                Height = 500,
                WindowStartupLocation = WindowStartupLocation.CenterOwner,
                Owner = Window.GetWindow(this),
                ResizeMode = ResizeMode.NoResize
            };

            var grid = new Grid { Margin = new Thickness(15) };
            grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });
            grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });
            grid.RowDefinitions.Add(new RowDefinition());
            grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });

            var searchBox = new TextBox { Margin = new Thickness(0, 0, 0, 10), Height = 30 };
            searchBox.TextChanged += (s, ev) =>
            {
                var q = searchBox.Text;
                ((ListBox)grid.Children[2]).ItemsSource = meals
                    .Where(m => m.Name.Contains(q, StringComparison.OrdinalIgnoreCase)).ToList();
            };
            grid.Children.Add(searchBox);

            var amountBox = new TextBox { Text = "100", Height = 30, Margin = new Thickness(0, 0, 0, 10) };
            grid.Children.Add(amountBox);

            var listBox = new ListBox { DisplayMemberPath = "Name" };
            listBox.ItemsSource = meals;
            listBox.SelectionChanged += (s, ev) => { ((Button)grid.Children[3]).IsEnabled = listBox.SelectedItem != null; };
            Grid.SetRow(listBox, 2);
            grid.Children.Add(listBox);

            var addBtn = new Button { Content = "Add", Height = 35, IsEnabled = false };
            addBtn.Click += async (s, ev) =>
            {
                if (listBox.SelectedItem is MealData meal && double.TryParse(amountBox.Text, out var amount))
                {
                    var ratio = amount / 100.0;
                    await _diaryService.AddConsumptionAsync(
                        _sessionService.UserId!, meal.MealId, meal.Name, amount,
                        meal.Calories * ratio, meal.Proteins * ratio, meal.Fats * ratio,
                        meal.Carbs * ratio, meal.Water * ratio);
                    dialog.Close();
                    LoadDiary();
                }
            };
            Grid.SetRow(addBtn, 3);
            grid.Children.Add(addBtn);

            dialog.Content = grid;
            dialog.ShowDialog();
        }
    }
}