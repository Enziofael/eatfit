package enzio.android.eatfit.ui.navigation

import androidx.compose.runtime.*
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import enzio.android.eatfit.ui.screens.forgot.ForgotPasswordScreen
import enzio.android.eatfit.ui.screens.forgot.ForgotPasswordViewModel
import enzio.android.eatfit.ui.screens.login.LoginScreen
import enzio.android.eatfit.ui.screens.login.LoginViewModel
import enzio.android.eatfit.ui.screens.main.MainScreen
import enzio.android.eatfit.ui.screens.main.MainViewModel
import enzio.android.eatfit.ui.screens.register.RegisterScreen
import enzio.android.eatfit.ui.screens.register.RegisterViewModel
import enzio.android.eatfit.ui.screens.reset.ResetPasswordScreen
import enzio.android.eatfit.ui.screens.reset.ResetPasswordViewModel
import enzio.android.eatfit.ui.screens.verify.VerifyEmailScreen
import enzio.android.eatfit.ui.screens.verify.VerifyEmailViewModel

object Routes {
    const val LOGIN = "login"
    const val REGISTER = "register"
    const val VERIFY = "verify/{userId}"
    const val FORGOT = "forgot"
    const val RESET = "reset/{resetToken}"
    const val MAIN = "main"
}

@Composable
fun NavGraph(
    navController: NavHostController,
    loginViewModel: LoginViewModel,
    registerViewModel: RegisterViewModel,
    verifyEmailViewModel: VerifyEmailViewModel,
    forgotPasswordViewModel: ForgotPasswordViewModel,
    resetPasswordViewModel: ResetPasswordViewModel,
    mainViewModel: MainViewModel,
    refreshToken: String
) {
    NavHost(navController = navController, startDestination = Routes.LOGIN) {
        composable(Routes.LOGIN) {
            LoginScreen(
                viewModel = loginViewModel,
                onLoginSuccess = { navController.navigate(Routes.MAIN) { popUpTo(0) } },
                onRegisterClick = { navController.navigate(Routes.REGISTER) },
                onForgotPasswordClick = { navController.navigate(Routes.FORGOT) }
            )
        }

        composable(Routes.REGISTER) {
            RegisterScreen(
                viewModel = registerViewModel,
                onRegisterSuccess = { userId, email, login ->
                    navController.navigate("verify/$userId")
                },
                onLoginClick = { navController.popBackStack() }
            )
        }

        composable(Routes.VERIFY) { backStackEntry ->
            val userId = backStackEntry.arguments?.getString("userId") ?: ""
            VerifyEmailScreen(
                viewModel = verifyEmailViewModel,
                userId = userId,
                onSuccess = {
                    navController.navigate(Routes.LOGIN) { popUpTo(Routes.LOGIN) { inclusive = true } }
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(Routes.FORGOT) {
            ForgotPasswordScreen(
                viewModel = forgotPasswordViewModel,
                onCodeSent = { resetToken ->
                    navController.navigate("reset/$resetToken")
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(Routes.RESET) { backStackEntry ->
            val resetToken = backStackEntry.arguments?.getString("resetToken") ?: ""
            ResetPasswordScreen(
                viewModel = resetPasswordViewModel,
                resetToken = resetToken,
                onSuccess = {
                    navController.navigate(Routes.LOGIN) { popUpTo(Routes.LOGIN) { inclusive = true } }
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(Routes.MAIN) {
            MainScreen(
                viewModel = mainViewModel,
                refreshToken = refreshToken,
                onLogout = {
                    navController.navigate(Routes.LOGIN) { popUpTo(0) }
                }
            )
        }
    }
}