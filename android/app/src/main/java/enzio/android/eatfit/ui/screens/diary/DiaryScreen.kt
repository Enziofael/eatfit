package enzio.android.eatfit.ui.screens.diary

import androidx.compose.foundation.*
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
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.graphics.nativeCanvas
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.proto.ConsumptionRecord
import enzio.android.eatfit.ui.theme.*
import java.time.LocalDate
import java.time.format.DateTimeFormatter
import androidx.compose.animation.AnimatedVisibility

@Composable
fun DiaryScreen(
    viewModel: DiaryViewModel,
    onMealClick: (String) -> Unit = {}
) {
    val state by viewModel.state.collectAsState()
    var showProgress by remember { mutableStateOf(true) }
    var showWeightChart by remember { mutableStateOf(true) }

    Column(modifier = Modifier.fillMaxSize().padding(8.dp)) {
        // Поиск и добавление
        Row(verticalAlignment = Alignment.CenterVertically) {
            OutlinedTextField(
                value = state.searchQuery,
                onValueChange = viewModel::onSearchQueryChange,
                modifier = Modifier.weight(1f),
                placeholder = { Text("Search meals...") },
                singleLine = true,
                trailingIcon = {
                    IconButton(onClick = { viewModel.searchMeals() }) {
                        Icon(Icons.Default.Search, null)
                    }
                }
            )
            Spacer(modifier = Modifier.width(8.dp))
            IconButton(onClick = { viewModel.showAllMeals() }) {
                Icon(Icons.Default.Add, null, tint = GreenSuccess)
            }
        }

        Spacer(modifier = Modifier.height(8.dp))

        LazyColumn(verticalArrangement = Arrangement.spacedBy(8.dp)) {
            // Сворачиваемые прогресс-бары
            item {
                Card(modifier = Modifier.fillMaxWidth()) {
                    Column {
                        Row(
                            modifier = Modifier.fillMaxWidth().clickable { showProgress = !showProgress }.padding(12.dp),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Daily Progress", fontWeight = FontWeight.SemiBold, fontSize = 15.sp)
                            Icon(
                                if (showProgress) Icons.Default.ExpandLess else Icons.Default.ExpandMore,
                                null,
                                tint = GrayText
                            )
                        }
                        AnimatedVisibility(visible = showProgress) {
                            Column(modifier = Modifier.padding(start = 12.dp, end = 12.dp, bottom = 12.dp)) {
                                val dayRecords = state.groups
                                    .filter { it.date == state.selectedDate.format(DateTimeFormatter.ISO_LOCAL_DATE) }
                                    .flatMap { it.recordsList }

                                val cal = dayRecords.sumOf { it.calories }
                                val prot = dayRecords.sumOf { it.proteins }
                                val fat = dayRecords.sumOf { it.fats }
                                val carb = dayRecords.sumOf { it.carbs }
                                val water = dayRecords.sumOf { it.water }

                                val norms = state.profile?.norms
                                ProgressBar("Calories", cal, norms?.calories ?: 2000.0, "kcal", Color(0xFF667eea))
                                ProgressBar("Proteins", prot, norms?.proteins ?: 150.0, "g", Color(0xFF11998e))
                                ProgressBar("Fats", fat, norms?.fats ?: 65.0, "g", Color(0xFFf5576c))
                                ProgressBar("Carbs", carb, norms?.carbs ?: 300.0, "g", Color(0xFFff9800))
                                ProgressBar("Water", water, norms?.water ?: 2500.0, "ml", Color(0xFF2196f3))
                            }
                        }
                    }
                }
            }

            // Сворачиваемый график веса
            item {
                Card(modifier = Modifier.fillMaxWidth()) {
                    Column {
                        Row(
                            modifier = Modifier.fillMaxWidth().clickable { showWeightChart = !showWeightChart }.padding(12.dp),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Weight Chart", fontWeight = FontWeight.SemiBold, fontSize = 15.sp)
                            Icon(
                                if (showWeightChart) Icons.Default.ExpandLess else Icons.Default.ExpandMore,
                                null,
                                tint = GrayText
                            )
                        }
                        AnimatedVisibility(visible = showWeightChart) {
                            Column(modifier = Modifier.padding(start = 12.dp, end = 12.dp, bottom = 12.dp)) {
                                Row(
                                    modifier = Modifier.fillMaxWidth(),
                                    horizontalArrangement = Arrangement.Center,
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    IconButton(onClick = { viewModel.prevMonth() }, modifier = Modifier.size(24.dp)) {
                                        Icon(Icons.Default.ChevronLeft, null)
                                    }
                                    Text(
                                        state.currentMonth.format(DateTimeFormatter.ofPattern("MMMM yyyy")),
                                        fontSize = 13.sp, fontWeight = FontWeight.Medium
                                    )
                                    IconButton(onClick = { viewModel.nextMonth() }, modifier = Modifier.size(24.dp)) {
                                        Icon(Icons.Default.ChevronRight, null)
                                    }
                                }
                                WeightChart(
                                    data = state.weightHistory.filter {
                                        it.date.year == state.currentMonth.year && it.date.month == state.currentMonth.month
                                    },
                                    modifier = Modifier.fillMaxWidth().height(150.dp)
                                )
                            }
                        }
                    }
                }
            }

            // Заголовок истории
            item {
                Text("History", fontWeight = FontWeight.SemiBold, fontSize = 15.sp, modifier = Modifier.padding(top = 4.dp))
            }

            // История потребления
            if (state.groups.isEmpty()) {
                item {
                    Text("No records yet", color = GrayText, fontSize = 13.sp, modifier = Modifier.padding(8.dp))
                }
            } else {
                state.groups.forEach { group ->
                    item {
                        Text(
                            group.date,
                            fontWeight = FontWeight.SemiBold,
                            fontSize = 13.sp,
                            modifier = Modifier.padding(top = 8.dp, bottom = 4.dp)
                        )
                    }
                    items(group.recordsList) { record ->
                        ConsumptionCard(
                            record = record,
                            onMealClick = { onMealClick(record.mealId) },
                            onDelete = { viewModel.deleteConsumption(record.recordId) }
                        )
                    }
                }
            }
        }
    }

    // Диалог выбора блюда
    if (state.showMealPicker) {
        AlertDialog(
            onDismissRequest = { viewModel.hideMealPicker() },
            title = { Text("Select Meal") },
            text = {
                Column {
                    OutlinedTextField(
                        value = state.amount,
                        onValueChange = viewModel::onAmountChange,
                        label = { Text("Amount (g)") },
                        singleLine = true
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    LazyColumn(modifier = Modifier.heightIn(max = 300.dp)) {
                        items(state.searchResults) { meal ->
                            Surface(
                                modifier = Modifier.fillMaxWidth().clickable { viewModel.selectMeal(meal) },
                                color = if (state.selectedMeal?.mealId == meal.mealId) PurplePrimary.copy(alpha = 0.1f) else Color.Transparent
                            ) {
                                Row(modifier = Modifier.padding(12.dp)) {
                                    Column(modifier = Modifier.weight(1f)) {
                                        Text(meal.name, fontWeight = FontWeight.Medium)
                                        Text("${meal.calories.toInt()} kcal/100g", color = GrayText, fontSize = 12.sp)
                                    }
                                }
                            }
                            HorizontalDivider()
                        }
                    }
                }
            },
            confirmButton = {
                TextButton(onClick = { viewModel.addConsumption() }, enabled = state.selectedMeal != null) {
                    Text("Add")
                }
            },
            dismissButton = {
                TextButton(onClick = { viewModel.hideMealPicker() }) { Text("Cancel") }
            }
        )
    }
}

@Composable
private fun ProgressBar(label: String, value: Double, max: Double, unit: String, color: Color) {
    val pct = if (max > 0) (value / max).coerceIn(0.0, 1.5).toFloat() else 0f
    Column(modifier = Modifier.padding(vertical = 4.dp)) {
        Text(label, fontSize = 11.sp, color = GrayText)
        Box(modifier = Modifier.fillMaxWidth().height(18.dp).clip(RoundedCornerShape(9.dp)).background(Color(0xFFEEEEEE))) {
            Box(
                modifier = Modifier
                    .fillMaxHeight()
                    .fillMaxWidth(pct)
                    .clip(RoundedCornerShape(9.dp))
                    .background(color)
            )
            Text(
                "${value.toInt()}/${max.toInt()} $unit",
                fontSize = 10.sp,
                fontWeight = FontWeight.SemiBold,
                modifier = Modifier.align(Alignment.Center)
            )
        }
    }
}

@Composable
private fun WeightChart(data: List<WeightPoint>, modifier: Modifier) {
    if (data.size < 2) {
        Box(modifier = modifier, contentAlignment = Alignment.Center) {
            Text("Not enough data", color = GrayText, fontSize = 12.sp)
        }
        return
    }

    Canvas(modifier = modifier.padding(8.dp)) {
        val min = data.minOf { it.weight }
        val max = data.maxOf { it.weight }
        val range = (max - min).coerceAtLeast(1.0)
        val pad = 8f

        val points = data.mapIndexed { i, point ->
            val x = pad + (i * (size.width - pad * 2) / (data.size - 1))
            val y = size.height - pad - ((point.weight - min) / range * (size.height - pad * 2)).toFloat()
            Offset(x, y)
        }

        // Линия
        val path = Path().apply {
            moveTo(points.first().x, points.first().y)
            points.drop(1).forEach { lineTo(it.x, it.y) }
        }
        drawPath(path, PurplePrimary, style = Stroke(width = 2.dp.toPx(), cap = StrokeCap.Round))

        // Точки
        points.forEach { point ->
            drawCircle(PurplePrimary, radius = 3.dp.toPx(), center = point)
            drawCircle(Color.White, radius = 1.5.dp.toPx(), center = point)
        }
    }
}

@Composable
private fun ConsumptionCard(
    record: ConsumptionRecord,
    onMealClick: () -> Unit,
    onDelete: () -> Unit
) {
    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(8.dp)) {
        Row(
            modifier = Modifier.padding(12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column(modifier = Modifier.weight(1f)) {
                Text("${record.mealName} — ${record.amount.toInt()}g", fontWeight = FontWeight.Medium)
                Text(
                    "Cal: ${record.calories.toInt()} | P: ${record.proteins.toInt()} | F: ${record.fats.toInt()} | C: ${record.carbs.toInt()}",
                    color = GrayText,
                    fontSize = 11.sp
                )
            }
            if (record.mealId.isNotEmpty()) {
                IconButton(onClick = onMealClick, modifier = Modifier.size(28.dp)) {
                    Icon(Icons.Default.Visibility, null, modifier = Modifier.size(18.dp), tint = PurplePrimary)
                }
            }
            IconButton(onClick = onDelete, modifier = Modifier.size(28.dp)) {
                Icon(Icons.Default.Close, null, modifier = Modifier.size(18.dp), tint = RedError)
            }
        }
    }
}

@Composable
private fun VerticalDivider() {
    Box(modifier = Modifier.fillMaxHeight().width(1.dp).background(Color(0xFFE0E0E0)))
}