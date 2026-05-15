package enzio.android.eatfit.ui.screens.main

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.domain.AuthRepository
import kotlinx.coroutines.launch

class MainViewModel(private val repository: AuthRepository) : ViewModel() {

    fun logout(refreshToken: String, onLogout: () -> Unit) {
        viewModelScope.launch {
            try {
                repository.logout(refreshToken)
            } catch (_: Exception) {
                // Игнорируем ошибки логаута
            } finally {
                onLogout()
            }
        }
    }
}