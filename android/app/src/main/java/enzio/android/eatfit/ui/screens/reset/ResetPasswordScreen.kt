package enzio.android.eatfit.ui.screens.reset

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import enzio.android.eatfit.ui.theme.*

@Composable
fun ResetPasswordScreen(
    viewModel: ResetPasswordViewModel,
    resetToken: String,
    onSuccess: () -> Unit,
    onBack: () -> Unit
) {
    val state by viewModel.state.collectAsState()

    Column(
        modifier = Modifier.fillMaxSize().padding(32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text("Reset Password", fontSize = 24.sp, fontWeight = FontWeight.Bold)
        Spacer(modifier = Modifier.height(24.dp))

        OutlinedTextField(
            value = state.code,
            onValueChange = viewModel::onCodeChange,
            label = { Text("Verification Code") },
            modifier = Modifier.fillMaxWidth(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number)
        )
        Spacer(modifier = Modifier.height(12.dp))

        OutlinedTextField(
            value = state.newPassword,
            onValueChange = viewModel::onPasswordChange,
            label = { Text("New Password") },
            modifier = Modifier.fillMaxWidth(),
            visualTransformation = PasswordVisualTransformation()
        )
        Spacer(modifier = Modifier.height(12.dp))

        OutlinedTextField(
            value = state.confirmPassword,
            onValueChange = viewModel::onConfirmChange,
            label = { Text("Confirm Password") },
            modifier = Modifier.fillMaxWidth(),
            visualTransformation = PasswordVisualTransformation()
        )

        val error = state.error
        if (error != null) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(error, color = RedError, fontSize = 12.sp)
        }

        Spacer(modifier = Modifier.height(24.dp))

        Button(
            onClick = { viewModel.resetPassword(resetToken, onSuccess) },
            modifier = Modifier.fillMaxWidth().height(48.dp),
            enabled = !state.isLoading
        ) {
            Text("Reset Password", fontWeight = FontWeight.Bold)
        }

        Spacer(modifier = Modifier.height(12.dp))

        TextButton(onClick = onBack) {
            Text("← Back to login", color = PurplePrimary)
        }
    }
}