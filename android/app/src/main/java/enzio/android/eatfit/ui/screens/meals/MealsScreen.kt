package enzio.android.eatfit.ui.screens.meals

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.proto.MealData
import enzio.android.eatfit.ui.theme.*
import java.util.Locale

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MealsScreen(
    viewModel: MealsViewModel,
    onMealClick: (String) -> Unit = {},
    onAddClick: () -> Unit = {},
    onEditClick: (String) -> Unit = {}
) {
    val state by viewModel.state.collectAsState()
    var showDeleteDialog by remember { mutableStateOf<MealData?>(null) }
    var sortExpanded by remember { mutableStateOf(false) }

    Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        // Поиск и сортировка
        Row(
            modifier = Modifier.fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically
        ) {
            OutlinedTextField(
                value = state.searchQuery,
                onValueChange = viewModel::onSearchQueryChange,
                modifier = Modifier.weight(1f),
                placeholder = { Text("Search meals...") },
                singleLine = true,
                trailingIcon = {
                    IconButton(onClick = { viewModel.search(state.searchQuery) }) {
                        Icon(Icons.Default.Search, "Search")
                    }
                }
            )
            Spacer(modifier = Modifier.width(8.dp))
            // Сортировка
            Box {
                IconButton(onClick = { sortExpanded = true }) {
                    Icon(Icons.Default.Sort, "Sort")
                }
                DropdownMenu(expanded = sortExpanded, onDismissRequest = { sortExpanded = false }) {
                    listOf(
                        "calories" to "Calories", "proteins" to "Proteins",
                        "fats" to "Fats", "carbs" to "Carbs",
                        "water" to "Water", "name" to "Name",
                        "created_at" to "Newest"
                    ).forEach { (key, label) ->
                        DropdownMenuItem(
                            text = { Text(label) },
                            onClick = {
                                viewModel.setSortBy(key)
                                sortExpanded = false
                            },
                            leadingIcon = {
                                if (state.sortBy == key) Icon(Icons.Default.Check, null)
                            }
                        )
                    }
                }
            }
            Spacer(modifier = Modifier.width(4.dp))
            // Кнопка добавления
            IconButton(onClick = onAddClick) {
                Icon(Icons.Default.Add, "Add meal", tint = GreenSuccess)
            }
        }

        Spacer(modifier = Modifier.height(12.dp))

        // Список блюд
        if (state.isLoading) {
            Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                CircularProgressIndicator()
            }
        } else if (state.meals.isEmpty()) {
            Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                Text("No meals yet", color = GrayText)
            }
        } else {
            LazyColumn(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                items(state.meals) { meal ->
                    MealCard(
                        meal = meal,
                        onClick = { onMealClick(meal.mealId) },
                        onEdit = { onEditClick(meal.mealId) },
                        onDelete = { showDeleteDialog = meal }
                    )
                }
            }
        }
    }

    // Диалог удаления
    showDeleteDialog?.let { meal ->
        AlertDialog(
            onDismissRequest = { showDeleteDialog = null },
            title = { Text("Delete meal") },
            text = { Text("Delete \"${meal.name}\"?") },
            confirmButton = {
                TextButton(onClick = {
                    viewModel.deleteMeal(meal.mealId)
                    showDeleteDialog = null
                }) {
                    Text("Delete", color = RedError)
                }
            },
            dismissButton = {
                TextButton(onClick = { showDeleteDialog = null }) {
                    Text("Cancel")
                }
            }
        )
    }
}

@Composable
private fun MealCard(
    meal: MealData,
    onClick: () -> Unit,
    onEdit: () -> Unit,
    onDelete: () -> Unit
) {
    Card(
        modifier = Modifier.fillMaxWidth().clickable(onClick = onClick),
        shape = RoundedCornerShape(12.dp)
    ) {
        Row(
            modifier = Modifier.padding(12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Иконка
            Surface(
                modifier = Modifier.size(48.dp),
                shape = RoundedCornerShape(8.dp),
                color = PurplePrimary
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Text("🍽", fontSize = 22.sp)
                }
            }

            Spacer(modifier = Modifier.width(12.dp))

            // Название и КБЖУ
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = meal.name,
                    fontWeight = FontWeight.SemiBold,
                    fontSize = 15.sp,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )
                Text(
                    text = "${meal.calories.toInt()} kcal | P:${meal.proteins.toInt()} F:${meal.fats.toInt()} C:${meal.carbs.toInt()}",
                    color = GrayText,
                    fontSize = 12.sp
                )
                if (meal.componentsCount > 0) {
                    Text(
                        text = "${meal.componentsCount} components",
                        color = GraySubText,
                        fontSize = 11.sp
                    )
                }
            }

            // Кнопки
            IconButton(onClick = onEdit, modifier = Modifier.size(32.dp)) {
                Icon(Icons.Default.Edit, "Edit", modifier = Modifier.size(18.dp), tint = PurplePrimary)
            }
            IconButton(onClick = onDelete, modifier = Modifier.size(32.dp)) {
                Icon(Icons.Default.Close, "Delete", modifier = Modifier.size(18.dp), tint = RedError)
            }
        }
    }
}