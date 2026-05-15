using System;
using System.IO;
using System.Text.Json;

namespace EatfitDesktop.Services
{
    public class SessionService
    {
        private readonly string _sessionFilePath;

        public string? AccessToken { get; private set; }
        public string? RefreshToken { get; private set; }
        public string? UserId { get; private set; }
        public string? Email { get; private set; }
        public string? Login { get; private set; }
        public bool IsLoggedIn => !string.IsNullOrEmpty(AccessToken);

        public SessionService()
        {
            var appData = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
            var eatfitFolder = Path.Combine(appData, "Eatfit");
            Directory.CreateDirectory(eatfitFolder);
            _sessionFilePath = Path.Combine(eatfitFolder, "session.json");
            LoadSession();
        }

        public void SaveSession(string accessToken, string refreshToken, string userId, string email, string login)
        {
            AccessToken = accessToken;
            RefreshToken = refreshToken;
            UserId = userId;
            Email = email;
            Login = login;

            var data = new
            {
                AccessToken = accessToken,
                RefreshToken = refreshToken,
                UserId = userId,
                Email = email,
                Login = login
            };

            var json = JsonSerializer.Serialize(data);
            File.WriteAllText(_sessionFilePath, json);
        }

        public void ClearSession()
        {
            AccessToken = null;
            RefreshToken = null;
            UserId = null;
            Email = null;
            Login = null;

            if (File.Exists(_sessionFilePath))
                File.Delete(_sessionFilePath);
        }

        private void LoadSession()
        {
            if (!File.Exists(_sessionFilePath)) return;

            try
            {
                var json = File.ReadAllText(_sessionFilePath);
                var data = JsonSerializer.Deserialize<SessionData>(json);

                if (data != null)
                {
                    AccessToken = data.AccessToken;
                    RefreshToken = data.RefreshToken;
                    UserId = data.UserId;
                    Email = data.Email;
                    Login = data.Login;
                }
            }
            catch
            {
                ClearSession();
            }
        }

        private class SessionData
        {
            public string AccessToken { get; set; } = string.Empty;
            public string RefreshToken { get; set; } = string.Empty;
            public string UserId { get; set; } = string.Empty;
            public string Email { get; set; } = string.Empty;
            public string Login { get; set; } = string.Empty;
        }
    }
}