package enzio.android.eatfit.ui.screens.forgot

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
fun ForgotPasswordScreen(
    viewModel: ForgotPasswordViewModel,
    onCodeSent: (String) -> Unit,
    onBack: () -> Unit
) {
    val state by viewModel.state.collectAsState()

    Column(
        modifier = Modifier.fillMaxSize().padding(32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text("Forgot Password", fontSize = 24.sp, fontWeight = FontWeight.Bold)
        Spacer(modifier = Modifier.height(16.dp))
        Text("Enter your email or login to receive reset code", color = GrayText)
        Spacer(modifier = Modifier.height(24.dp))

        OutlinedTextField(
            value = state.loginIdentifier,
            onValueChange = viewModel::onIdentifierChange,
            label = { Text("Email or Login") },
            modifier = Modifier.fillMaxWidth()
        )

        val error = state.error
        if (error != null) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(error, color = RedError, fontSize = 12.sp)
        }
        val message = state.message
        if (message != null) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(message, color = GreenSuccess, fontSize = 12.sp)
        }

        Spacer(modifier = Modifier.height(24.dp))

        Button(
            onClick = { viewModel.sendResetCode(onCodeSent) },
            modifier = Modifier.fillMaxWidth().height(48.dp),
            enabled = !state.isLoading
        ) {
            Text("Send Reset Code", fontWeight = FontWeight.Bold)
        }

        Spacer(modifier = Modifier.height(12.dp))

        TextButton(onClick = onBack) {
            Text("← Back to login", color = PurplePrimary)
        }
    }
}