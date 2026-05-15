using System.Windows;

namespace EatfitDesktop.Views
{
    public partial class NormsDialog : Window
    {
        public double Calories { get; private set; }
        public double Proteins { get; private set; }
        public double Fats { get; private set; }
        public double Carbs { get; private set; }
        public double Water { get; private set; }

        public NormsDialog()
        {
            InitializeComponent();
        }

        private void Save_Click(object sender, RoutedEventArgs e)
        {
            if (double.TryParse(CaloriesBox.Text, out double cal) &&
                double.TryParse(ProteinsBox.Text, out double prot) &&
                double.TryParse(FatsBox.Text, out double fat) &&
                double.TryParse(CarbsBox.Text, out double carb) &&
                double.TryParse(WaterBox.Text, out double water))
            {
                Calories = cal;
                Proteins = prot;
                Fats = fat;
                Carbs = carb;
                Water = water;
                DialogResult = true;
            }
            else
            {
                MessageBox.Show("Please enter valid numbers", "Error",
                    MessageBoxButton.OK, MessageBoxImage.Warning);
            }
        }

        private void Cancel_Click(object sender, RoutedEventArgs e)
        {
            DialogResult = false;
        }
    }
}