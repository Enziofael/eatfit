package enzio.android.eatfit.ui.screens.profile

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.domain.ProfileRepository
import enzio.android.eatfit.proto.ProfileData
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

data class ProfileState(
    val displayName: String = "Loading...",
    val login: String = "",
    val bio: String = "",
    val height: String = "—",
    val weight: String = "—",
    val age: String = "—",
    val gender: String = "—",
    val birthDate: String = "—",
    val calories: String = "—",
    val proteins: String = "—",
    val fats: String = "—",
    val carbs: String = "—",
    val water: String = "—",
    val isLoading: Boolean = true
)

class ProfileViewModel(
    private val profileRepository: ProfileRepository,
    private val sessionManager: SessionManager
) : ViewModel() {

    private val _state = MutableStateFlow(ProfileState())
    val state: StateFlow<ProfileState> = _state

    init {
        loadProfile()
    }

    fun loadProfile() {
        viewModelScope.launch {
            val userId = sessionManager.userId ?: ""
            if (userId.isEmpty()) return@launch

            val profile = profileRepository.getProfile(userId)
            if (profile != null) {
                _state.value = _state.value.copy(
                    displayName = profile.displayName(),
                    login = profile.login,
                    bio = profile.bio,
                    height = if (profile.height > 0) "${profile.height} cm" else "—",
                    weight = if (profile.currentWeight > 0) "${profile.currentWeight} kg" else "—",
                    age = if (profile.age > 0) profile.age.toString() else "—",
                    gender = when (profile.gender) {
                        "male" -> "Male"
                        "female" -> "Female"
                        "other" -> "Other"
                        else -> "—"
                    },
                    birthDate = profile.birthDate.ifEmpty { "—" },
                    calories = if (profile.hasNorms()) "${profile.norms.calories.toInt()} kcal" else "—",
                    proteins = if (profile.hasNorms()) "${profile.norms.proteins.toInt()} g" else "—",
                    fats = if (profile.hasNorms()) "${profile.norms.fats.toInt()} g" else "—",
                    carbs = if (profile.hasNorms()) "${profile.norms.carbs.toInt()} g" else "—",
                    water = if (profile.hasNorms()) "${profile.norms.water.toInt()} ml" else "—",
                    isLoading = false
                )
            } else {
                _state.value = _state.value.copy(isLoading = false)
            }
        }
    }

    fun updateWeight(weight: Double) {
        viewModelScope.launch {
            val userId = sessionManager.userId ?: return@launch
            if (profileRepository.setWeight(userId, weight)) {
                // Перезагружаем профиль, чтобы вес обновился
                loadProfile()
            }
        }
    }

    fun updateNorms(calories: Double, proteins: Double, fats: Double, carbs: Double, water: Double) {
        viewModelScope.launch {
            val userId = sessionManager.userId ?: return@launch
            if (profileRepository.setNorms(userId, calories, proteins, fats, carbs, water)) {
                loadProfile()
            }
        }
    }

    private fun ProfileData.displayName(): String {
        return if (!firstName.isNullOrEmpty() || !lastName.isNullOrEmpty()) {
            "${firstName ?: ""} ${lastName ?: ""}".trim()
        } else {
            login
        }
    }
}