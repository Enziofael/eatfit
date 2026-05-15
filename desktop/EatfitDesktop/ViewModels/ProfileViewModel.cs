using Eatfit.V1;
using EatfitDesktop.Helpers;
using EatfitDesktop.Services;
using System;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Input;

namespace EatfitDesktop.ViewModels
{
    public class ProfileViewModel : INotifyPropertyChanged
    {
        private readonly GrpcProfileService _profileService;
        private readonly SessionService _sessionService;

        public ProfileViewModel(GrpcProfileService profileService, SessionService sessionService)
        {
            _profileService = profileService;
            _sessionService = sessionService;

            UpdateWeightCommand = new RelayCommand(async _ => await UpdateWeightAsync());
            UpdateNormsCommand = new RelayCommand(async _ => await UpdateNormsAsync());
        }

        private string _displayName = "Loading...";
        public string DisplayName
        {
            get => _displayName;
            set { _displayName = value; OnPropertyChanged(); }
        }

        private string _bio = "";
        public string Bio
        {
            get => _bio;
            set { _bio = value; OnPropertyChanged(); }
        }

        private string _height = "—";
        public string Height
        {
            get => _height;
            set { _height = value; OnPropertyChanged(); }
        }

        private string _weight = "—";
        public string Weight
        {
            get => _weight;
            set { _weight = value; OnPropertyChanged(); }
        }

        private string _age = "—";
        public string Age
        {
            get => _age;
            set { _age = value; OnPropertyChanged(); }
        }

        private string _gender = "—";
        public string Gender
        {
            get => _gender;
            set { _gender = value; OnPropertyChanged(); }
        }

        private string _birthDate = "—";
        public string BirthDate
        {
            get => _birthDate;
            set { _birthDate = value; OnPropertyChanged(); }
        }

        // Нормы
        private string _calories = "—";
        public string Calories
        {
            get => _calories;
            set { _calories = value; OnPropertyChanged(); }
        }

        private string _proteins = "—";
        public string Proteins
        {
            get => _proteins;
            set { _proteins = value; OnPropertyChanged(); }
        }

        private string _fats = "—";
        public string Fats
        {
            get => _fats;
            set { _fats = value; OnPropertyChanged(); }
        }

        private string _carbs = "—";
        public string Carbs
        {
            get => _carbs;
            set { _carbs = value; OnPropertyChanged(); }
        }

        private string _water = "—";
        public string Water
        {
            get => _water;
            set { _water = value; OnPropertyChanged(); }
        }

        public ICommand UpdateWeightCommand { get; }
        public ICommand UpdateNormsCommand { get; }

        public async Task LoadProfileAsync()
        {
            var userId = _sessionService.UserId;
            if (string.IsNullOrEmpty(userId)) return;

            var profile = await _profileService.GetProfileAsync(userId);
            if (profile == null) return;

            // Имя
            if (!string.IsNullOrEmpty(profile.FirstName) || !string.IsNullOrEmpty(profile.LastName))
            {
                DisplayName = $"{profile.FirstName} {profile.LastName}".Trim();
            }
            else
            {
                DisplayName = profile.Login;
            }

            Bio = profile.Bio;
            Height = profile.Height > 0 ? $"{profile.Height} cm" : "—";
            Weight = profile.CurrentWeight > 0 ? $"{profile.CurrentWeight} kg" : "—";
            Age = profile.Age > 0 ? profile.Age.ToString() : "—";
            Gender = profile.Gender switch
            {
                "male" => "Male",
                "female" => "Female",
                "other" => "Other",
                _ => "—"
            };
            BirthDate = !string.IsNullOrEmpty(profile.BirthDate) ? profile.BirthDate : "—";

            // Нормы
            if (profile.Norms != null)
            {
                Calories = $"{profile.Norms.Calories:F0} kcal";
                Proteins = $"{profile.Norms.Proteins:F0} g";
                Fats = $"{profile.Norms.Fats:F0} g";
                Carbs = $"{profile.Norms.Carbs:F0} g";
                Water = $"{profile.Norms.Water:F0} ml";
            }
        }

        private async Task UpdateWeightAsync()
        {
            var input = Microsoft.VisualBasic.Interaction.InputBox(
                "Enter new weight (kg):", "Update Weight", "70");

            if (double.TryParse(input, out double weight) && weight > 0)
            {
                var userId = _sessionService.UserId;
                if (string.IsNullOrEmpty(userId)) return;

                var success = await _profileService.UpdateWeightAsync(userId, weight);
                if (success)
                {
                    Weight = $"{weight} kg";
                    MessageBox.Show("Weight updated!", "Success", MessageBoxButton.OK, MessageBoxImage.Information);
                }
                else
                {
                    MessageBox.Show("Failed to update weight", "Error", MessageBoxButton.OK, MessageBoxImage.Error);
                }
            }
        }

        private async Task UpdateNormsAsync()
        {
            // Создаём простое диалоговое окно для ввода норм
            var dialog = new Views.NormsDialog();
            if (dialog.ShowDialog() == true)
            {
                var userId = _sessionService.UserId;
                if (string.IsNullOrEmpty(userId)) return;

                var success = await _profileService.UpdateNormsAsync(
                    userId,
                    dialog.Calories,
                    dialog.Proteins,
                    dialog.Fats,
                    dialog.Carbs,
                    dialog.Water
                );

                if (success)
                {
                    await LoadProfileAsync();
                    MessageBox.Show("Norms updated!", "Success", MessageBoxButton.OK, MessageBoxImage.Information);
                }
            }
        }

        public event PropertyChangedEventHandler? PropertyChanged;
        protected void OnPropertyChanged([CallerMemberName] string? propertyName = null)
            => PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
    }
}