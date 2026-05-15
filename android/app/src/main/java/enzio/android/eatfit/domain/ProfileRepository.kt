package enzio.android.eatfit.domain

import enzio.android.eatfit.data.remote.GrpcProfileService
import enzio.android.eatfit.proto.ProfileData

class ProfileRepository(private val grpcProfileService: GrpcProfileService) {
    suspend fun getProfile(userId: String): ProfileData? {
        return grpcProfileService.getProfile(userId)
    }

    suspend fun setWeight(userId: String, weight: Double): Boolean {
        return grpcProfileService.setWeight(userId, weight)
    }

    suspend fun setNorms(userId: String, calories: Double, proteins: Double, fats: Double, carbs: Double, water: Double): Boolean {
        return grpcProfileService.setNorms(userId, calories, proteins, fats, carbs, water)
    }
}