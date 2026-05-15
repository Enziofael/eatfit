package enzio.android.eatfit.ui.screens.main

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.ui.screens.diary.DiaryScreen
import enzio.android.eatfit.ui.screens.feed.FeedScreen
import enzio.android.eatfit.ui.screens.meals.MealsScreen
import enzio.android.eatfit.ui.screens.messages.MessagesScreen
import enzio.android.eatfit.ui.screens.profile.ProfileScreen
import enzio.android.eatfit.ui.screens.profile.ProfileViewModel
import enzio.android.eatfit.ui.screens.meals.MealsScreen
import enzio.android.eatfit.ui.screens.meals.MealsViewModel
import enzio.android.eatfit.ui.screens.diary.DiaryScreen
import enzio.android.eatfit.ui.screens.diary.DiaryViewModel
import enzio.android.eatfit.ui.theme.*

enum class BottomNavItem(
    val label: String,
    val icon: ImageVector
) {
    FEED("Feed", Icons.Default.Home),
    DIARY("Diary", Icons.Default.DateRange),
    MEALS("Meals", Icons.Default.Restaurant),
    MESSAGES("Messages", Icons.Default.Email),
    PROFILE("Profile", Icons.Default.Person)
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MainScreen(
    viewModel: MainViewModel,
    refreshToken: String,
    profileViewModel: ProfileViewModel,
    mealsViewModel: MealsViewModel,
    diaryViewModel: DiaryViewModel,

    onLogout: () -> Unit,
    onSettingsClick: () -> Unit = {},
    onMealAdd: () -> Unit = {},
    onMealEdit: (String) -> Unit = {},
    onMealClick: (String) -> Unit = {}
) {
    var selectedTab by remember { mutableStateOf(BottomNavItem.FEED) }
    var showLogoutDialog by remember { mutableStateOf(false) }

    Scaffold(
        bottomBar = {
            NavigationBar {
                BottomNavItem.entries.forEach { item ->
                    NavigationBarItem(
                        icon = { Icon(item.icon, contentDescription = item.label) },
                        label = { Text(item.label, fontSize = 11.sp) },
                        selected = selectedTab == item,
                        onClick = { selectedTab = item }
                    )
                }
            }
        }
    ) { paddingValues ->
        Box(modifier = Modifier.padding(paddingValues)) {
            when (selectedTab) {
                BottomNavItem.FEED -> FeedScreen()
                BottomNavItem.DIARY -> DiaryScreen(
                    viewModel = diaryViewModel,
                    onMealClick = onMealClick  // ← использовать переданный колбэк
                )
                BottomNavItem.MEALS -> MealsScreen(
                    viewModel = mealsViewModel,
                    onAddClick = onMealAdd,
                    onEditClick = onMealEdit,
                    onMealClick = onMealClick,
                )
                BottomNavItem.MESSAGES -> MessagesScreen()
                BottomNavItem.PROFILE -> ProfileScreen(
                    viewModel = profileViewModel,
                    onLogoutClick = { showLogoutDialog = true },
                    onSettingsClick = onSettingsClick
                )
            }
        }
    }

    if (showLogoutDialog) {
        AlertDialog(
            onDismissRequest = { showLogoutDialog = false },
            title = { Text("Logout") },
            text = { Text("Are you sure you want to logout?") },
            confirmButton = {
                TextButton(onClick = {
                    showLogoutDialog = false
                    viewModel.logout(refreshToken, onLogout)
                }) {
                    Text("Yes")
                }
            },
            dismissButton = {
                TextButton(onClick = { showLogoutDialog = false }) {
                    Text("No")
                }
            }
        )
    }
}