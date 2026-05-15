// ui/screens/verify/VerifyEmailViewModel.kt
package enzio.android.eatfit.ui.screens.verify

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.domain.AuthRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class VerifyEmailViewModel(private val repository: AuthRepository) : ViewModel() {
    private val _state = MutableStateFlow(VerifyState())
    val state: StateFlow<VerifyState> = _state

    fun onCodeChange(value: String) { _state.value = _state.value.copy(code = value, error = null) }

    fun verify(userId: String, onSuccess: () -> Unit) {
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true)
            val result = repository.verifyEmail(userId, _state.value.code)
            if (result.success) onSuccess()
            else _state.value = _state.value.copy(isLoading = false, error = result.message)
        }
    }
}

data class VerifyState(val code: String = "", val isLoading: Boolean = false, val error: String? = null)