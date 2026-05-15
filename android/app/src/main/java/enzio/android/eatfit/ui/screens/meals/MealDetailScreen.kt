package enzio.android.eatfit.ui.screens.meals

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.data.remote.GrpcMealService
import enzio.android.eatfit.proto.MealData
import enzio.android.eatfit.ui.theme.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MealDetailScreen(
    mealService: GrpcMealService,
    mealId: String,
    onBack: () -> Unit
) {
    var meal by remember { mutableStateOf<MealData?>(null) }
    var isLoading by remember { mutableStateOf(true) }

    LaunchedEffect(mealId) {
        meal = mealService.getMeal(mealId)
        isLoading = false
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(meal?.name ?: "Meal") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, "Back")
                    }
                }
            )
        }
    ) { padding ->
        if (isLoading) {
            Box(modifier = Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                CircularProgressIndicator()
            }
        } else if (meal == null) {
            Box(modifier = Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                Text("Meal not found", color = GrayText)
            }
        } else {
            val m = meal!!
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .padding(16.dp)
                    .verticalScroll(rememberScrollState())
            ) {
                // Заголовок
                Surface(
                    modifier = Modifier.size(72.dp),
                    shape = MaterialTheme.shapes.extraLarge,
                    color = PurplePrimary
                ) {
                    Box(contentAlignment = Alignment.Center) {
                        Text("🍽", fontSize = 32.sp)
                    }
                }
                Spacer(modifier = Modifier.height(12.dp))
                Text(m.name, fontSize = 26.sp, fontWeight = FontWeight.Bold)
                Text("${m.calories.toInt()} kcal per 100g", color = GrayText, fontSize = 15.sp)
                Spacer(modifier = Modifier.height(20.dp))

                // КБЖУ
                Card(modifier = Modifier.fillMaxWidth()) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(16.dp),
                        horizontalArrangement = Arrangement.SpaceEvenly
                    ) {
                        NutritionItem("Calories", m.calories)
                        NutritionItem("Proteins", m.proteins)
                        NutritionItem("Fats", m.fats)
                        NutritionItem("Carbs", m.carbs)
                        NutritionItem("Water", m.water)
                    }
                }

                // Состав
                if (m.componentsList.isNotEmpty()) {
                    Spacer(modifier = Modifier.height(16.dp))
                    Card(modifier = Modifier.fillMaxWidth()) {
                        Column(modifier = Modifier.padding(16.dp)) {
                            Text("Components", fontWeight = FontWeight.SemiBold, fontSize = 16.sp)
                            Spacer(modifier = Modifier.height(8.dp))
                            m.componentsList.forEach { comp ->
                                Row(
                                    modifier = Modifier.fillMaxWidth().padding(vertical = 6.dp),
                                    horizontalArrangement = Arrangement.SpaceBetween
                                ) {
                                    Column(modifier = Modifier.weight(1f)) {
                                        Text(comp.componentName, fontWeight = FontWeight.Medium)
                                        Text(
                                            "Cal: ${comp.calories.toInt()} | P: ${comp.proteins.toInt()} | F: ${comp.fats.toInt()} | C: ${comp.carbs.toInt()}",
                                            color = GrayText,
                                            fontSize = 11.sp
                                        )
                                    }
                                    Text("${comp.amount.toInt()}g", color = GraySubText, fontSize = 13.sp)
                                }
                                HorizontalDivider()
                            }
                        }
                    }
                }

                // Описание
                if (m.description.isNotBlank()) {
                    Spacer(modifier = Modifier.height(16.dp))
                    Card(modifier = Modifier.fillMaxWidth()) {
                        Column(modifier = Modifier.padding(16.dp)) {
                            Text("Description", fontWeight = FontWeight.SemiBold, fontSize = 16.sp)
                            Spacer(modifier = Modifier.height(4.dp))
                            Text(m.description, color = GrayText)
                        }
                    }
                }

                // Рецепт
                if (m.recipe.isNotBlank()) {
                    Spacer(modifier = Modifier.height(16.dp))
                    Card(modifier = Modifier.fillMaxWidth()) {
                        Column(modifier = Modifier.padding(16.dp)) {
                            Text("Recipe", fontWeight = FontWeight.SemiBold, fontSize = 16.sp)
                            Spacer(modifier = Modifier.height(4.dp))
                            Text(m.recipe, color = GrayText)
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun NutritionItem(label: String, value: Double) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text(label, color = GrayText, fontSize = 11.sp)
        Text(value.toInt().toString(), fontSize = 18.sp, fontWeight = FontWeight.SemiBold)
    }
}