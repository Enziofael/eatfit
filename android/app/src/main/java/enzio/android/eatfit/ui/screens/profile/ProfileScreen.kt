package enzio.android.eatfit.ui.screens.profile

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.ExitToApp
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.ui.theme.*

@Composable
fun ProfileScreen(
    viewModel: ProfileViewModel,
    onLogoutClick: () -> Unit,
    onSettingsClick: () -> Unit = {}
) {
    val state by viewModel.state.collectAsState()
    var showWeightDialog by remember { mutableStateOf(false) }
    var showNormsDialog by remember { mutableStateOf(false) }

    if (state.isLoading) {
        Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            CircularProgressIndicator()
        }
        return
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp)
            .verticalScroll(rememberScrollState())
    ) {
        // Шапка профиля
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier.fillMaxWidth()
        ) {
            Surface(
                modifier = Modifier.size(64.dp),
                shape = MaterialTheme.shapes.extraLarge,
                color = PurplePrimary
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Text("👤", fontSize = 28.sp)
                }
            }
            Spacer(modifier = Modifier.width(12.dp))
            Column(modifier = Modifier.weight(1f)) {
                Text(state.displayName, fontSize = 22.sp, fontWeight = FontWeight.Bold)
                if (state.bio.isNotEmpty()) {
                    Text(state.bio, color = GrayText, fontSize = 14.sp)
                }
            }
            // Иконки настроек и выхода
            IconButton(onClick = onSettingsClick, modifier = Modifier.size(32.dp)) {
                Icon(
                    Icons.Default.Settings,
                    contentDescription = "Settings",
                    modifier = Modifier.size(20.dp),
                    tint = GraySubText
                )
            }
            IconButton(onClick = onLogoutClick, modifier = Modifier.size(32.dp)) {
                Icon(
                    Icons.Default.ExitToApp,
                    contentDescription = "Logout",
                    modifier = Modifier.size(20.dp),
                    tint = RedError
                )
            }
        }

        Spacer(modifier = Modifier.height(24.dp))

        // Информация
        Card(modifier = Modifier.fillMaxWidth()) {
            Column(modifier = Modifier.padding(16.dp)) {
                Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                    ProfileInfoItem("Height", state.height)
                    Column {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text("Weight", color = GrayText, fontSize = 12.sp)
                            IconButton(
                                onClick = { showWeightDialog = true },
                                modifier = Modifier.size(20.dp)
                            ) {
                                Icon(
                                    Icons.Default.Edit,
                                    contentDescription = "Update weight",
                                    modifier = Modifier.size(14.dp),
                                    tint = PurplePrimary
                                )
                            }
                        }
                        Text(state.weight, fontSize = 16.sp, fontWeight = FontWeight.Medium)
                    }
                }
                Spacer(modifier = Modifier.height(12.dp))
                Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                    ProfileInfoItem("Age", state.age)
                    ProfileInfoItem("Gender", state.gender)
                }
                Spacer(modifier = Modifier.height(12.dp))
                ProfileInfoItem("Birth Date", state.birthDate)
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        // Нормы
        Card(modifier = Modifier.fillMaxWidth()) {
            Column(modifier = Modifier.padding(16.dp)) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Text("Nutrition Norms", fontWeight = FontWeight.SemiBold, fontSize = 16.sp)
                    Spacer(modifier = Modifier.width(8.dp))
                    IconButton(
                        onClick = { showNormsDialog = true },
                        modifier = Modifier.size(20.dp)
                    ) {
                        Icon(
                            Icons.Default.Edit,
                            contentDescription = "Update norms",
                            modifier = Modifier.size(14.dp),
                            tint = PurplePrimary
                        )
                    }
                }
                Spacer(modifier = Modifier.height(12.dp))
                Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceEvenly) {
                    NormItem("Calories", state.calories)
                    NormItem("Proteins", state.proteins)
                    NormItem("Fats", state.fats)
                    NormItem("Carbs", state.carbs)
                    NormItem("Water", state.water)
                }
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        // Публикации
        Card(modifier = Modifier.fillMaxWidth()) {
            Column(modifier = Modifier.padding(16.dp)) {
                Text("My Publications", fontWeight = FontWeight.SemiBold, fontSize = 16.sp)
                Spacer(modifier = Modifier.height(8.dp))
                Text("Your publications will appear here", color = GrayText, fontSize = 13.sp)
            }
        }

        Spacer(modifier = Modifier.height(16.dp))
    }

    // Диалог обновления веса
    if (showWeightDialog) {
        var weightInput by remember { mutableStateOf("") }
        AlertDialog(
            onDismissRequest = { showWeightDialog = false },
            title = { Text("Update Weight") },
            text = {
                OutlinedTextField(
                    value = weightInput,
                    onValueChange = { weightInput = it },
                    label = { Text("Weight (kg)") },
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                    singleLine = true
                )
            },
            confirmButton = {
                TextButton(onClick = {
                    weightInput.toDoubleOrNull()?.let { weight ->
                        viewModel.updateWeight(weight)
                    }
                    showWeightDialog = false
                }) {
                    Text("Save")
                }
            },
            dismissButton = {
                TextButton(onClick = { showWeightDialog = false }) {
                    Text("Cancel")
                }
            }
        )
    }

    // Диалог обновления норм
    if (showNormsDialog) {
        var calories by remember { mutableStateOf("") }
        var proteins by remember { mutableStateOf("") }
        var fats by remember { mutableStateOf("") }
        var carbs by remember { mutableStateOf("") }
        var water by remember { mutableStateOf("") }

        AlertDialog(
            onDismissRequest = { showNormsDialog = false },
            title = { Text("Update Norms") },
            text = {
                Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    OutlinedTextField(
                        value = calories,
                        onValueChange = { calories = it },
                        label = { Text("Calories (kcal)") },
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        singleLine = true
                    )
                    OutlinedTextField(
                        value = proteins,
                        onValueChange = { proteins = it },
                        label = { Text("Proteins (g)") },
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        singleLine = true
                    )
                    OutlinedTextField(
                        value = fats,
                        onValueChange = { fats = it },
                        label = { Text("Fats (g)") },
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        singleLine = true
                    )
                    OutlinedTextField(
                        value = carbs,
                        onValueChange = { carbs = it },
                        label = { Text("Carbs (g)") },
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        singleLine = true
                    )
                    OutlinedTextField(
                        value = water,
                        onValueChange = { water = it },
                        label = { Text("Water (ml)") },
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        singleLine = true
                    )
                }
            },
            confirmButton = {
                TextButton(onClick = {
                    val cal = calories.toDoubleOrNull() ?: 0.0
                    val prot = proteins.toDoubleOrNull() ?: 0.0
                    val fat = fats.toDoubleOrNull() ?: 0.0
                    val carb = carbs.toDoubleOrNull() ?: 0.0
                    val wat = water.toDoubleOrNull() ?: 0.0
                    viewModel.updateNorms(cal, prot, fat, carb, wat)
                    showNormsDialog = false
                }) {
                    Text("Save")
                }
            },
            dismissButton = {
                TextButton(onClick = { showNormsDialog = false }) {
                    Text("Cancel")
                }
            }
        )
    }
}

@Composable
private fun ProfileInfoItem(label: String, value: String) {
    Column {
        Text(label, color = GrayText, fontSize = 12.sp)
        Text(value, fontSize = 16.sp, fontWeight = FontWeight.Medium)
    }
}

@Composable
private fun NormItem(label: String, value: String) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text(label, color = GrayText, fontSize = 11.sp)
        Text(value, fontSize = 13.sp, fontWeight = FontWeight.Medium)
    }
}