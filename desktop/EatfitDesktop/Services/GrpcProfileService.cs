using System;
using System.Net.Http;
using System.Threading.Tasks;
using Eatfit.V1;
using Grpc.Core;
using Grpc.Net.Client;

namespace EatfitDesktop.Services
{
    public class GrpcProfileService : IDisposable
    {
        private readonly GrpcChannel _channel;
        private readonly Eatfit.V1.ProfileService.ProfileServiceClient _client;

        public GrpcProfileService(string serverUrl = "http://localhost:50051")
        {
            var handler = new HttpClientHandler
            {
                ServerCertificateCustomValidationCallback =
                    HttpClientHandler.DangerousAcceptAnyServerCertificateValidator
            };

            _channel = GrpcChannel.ForAddress(serverUrl, new GrpcChannelOptions
            {
                HttpHandler = handler
            });

            _client = new Eatfit.V1.ProfileService.ProfileServiceClient(_channel);
        }

        public async Task<ProfileData?> GetProfileAsync(string userId)
        {
            try
            {
                var request = new GetProfileRequest { UserId = userId };
                var response = await _client.GetProfileAsync(request);
                return response.Profile;
            }
            catch (RpcException ex)
            {
                System.Diagnostics.Debug.WriteLine($"gRPC Error: {ex.Status.Detail}");
                return null;
            }
        }

        public async Task<bool> UpdateWeightAsync(string userId, double weight, string note = "")
        {
            try
            {
                var request = new SetWeightRequest
                {
                    UserId = userId,
                    Weight = weight,
                    Note = note
                };
                var response = await _client.SetWeightAsync(request);
                return response.Success;
            }
            catch
            {
                return false;
            }
        }

        public async Task<bool> UpdateNormsAsync(string userId, double calories, double proteins, double fats, double carbs, double water)
        {
            try
            {
                var request = new SetNormsRequest
                {
                    UserId = userId,
                    Calories = calories,
                    Proteins = proteins,
                    Fats = fats,
                    Carbs = carbs,
                    Water = water
                };
                var response = await _client.SetNormsAsync(request);
                return response.Success;
            }
            catch
            {
                return false;
            }
        }

        public void Dispose()
        {
            _channel?.Dispose();
        }
    }
}