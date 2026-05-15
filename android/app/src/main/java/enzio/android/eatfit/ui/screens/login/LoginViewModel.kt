package enzio.android.eatfit.ui.screens.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.data.model.LoginRequest
import enzio.android.eatfit.domain.AuthRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class LoginViewModel(private val repository: AuthRepository) : ViewModel() {

    private val _state = MutableStateFlow(LoginState())
    val state: StateFlow<LoginState> = _state

    fun onLoginIdentifierChange(value: String) {
        _state.value = _state.value.copy(loginIdentifier = value, error = null)
    }

    fun onPasswordChange(value: String) {
        _state.value = _state.value.copy(password = value, error = null)
    }

    fun login(onSuccess: () -> Unit) {
        val currentState = _state.value
        if (currentState.loginIdentifier.isBlank() || currentState.password.isBlank()) {
            _state.value = currentState.copy(error = "Please fill in all fields")
            return
        }

        viewModelScope.launch {
            _state.value = currentState.copy(isLoading = true)

            val result = repository.login(
                LoginRequest(currentState.loginIdentifier, currentState.password)
            )

            if (result.success) {
                onSuccess()
            } else {
                _state.value = _state.value.copy(
                    isLoading = false,
                    error = result.message
                )
            }
        }
    }
}

data class LoginState(
    val loginIdentifier: String = "",
    val password: String = "",
    val isLoading: Boolean = false,
    val error: String? = null
)