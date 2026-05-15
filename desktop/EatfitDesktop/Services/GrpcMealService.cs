using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using Eatfit.V1;
using Grpc.Core;
using Grpc.Net.Client;

namespace EatfitDesktop.Services
{
    public class GrpcMealService : IDisposable
    {
        private readonly GrpcChannel _channel;
        private readonly Eatfit.V1.MealService.MealServiceClient _client;
        private bool _disposed;

        public GrpcMealService(string serverUrl = "http://localhost:50051")
        {
            var handler = new HttpClientHandler
            {
                ServerCertificateCustomValidationCallback =
                    HttpClientHandler.DangerousAcceptAnyServerCertificateValidator
            };

            _channel = GrpcChannel.ForAddress(serverUrl, new GrpcChannelOptions { HttpHandler = handler });
            _client = new Eatfit.V1.MealService.MealServiceClient(_channel);
        }

        public async Task<List<MealData>> ListMealsAsync(string userId, string sortBy = "created_at", string sortOrder = "desc")
        {
            try
            {
                var request = new ListMealsRequest
                {
                    UserId = userId,
                    SortBy = sortBy,
                    SortOrder = sortOrder,
                    Page = 1,
                    PageSize = 100
                };
                var response = await _client.ListMealsAsync(request);
                return new List<MealData>(response.Meals);
            }
            catch (RpcException)
            {
                return new List<MealData>();
            }
        }

        public async Task<List<MealData>> SearchMealsAsync(string userId, string query, string sortBy = "created_at", string sortOrder = "desc")
        {
            try
            {
                var request = new SearchMealsRequest
                {
                    UserId = userId,
                    Query = query,
                    SortBy = sortBy,
                    SortOrder = sortOrder,
                    Page = 1,
                    PageSize = 100
                };
                var response = await _client.SearchMealsAsync(request);
                return new List<MealData>(response.Meals);
            }
            catch
            {
                return new List<MealData>();
            }
        }

        public async Task<MealData?> GetMealAsync(string mealId)
        {
            try
            {
                var response = await _client.GetMealAsync(new GetMealRequest { MealId = mealId });
                return response.Meal;
            }
            catch
            {
                return null;
            }
        }

        public async Task<bool> CreateMealAsync(string userId, string name, string description, string recipe,
            string imageUrl, double calories, double proteins, double fats, double carbs, double water,
            List<MealComponentInput> components)
        {
            try
            {
                var request = new CreateMealRequest
                {
                    UserId = userId,
                    Name = name,
                    Description = description,
                    Recipe = recipe,
                    ImageUrl = imageUrl,
                    Calories = calories,
                    Proteins = proteins,
                    Fats = fats,
                    Carbs = carbs,
                    Water = water
                };
                request.Components.AddRange(components);
                await _client.CreateMealAsync(request);
                return true;
            }
            catch
            {
                return false;
            }
        }

        public async Task<bool> UpdateMealAsync(string mealId, string name, string description, string recipe,
            string imageUrl, double calories, double proteins, double fats, double carbs, double water,
            List<MealComponentInput> components)
        {
            try
            {
                var request = new UpdateMealRequest
                {
                    MealId = mealId,
                    Name = name,
                    Description = description,
                    Recipe = recipe,
                    ImageUrl = imageUrl,
                    Calories = calories,
                    Proteins = proteins,
                    Fats = fats,
                    Carbs = carbs,
                    Water = water
                };
                request.Components.AddRange(components);
                await _client.UpdateMealAsync(request);
                return true;
            }
            catch
            {
                return false;
            }
        }

        public async Task<bool> DeleteMealAsync(string mealId)
        {
            try
            {
                await _client.DeleteMealAsync(new DeleteMealRequest { MealId = mealId });
                return true;
            }
            catch
            {
                return false;
            }
        }

        public void Dispose()
        {
            if (_disposed) return;
            _disposed = true;
            _channel?.Dispose();
            GC.SuppressFinalize(this);
        }
    }
}