package enzio.android.eatfit.ui.screens.register

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.data.model.RegisterRequest
import enzio.android.eatfit.domain.AuthRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class RegisterViewModel(private val repository: AuthRepository) : ViewModel() {

    private val _state = MutableStateFlow(RegisterState())
    val state: StateFlow<RegisterState> = _state

    fun onEmailChange(value: String) { _state.value = _state.value.copy(email = value, error = null) }
    fun onLoginChange(value: String) { _state.value = _state.value.copy(login = value, error = null) }
    fun onPasswordChange(value: String) { _state.value = _state.value.copy(password = value, error = null) }
    fun onConfirmPasswordChange(value: String) { _state.value = _state.value.copy(confirmPassword = value, error = null) }

    fun register(onSuccess: (String, String, String) -> Unit) {
        val current = _state.value

        if (current.email.isBlank() || current.login.isBlank() || current.password.isBlank()) {
            _state.value = current.copy(error = "Please fill in all fields")
            return
        }
        if (current.password != current.confirmPassword) {
            _state.value = current.copy(error = "Passwords do not match")
            return
        }

        viewModelScope.launch {
            _state.value = current.copy(isLoading = true)

            val result = repository.register(
                RegisterRequest(current.email, current.login, current.password, current.confirmPassword)
            )

            if (result.success) {
                onSuccess(result.userId, current.email, current.login)
            } else {
                _state.value = _state.value.copy(isLoading = false, error = result.message)
            }
        }
    }
}

data class RegisterState(
    val email: String = "",
    val login: String = "",
    val password: String = "",
    val confirmPassword: String = "",
    val isLoading: Boolean = false,
    val error: String? = null
)