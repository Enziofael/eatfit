package enzio.android.eatfit.ui.screens.diary

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.data.remote.GrpcDiaryService
import enzio.android.eatfit.data.remote.GrpcMealService
import enzio.android.eatfit.data.remote.GrpcProfileService
import enzio.android.eatfit.proto.ConsumptionGroup
import enzio.android.eatfit.proto.MealData
import enzio.android.eatfit.proto.ProfileData
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.YearMonth
import java.time.ZoneId
import java.time.format.DateTimeFormatter

data class DiaryState(
    val groups: List<ConsumptionGroup> = emptyList(),
    val searchQuery: String = "",
    val searchResults: List<MealData> = emptyList(),
    val selectedDate: LocalDate = LocalDate.now(),
    val currentMonth: YearMonth = YearMonth.now(),
    val profile: ProfileData? = null,
    val weightHistory: List<WeightPoint> = emptyList(),
    val isLoading: Boolean = false,
    val showMealPicker: Boolean = false,
    val selectedMeal: MealData? = null,
    val amount: String = "100"
)

data class WeightPoint(val date: LocalDate, val weight: Double)

class DiaryViewModel(
    private val diaryService: GrpcDiaryService,
    private val mealService: GrpcMealService,
    private val profileService: GrpcProfileService,
    private val sessionManager: SessionManager
) : ViewModel() {

    private val _state = MutableStateFlow(DiaryState())
    val state: StateFlow<DiaryState> = _state

    init { loadData() }

    fun loadData() {
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true)
            val userId = sessionManager.userId ?: return@launch

            val groups = diaryService.listConsumptions(userId)
            val profile = profileService.getProfile(userId)
            val history = profileService.getWeightHistory(userId)

            _state.value = _state.value.copy(
                groups = groups,
                profile = profile,
                weightHistory = history.map { entry ->
                    val instant = Instant.ofEpochSecond(entry.recordedAt.seconds, entry.recordedAt.nanos.toLong())
                    val date = LocalDateTime.ofInstant(instant, ZoneId.systemDefault()).toLocalDate()
                    WeightPoint(date, entry.weight)
                },
                isLoading = false
            )
        }
    }

    fun onSearchQueryChange(query: String) {
        _state.value = _state.value.copy(searchQuery = query)
    }

    fun searchMeals() {
        viewModelScope.launch {
            val userId = sessionManager.userId ?: return@launch
            val query = _state.value.searchQuery
            val meals = if (query.isBlank()) {
                mealService.listMeals(userId)
            } else {
                mealService.searchMeals(userId, query)
            }
            _state.value = _state.value.copy(searchResults = meals, showMealPicker = true)
        }
    }

    fun showAllMeals() {
        viewModelScope.launch {
            val userId = sessionManager.userId ?: return@launch
            val meals = mealService.listMeals(userId)
            _state.value = _state.value.copy(searchResults = meals, showMealPicker = true)
        }
    }

    fun hideMealPicker() {
        _state.value = _state.value.copy(showMealPicker = false, selectedMeal = null)
    }

    fun selectMeal(meal: MealData) {
        _state.value = _state.value.copy(selectedMeal = meal)
    }

    fun onAmountChange(amount: String) {
        _state.value = _state.value.copy(amount = amount)
    }

    fun addConsumption() {
        val meal = _state.value.selectedMeal ?: return
        val amount = _state.value.amount.toDoubleOrNull() ?: 100.0

        viewModelScope.launch {
            val userId = sessionManager.userId ?: return@launch
            val ratio = amount / 100.0
            diaryService.addConsumption(
                userId, meal.mealId, meal.name, amount,
                meal.calories * ratio, meal.proteins * ratio,
                meal.fats * ratio, meal.carbs * ratio, meal.water * ratio
            )
            hideMealPicker()
            loadData()
        }
    }

    fun deleteConsumption(recordId: String) {
        viewModelScope.launch {
            diaryService.deleteConsumption(recordId)
            loadData()
        }
    }

    fun setSelectedDate(date: LocalDate) {
        _state.value = _state.value.copy(selectedDate = date)
    }

    fun prevMonth() {
        _state.value = _state.value.copy(currentMonth = _state.value.currentMonth.minusMonths(1))
    }

    fun nextMonth() {
        _state.value = _state.value.copy(currentMonth = _state.value.currentMonth.plusMonths(1))
    }
}