package enzio.android.eatfit.ui.screens.meals

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.data.remote.GrpcMealService
import enzio.android.eatfit.proto.MealComponentInput
import enzio.android.eatfit.proto.MealData
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

data class MealsState(
    val meals: List<MealData> = emptyList(),
    val searchQuery: String = "",
    val sortBy: String = "created_at",
    val sortOrder: String = "desc",
    val isLoading: Boolean = false
)

class MealsViewModel(
    private val mealService: GrpcMealService,
    private val sessionManager: SessionManager
) : ViewModel() {

    private val _state = MutableStateFlow(MealsState())
    val state: StateFlow<MealsState> = _state

    init {
        loadMeals()
    }

    fun loadMeals() {
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true)
            val userId = sessionManager.userId ?: return@launch
            val meals = mealService.listMeals(userId, _state.value.sortBy, _state.value.sortOrder)
            _state.value = _state.value.copy(meals = meals, isLoading = false)
        }
    }

    fun search(query: String) {
        _state.value = _state.value.copy(searchQuery = query)
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true)
            val userId = sessionManager.userId ?: return@launch
            val meals = if (query.isBlank()) {
                mealService.listMeals(userId, _state.value.sortBy, _state.value.sortOrder)
            } else {
                mealService.searchMeals(userId, query, _state.value.sortBy, _state.value.sortOrder)
            }
            _state.value = _state.value.copy(meals = meals, isLoading = false)
        }
    }

    fun setSortBy(sortBy: String) {
        _state.value = _state.value.copy(sortBy = sortBy)
        loadMeals()
    }

    fun deleteMeal(mealId: String) {
        viewModelScope.launch {
            mealService.deleteMeal(mealId)
            loadMeals()
        }
    }

    fun onSearchQueryChange(query: String) {
        _state.value = _state.value.copy(searchQuery = query)
    }
}