package enzio.android.eatfit.ui.screens.meals

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.data.remote.GrpcMealService
import enzio.android.eatfit.proto.MealComponentInput
import enzio.android.eatfit.proto.MealData
import enzio.android.eatfit.ui.theme.*
import kotlinx.coroutines.launch
import androidx.compose.foundation.clickable
import androidx.compose.material3.HorizontalDivider

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MealEditorScreen(
    mealService: GrpcMealService,
    sessionManager: SessionManager,
    mealId: String? = null,
    onSaved: () -> Unit,
    onCancel: () -> Unit
) {
    var name by remember { mutableStateOf("") }
    var description by remember { mutableStateOf("") }
    var recipe by remember { mutableStateOf("") }
    var calories by remember { mutableStateOf("") }
    var proteins by remember { mutableStateOf("") }
    var fats by remember { mutableStateOf("") }
    var carbs by remember { mutableStateOf("") }
    var water by remember { mutableStateOf("") }
    var components by remember { mutableStateOf<List<ComponentItem>>(emptyList()) }
    var error by remember { mutableStateOf<String?>(null) }
    var isLoading by remember { mutableStateOf(false) }
    var showComponentPicker by remember { mutableStateOf(false) }
    var allMeals by remember { mutableStateOf<List<MealData>>(emptyList()) }
    var componentSearch by remember { mutableStateOf("") }
    var isEditMode by remember { mutableStateOf(mealId != null) }

    // Загрузка блюда для редактирования
    LaunchedEffect(mealId) {
        if (mealId != null) {
            val meal = mealService.getMeal(mealId)
            if (meal != null) {
                name = meal.name
                description = meal.description
                recipe = meal.recipe
                calories = if (meal.calories > 0) meal.calories.toInt().toString() else ""
                proteins = if (meal.proteins > 0) meal.proteins.toInt().toString() else ""
                fats = if (meal.fats > 0) meal.fats.toInt().toString() else ""
                carbs = if (meal.carbs > 0) meal.carbs.toInt().toString() else ""
                water = if (meal.water > 0) meal.water.toInt().toString() else ""
                components = meal.componentsList.map {
                    ComponentItem(it.componentMealId, it.componentName, it.amount)
                }
            }
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(if (isEditMode) "Edit Meal" else "Create Meal") },
                navigationIcon = {
                    IconButton(onClick = onCancel) {
                        Icon(Icons.Default.ArrowBack, "Back")
                    }
                }
            )
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(16.dp)
                .verticalScroll(rememberScrollState())
        ) {
            OutlinedTextField(value = name, onValueChange = { name = it }, label = { Text("Name *") }, modifier = Modifier.fillMaxWidth())
            Spacer(modifier = Modifier.height(8.dp))
            OutlinedTextField(value = description, onValueChange = { description = it }, label = { Text("Description") }, modifier = Modifier.fillMaxWidth(), maxLines = 3)
            Spacer(modifier = Modifier.height(8.dp))
            OutlinedTextField(value = recipe, onValueChange = { recipe = it }, label = { Text("Recipe") }, modifier = Modifier.fillMaxWidth(), maxLines = 4)
            Spacer(modifier = Modifier.height(16.dp))

            // Компоненты
            Text("Components", fontWeight = FontWeight.SemiBold, fontSize = 16.sp)
            Text("If added, nutrition auto-calculated", color = GrayText, fontSize = 12.sp)
            Spacer(modifier = Modifier.height(8.dp))

            components.forEach { comp ->
                Row(
                    modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(comp.name, modifier = Modifier.weight(1f), fontSize = 14.sp)
                    OutlinedTextField(
                        value = comp.amount.toInt().toString(),
                        onValueChange = { value ->
                            components = components.map {
                                if (it.id == comp.id) it.copy(amount = value.toDoubleOrNull() ?: 100.0) else it
                            }
                        },
                        modifier = Modifier.width(70.dp),
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        singleLine = true
                    )
                    Text("g", fontSize = 12.sp, color = GrayText)
                    IconButton(onClick = {
                        components = components.filter { it.id != comp.id }
                    }, modifier = Modifier.size(28.dp)) {
                        Icon(Icons.Default.Close, "Remove", modifier = Modifier.size(18.dp), tint = RedError)
                    }
                }
            }

            OutlinedButton(
                onClick = { showComponentPicker = true },
                modifier = Modifier.fillMaxWidth()
            ) {
                Icon(Icons.Default.Add, null, modifier = Modifier.size(18.dp))
                Spacer(modifier = Modifier.width(4.dp))
                Text("Add Component")
            }

            Spacer(modifier = Modifier.height(16.dp))

            // КБЖУ
            Text("Nutrition (per 100g)", fontWeight = FontWeight.SemiBold, fontSize = 16.sp)
            Text("Leave 0 to auto-calculate", color = GrayText, fontSize = 12.sp)
            Spacer(modifier = Modifier.height(8.dp))
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                OutlinedTextField(value = calories, onValueChange = { calories = it }, label = { Text("Cal") }, modifier = Modifier.weight(1f), keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number), singleLine = true)
                OutlinedTextField(value = proteins, onValueChange = { proteins = it }, label = { Text("Prot") }, modifier = Modifier.weight(1f), keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number), singleLine = true)
                OutlinedTextField(value = fats, onValueChange = { fats = it }, label = { Text("Fats") }, modifier = Modifier.weight(1f), keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number), singleLine = true)
            }
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                OutlinedTextField(value = carbs, onValueChange = { carbs = it }, label = { Text("Carbs") }, modifier = Modifier.weight(1f), keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number), singleLine = true)
                OutlinedTextField(value = water, onValueChange = { water = it }, label = { Text("Water") }, modifier = Modifier.weight(1f), keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number), singleLine = true)
            }

            error?.let {
                Spacer(modifier = Modifier.height(8.dp))
                Text(it, color = RedError, fontSize = 12.sp)
            }

            Spacer(modifier = Modifier.height(16.dp))

            Button(
                onClick = {
                    if (name.isBlank()) {
                        error = "Name is required"
                        return@Button
                    }
                    isLoading = true
                    kotlinx.coroutines.MainScope().launch {
                        val userId = sessionManager.userId ?: return@launch
                        val compInputs = components.map {
                            MealComponentInput.newBuilder()
                                .setComponentMealId(it.id)
                                .setAmount(it.amount)
                                .build()
                        }
                        val success = if (mealId == null) {
                            mealService.createMeal(
                                userId, name, description, recipe, "",
                                calories.toDoubleOrNull() ?: 0.0,
                                proteins.toDoubleOrNull() ?: 0.0,
                                fats.toDoubleOrNull() ?: 0.0,
                                carbs.toDoubleOrNull() ?: 0.0,
                                water.toDoubleOrNull() ?: 0.0,
                                compInputs
                            )
                        } else {
                            mealService.updateMeal(
                                mealId, name, description, recipe, "",
                                calories.toDoubleOrNull() ?: 0.0,
                                proteins.toDoubleOrNull() ?: 0.0,
                                fats.toDoubleOrNull() ?: 0.0,
                                carbs.toDoubleOrNull() ?: 0.0,
                                water.toDoubleOrNull() ?: 0.0,
                                compInputs
                            )
                        }
                        isLoading = false
                        if (success) onSaved() else error = "Failed to save"
                    }
                },
                modifier = Modifier.fillMaxWidth().height(48.dp),
                enabled = !isLoading
            ) {
                Text(if (isLoading) "Saving..." else "Save", fontWeight = FontWeight.Bold)
            }
        }
    }

// Диалог выбора компонента
    if (showComponentPicker) {
        val addedIds = components.map { it.id }.toSet()

        // Загружаем все блюда при открытии диалога
        LaunchedEffect(showComponentPicker) {
            val userId = sessionManager.userId ?: return@LaunchedEffect
            allMeals = mealService.listMeals(userId)
        }

        val available = allMeals.filter { it.mealId != mealId && it.mealId !in addedIds }
        val filtered = if (componentSearch.isBlank()) {
            available
        } else {
            available.filter { it.name.contains(componentSearch, ignoreCase = true) }
        }

        AlertDialog(
            onDismissRequest = {
                showComponentPicker = false
                componentSearch = ""
            },
            title = { Text("Select Component") },
            text = {
                Column {
                    OutlinedTextField(
                        value = componentSearch,
                        onValueChange = { componentSearch = it },
                        placeholder = { Text("Search...") },
                        modifier = Modifier.fillMaxWidth(),
                        singleLine = true
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    if (available.isEmpty()) {
                        Text(
                            "No other meals available",
                            color = GrayText,
                            modifier = Modifier.padding(16.dp)
                        )
                    } else if (filtered.isEmpty()) {
                        Text(
                            "No meals found",
                            color = GrayText,
                            modifier = Modifier.padding(16.dp)
                        )
                    } else {
                        LazyColumn(modifier = Modifier.heightIn(max = 300.dp)) {
                            items(filtered) { meal ->
                                Surface(
                                    modifier = Modifier
                                        .fillMaxWidth()
                                        .clickable {
                                            components = components + ComponentItem(
                                                meal.mealId, meal.name, 100.0
                                            )
                                            showComponentPicker = false
                                            componentSearch = ""
                                        },
                                    shape = MaterialTheme.shapes.medium
                                ) {
                                    Row(
                                        modifier = Modifier.padding(12.dp),
                                        verticalAlignment = Alignment.CenterVertically
                                    ) {
                                        Column(modifier = Modifier.weight(1f)) {
                                            Text(meal.name, fontWeight = FontWeight.Medium)
                                            Text(
                                                "${meal.calories.toInt()} kcal",
                                                color = GrayText,
                                                fontSize = 12.sp
                                            )
                                        }
                                    }
                                }
                                HorizontalDivider()
                            }
                        }
                    }
                }
            },
            confirmButton = {},
            dismissButton = {
                TextButton(onClick = {
                    showComponentPicker = false
                    componentSearch = ""
                }) {
                    Text("Close")
                }
            }
        )
    }
}

private data class ComponentItem(
    val id: String,
    val name: String,
    val amount: Double
)