package enzio.android.eatfit.data.remote

import enzio.android.eatfit.proto.*
import io.grpc.ManagedChannel
import io.grpc.okhttp.OkHttpChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.util.concurrent.TimeUnit

class GrpcDiaryService(
    private val serverUrl: String = "192.168.0.10",
    private val port: Int = 50051
) {
    private val channel: ManagedChannel = OkHttpChannelBuilder
        .forAddress(serverUrl, port)
        .usePlaintext()
        .keepAliveTime(30, TimeUnit.SECONDS)
        .build()

    private val stub = DiaryServiceGrpc.newBlockingStub(channel)

    suspend fun listConsumptions(userId: String, days: Int = 30): List<ConsumptionGroup> = withContext(Dispatchers.IO) {
        try {
            val request = ListConsumptionsRequest.newBuilder()
                .setUserId(userId)
                .setLimit(days)
                .build()
            stub.listConsumptions(request).groupsList
        } catch (e: Exception) {
            emptyList()
        }
    }

    suspend fun addConsumption(
        userId: String, mealId: String, mealName: String, amount: Double,
        calories: Double, proteins: Double, fats: Double, carbs: Double, water: Double
    ): Boolean = withContext(Dispatchers.IO) {
        try {
            val request = AddConsumptionRequest.newBuilder()
                .setUserId(userId)
                .setMealId(mealId)
                .setMealName(mealName)
                .setAmount(amount)
                .setCalories(calories)
                .setProteins(proteins)
                .setFats(fats)
                .setCarbs(carbs)
                .setWater(water)
                .build()
            stub.addConsumption(request).success
        } catch (e: Exception) {
            false
        }
    }

    suspend fun deleteConsumption(recordId: String): Boolean = withContext(Dispatchers.IO) {
        try {
            stub.deleteConsumption(DeleteConsumptionRequest.newBuilder().setRecordId(recordId).build()).success
        } catch (e: Exception) {
            false
        }
    }

    fun close() { channel.shutdown() }
}