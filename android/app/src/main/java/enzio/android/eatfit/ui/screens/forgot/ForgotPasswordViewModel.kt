// ui/screens/forgot/ForgotPasswordViewModel.kt
package enzio.android.eatfit.ui.screens.forgot

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.domain.AuthRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class ForgotPasswordViewModel(private val repository: AuthRepository) : ViewModel() {
    private val _state = MutableStateFlow(ForgotPasswordState())
    val state: StateFlow<ForgotPasswordState> = _state

    fun onIdentifierChange(value: String) { _state.value = _state.value.copy(loginIdentifier = value, error = null) }

    fun sendResetCode(onSuccess: (String) -> Unit) {
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true)
            val result = repository.forgotPassword(_state.value.loginIdentifier)
            if (result.success && result.resetToken.isNotEmpty()) {
                onSuccess(result.resetToken)
            } else {
                _state.value = _state.value.copy(isLoading = false, message = "If account exists, reset code sent")
            }
        }
    }
}

data class ForgotPasswordState(
    val loginIdentifier: String = "",
    val isLoading: Boolean = false,
    val error: String? = null,
    val message: String? = null
)