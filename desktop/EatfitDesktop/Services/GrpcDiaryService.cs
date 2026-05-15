using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using Eatfit.V1;
using Grpc.Net.Client;

namespace EatfitDesktop.Services
{
    public class GrpcDiaryService : IDisposable
    {
        private readonly GrpcChannel _channel;
        private readonly Eatfit.V1.DiaryService.DiaryServiceClient _client;

        public GrpcDiaryService(string serverUrl = "http://localhost:50051")
        {
            var handler = new HttpClientHandler
            {
                ServerCertificateCustomValidationCallback = HttpClientHandler.DangerousAcceptAnyServerCertificateValidator
            };
            _channel = GrpcChannel.ForAddress(serverUrl, new GrpcChannelOptions { HttpHandler = handler });
            _client = new Eatfit.V1.DiaryService.DiaryServiceClient(_channel);
        }

        public async Task<List<ConsumptionGroup>> ListConsumptionsAsync(string userId, int days = 30)
        {
            try
            {
                var response = await _client.ListConsumptionsAsync(new ListConsumptionsRequest
                {
                    UserId = userId,
                    Limit = days
                });
                return new List<ConsumptionGroup>(response.Groups);
            }
            catch { return new(); }
        }

        public async Task AddConsumptionAsync(string userId, string mealId, string mealName,
            double amount, double calories, double proteins, double fats, double carbs, double water)
        {
            try
            {
                await _client.AddConsumptionAsync(new AddConsumptionRequest
                {
                    UserId = userId,
                    MealId = mealId,
                    MealName = mealName,
                    Amount = amount,
                    Calories = calories,
                    Proteins = proteins,
                    Fats = fats,
                    Carbs = carbs,
                    Water = water
                });
            }
            catch { }
        }

        public async Task DeleteConsumptionAsync(string recordId)
        {
            try
            {
                await _client.DeleteConsumptionAsync(new DeleteConsumptionRequest { RecordId = recordId });
            }
            catch { }
        }

        public void Dispose() => _channel?.Dispose();
    }
}