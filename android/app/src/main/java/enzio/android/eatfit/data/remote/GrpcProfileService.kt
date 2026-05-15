package enzio.android.eatfit.data.remote

import enzio.android.eatfit.proto.GetProfileRequest
import enzio.android.eatfit.proto.ProfileData
import enzio.android.eatfit.proto.ProfileServiceGrpc
import enzio.android.eatfit.proto.SetNormsRequest
import enzio.android.eatfit.proto.SetWeightRequest
import io.grpc.ManagedChannel
import io.grpc.okhttp.OkHttpChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.util.concurrent.TimeUnit

class GrpcProfileService(
    private val serverUrl: String = "192.168.0.10",
    private val port: Int = 50051
) {
    private val channel: ManagedChannel = OkHttpChannelBuilder
        .forAddress(serverUrl, port)
        .usePlaintext()
        .keepAliveTime(30, TimeUnit.SECONDS)
        .build()

    private val stub = ProfileServiceGrpc.newBlockingStub(channel)

    suspend fun getProfile(userId: String): ProfileData? = withContext(Dispatchers.IO) {
        try {
            val request = GetProfileRequest.newBuilder().setUserId(userId).build()
            val response = stub.getProfile(request)
            response.profile
        } catch (e: Exception) {
            null
        }
    }

    suspend fun setWeight(userId: String, weight: Double, note: String = ""): Boolean = withContext(Dispatchers.IO) {
        try {
            val request = SetWeightRequest.newBuilder()
                .setUserId(userId)
                .setWeight(weight)
                .setNote(note)
                .build()
            stub.setWeight(request).success
        } catch (e: Exception) {
            false
        }
    }

    suspend fun setNorms(userId: String, calories: Double, proteins: Double, fats: Double, carbs: Double, water: Double): Boolean = withContext(Dispatchers.IO) {
        try {
            val request = SetNormsRequest.newBuilder()
                .setUserId(userId)
                .setCalories(calories)
                .setProteins(proteins)
                .setFats(fats)
                .setCarbs(carbs)
                .setWater(water)
                .build()
            stub.setNorms(request).success
        } catch (e: Exception) {
            false
        }
    }

    fun close() {
        channel.shutdown()
    }
}