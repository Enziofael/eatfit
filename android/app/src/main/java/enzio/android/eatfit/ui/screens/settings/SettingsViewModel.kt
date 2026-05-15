package enzio.android.eatfit.ui.screens.settings

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.data.remote.GrpcAuthService
import enzio.android.eatfit.data.remote.GrpcProfileService
import enzio.android.eatfit.proto.ProfileData
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

data class SettingsState(
    // Профиль
    val firstName: String = "",
    val lastName: String = "",
    val height: String = "",
    val birthDate: String = "",
    val gender: String = "",
    // Логин
    val newLogin: String = "",
    val loginPassword: String = "",
    // Пароль
    val currentPassword: String = "",
    val newPassword: String = "",
    val confirmPassword: String = "",
    // Сообщения
    val profileMessage: String? = null,
    val loginMessage: String? = null,
    val passwordMessage: String? = null,
    val isLoading: Boolean = false
)

class SettingsViewModel(
    private val profileService: GrpcProfileService,
    private val authService: GrpcAuthService,
    private val sessionManager: SessionManager
) : ViewModel() {

    private val _state = MutableStateFlow(SettingsState())
    val state: StateFlow<SettingsState> = _state

    init {
        loadProfile()
    }

    fun loadProfile() {
        viewModelScope.launch {
            val userId = sessionManager.userId ?: return@launch
            val profile = profileService.getProfile(userId) ?: return@launch

            _state.value = _state.value.copy(
                firstName = profile.firstName,
                lastName = profile.lastName,
                height = if (profile.height > 0) profile.height.toInt().toString() else "",
                birthDate = profile.birthDate,
                gender = profile.gender
            )
        }
    }

    fun onFirstNameChange(value: String) { _state.value = _state.value.copy(firstName = value) }
    fun onLastNameChange(value: String) { _state.value = _state.value.copy(lastName = value) }
    fun onHeightChange(value: String) { _state.value = _state.value.copy(height = value) }
    fun onBirthDateChange(value: String) { _state.value = _state.value.copy(birthDate = value) }
    fun onGenderChange(value: String) { _state.value = _state.value.copy(gender = value) }
    fun onNewLoginChange(value: String) { _state.value = _state.value.copy(newLogin = value) }
    fun onLoginPasswordChange(value: String) { _state.value = _state.value.copy(loginPassword = value) }
    fun onCurrentPasswordChange(value: String) { _state.value = _state.value.copy(currentPassword = value) }
    fun onNewPasswordChange(value: String) { _state.value = _state.value.copy(newPassword = value) }
    fun onConfirmPasswordChange(value: String) { _state.value = _state.value.copy(confirmPassword = value) }

    fun saveProfile() {
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true, profileMessage = null)
            val userId = sessionManager.userId ?: return@launch

            val height = _state.value.height.toDoubleOrNull() ?: 0.0
            val success = profileService.updateProfile(
                userId, _state.value.firstName, _state.value.lastName,
                height, _state.value.birthDate, _state.value.gender
            )

            _state.value = _state.value.copy(
                isLoading = false,
                profileMessage = if (success) "Profile saved!" else "Failed to save profile"
            )
        }
    }

    fun changeLogin() {
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true, loginMessage = null)
            val userId = sessionManager.userId ?: return@launch

            val result = authService.changeLogin(
                userId, _state.value.newLogin, _state.value.loginPassword
            )

            _state.value = _state.value.copy(
                isLoading = false,
                loginMessage = if (result.success) "Login changed!" else result.message,
                newLogin = if (result.success) "" else _state.value.newLogin,
                loginPassword = ""
            )
        }
    }

    fun changePassword() {
        viewModelScope.launch {
            val s = _state.value
            if (s.newPassword != s.confirmPassword) {
                _state.value = _state.value.copy(passwordMessage = "Passwords do not match")
                return@launch
            }

            _state.value = _state.value.copy(isLoading = true, passwordMessage = null)
            val userId = sessionManager.userId ?: return@launch

            val result = authService.changePassword(
                userId, s.currentPassword, s.newPassword, s.confirmPassword
            )

            _state.value = _state.value.copy(
                isLoading = false,
                passwordMessage = if (result.success) "Password changed!" else result.message,
                currentPassword = "",
                newPassword = "",
                confirmPassword = ""
            )
        }
    }
}