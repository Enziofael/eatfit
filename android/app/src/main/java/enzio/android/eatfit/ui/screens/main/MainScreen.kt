package enzio.android.eatfit.ui.screens.main

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.ui.theme.*

@Composable
fun MainScreen(
    viewModel: MainViewModel,
    refreshToken: String,
    onLogout: () -> Unit
) {
    var isLoading by remember { mutableStateOf(false) }

    Column(
        modifier = Modifier.fillMaxSize(),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = "Welcome to Eatfit!",
            fontSize = 32.sp,
            fontWeight = FontWeight.Bold,
            color = PurplePrimary
        )

        Spacer(modifier = Modifier.height(32.dp))

        Button(
            onClick = {
                isLoading = true
                viewModel.logout(refreshToken) {
                    isLoading = false
                    onLogout()
                }
            },
            modifier = Modifier
                .width(200.dp)
                .height(48.dp),
            enabled = !isLoading
        ) {
            Text(
                text = if (isLoading) "Logging out..." else "Logout",
                fontWeight = FontWeight.Bold
            )
        }
    }
}