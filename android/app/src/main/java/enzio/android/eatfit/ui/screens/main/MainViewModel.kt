// ui/screens/main/MainViewModel.kt
package enzio.android.eatfit.ui.screens.main

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.domain.AuthRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.launch

class MainViewModel(private val repository: AuthRepository) : ViewModel() {
    fun logout(refreshToken: String, onLogout: () -> Unit) {
        viewModelScope.launch {
            repository.logout(refreshToken)
            onLogout()
        }
    }
}