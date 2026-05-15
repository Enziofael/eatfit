package enzio.android.eatfit.data.remote

import enzio.android.eatfit.proto.*
import io.grpc.ManagedChannel
import io.grpc.okhttp.OkHttpChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.util.concurrent.TimeUnit

class GrpcMealService(
    private val serverUrl: String = "192.168.0.10",
    private val port: Int = 50051
) {
    private val channel: ManagedChannel = OkHttpChannelBuilder
        .forAddress(serverUrl, port)
        .usePlaintext()
        .keepAliveTime(30, TimeUnit.SECONDS)
        .build()

    private val stub = MealServiceGrpc.newBlockingStub(channel)

    suspend fun listMeals(userId: String, sortBy: String = "created_at", sortOrder: String = "desc"): List<MealData> = withContext(Dispatchers.IO) {
        try {
            val request = ListMealsRequest.newBuilder()
                .setUserId(userId)
                .setSortBy(sortBy)
                .setSortOrder(sortOrder)
                .setPage(1)
                .setPageSize(100)
                .build()
            stub.listMeals(request).mealsList
        } catch (e: Exception) {
            emptyList()
        }
    }

    suspend fun searchMeals(userId: String, query: String, sortBy: String = "created_at", sortOrder: String = "desc"): List<MealData> = withContext(Dispatchers.IO) {
        try {
            val request = SearchMealsRequest.newBuilder()
                .setUserId(userId)
                .setQuery(query)
                .setSortBy(sortBy)
                .setSortOrder(sortOrder)
                .setPage(1)
                .setPageSize(100)
                .build()
            stub.searchMeals(request).mealsList
        } catch (e: Exception) {
            emptyList()
        }
    }

    suspend fun getMeal(mealId: String): MealData? = withContext(Dispatchers.IO) {
        try {
            val request = GetMealRequest.newBuilder().setMealId(mealId).build()
            stub.getMeal(request).meal
        } catch (e: Exception) {
            null
        }
    }

    suspend fun createMeal(
        userId: String, name: String, description: String, recipe: String,
        imageUrl: String, calories: Double, proteins: Double, fats: Double,
        carbs: Double, water: Double, components: List<MealComponentInput>
    ): Boolean = withContext(Dispatchers.IO) {
        try {
            val request = CreateMealRequest.newBuilder()
                .setUserId(userId)
                .setName(name)
                .setDescription(description)
                .setRecipe(recipe)
                .setImageUrl(imageUrl)
                .setCalories(calories)
                .setProteins(proteins)
                .setFats(fats)
                .setCarbs(carbs)
                .setWater(water)
                .addAllComponents(components)
                .build()
            stub.createMeal(request).success
        } catch (e: Exception) {
            false
        }
    }

    suspend fun updateMeal(
        mealId: String, name: String, description: String, recipe: String,
        imageUrl: String, calories: Double, proteins: Double, fats: Double,
        carbs: Double, water: Double, components: List<MealComponentInput>
    ): Boolean = withContext(Dispatchers.IO) {
        try {
            val request = UpdateMealRequest.newBuilder()
                .setMealId(mealId)
                .setName(name)
                .setDescription(description)
                .setRecipe(recipe)
                .setImageUrl(imageUrl)
                .setCalories(calories)
                .setProteins(proteins)
                .setFats(fats)
                .setCarbs(carbs)
                .setWater(water)
                .addAllComponents(components)
                .build()
            stub.updateMeal(request).success
        } catch (e: Exception) {
            false
        }
    }

    suspend fun deleteMeal(mealId: String): Boolean = withContext(Dispatchers.IO) {
        try {
            stub.deleteMeal(DeleteMealRequest.newBuilder().setMealId(mealId).build()).success
        } catch (e: Exception) {
            false
        }
    }

    fun close() {
        channel.shutdown()
    }
}