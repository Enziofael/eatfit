// ui/screens/reset/ResetPasswordViewModel.kt
package enzio.android.eatfit.ui.screens.reset

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.domain.AuthRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class ResetPasswordViewModel(private val repository: AuthRepository) : ViewModel() {
    private val _state = MutableStateFlow(ResetPasswordState())
    val state: StateFlow<ResetPasswordState> = _state

    fun onCodeChange(value: String) { _state.value = _state.value.copy(code = value) }
    fun onPasswordChange(value: String) { _state.value = _state.value.copy(newPassword = value) }
    fun onConfirmChange(value: String) { _state.value = _state.value.copy(confirmPassword = value) }

    fun resetPassword(resetToken: String, onSuccess: () -> Unit) {
        val current = _state.value
        if (current.newPassword != current.confirmPassword) {
            _state.value = current.copy(error = "Passwords do not match")
            return
        }

        viewModelScope.launch {
            _state.value = current.copy(isLoading = true)
            val result = repository.resetPassword(resetToken, current.code, current.newPassword)
            if (result.success) onSuccess()
            else _state.value = _state.value.copy(isLoading = false, error = result.message)
        }
    }
}

data class ResetPasswordState(
    val code: String = "",
    val newPassword: String = "",
    val confirmPassword: String = "",
    val isLoading: Boolean = false,
    val error: String? = null
)
