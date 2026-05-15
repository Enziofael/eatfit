package enzio.android.eatfit

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.data.remote.GrpcAuthService
import enzio.android.eatfit.domain.AuthRepository
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
import enzio.android.eatfit.ui.theme.EatfitTheme
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val sessionManager = SessionManager(applicationContext)
        val grpcService = GrpcAuthService()
        val repository = AuthRepository(grpcService, sessionManager)

        val loginViewModel = LoginViewModel(repository)
        val registerViewModel = RegisterViewModel(repository)
        val verifyEmailViewModel = VerifyEmailViewModel(repository)
        val forgotPasswordViewModel = ForgotPasswordViewModel(repository)
        val resetPasswordViewModel = ResetPasswordViewModel(repository)
        val mainViewModel = MainViewModel(repository)

        // Проверяем, залогинен ли пользователь
        val isLoggedIn = runBlocking { sessionManager.isLoggedIn.first() }
        val refreshToken = runBlocking { sessionManager.refreshToken.first() ?: "" }
        val startDestination = if (isLoggedIn) "main" else "login"

        setContent {
            EatfitTheme {
                Surface(modifier = Modifier.fillMaxSize()) {
                    val navController = rememberNavController()

                    NavHost(
                        navController = navController,
                        startDestination = startDestination
                    ) {
                        composable("login") {
                            LoginScreen(
                                viewModel = loginViewModel,
                                onLoginSuccess = {
                                    navController.navigate("main") {
                                        popUpTo(0) { inclusive = true }
                                    }
                                },
                                onRegisterClick = {
                                    navController.navigate("register")
                                },
                                onForgotPasswordClick = {
                                    navController.navigate("forgot")
                                }
                            )
                        }

                        composable("register") {
                            RegisterScreen(
                                viewModel = registerViewModel,
                                onRegisterSuccess = { userId, _, _ ->
                                    navController.navigate("verify/$userId")
                                },
                                onLoginClick = {
                                    navController.popBackStack()
                                }
                            )
                        }

                        composable("verify/{userId}") { backStackEntry ->
                            val userId = backStackEntry.arguments?.getString("userId") ?: ""
                            VerifyEmailScreen(
                                viewModel = verifyEmailViewModel,
                                userId = userId,
                                onSuccess = {
                                    navController.navigate("login") {
                                        popUpTo("login") { inclusive = true }
                                    }
                                },
                                onBack = {
                                    navController.popBackStack()
                                }
                            )
                        }

                        composable("forgot") {
                            ForgotPasswordScreen(
                                viewModel = forgotPasswordViewModel,
                                onCodeSent = { resetToken ->
                                    navController.navigate("reset/$resetToken")
                                },
                                onBack = {
                                    navController.popBackStack()
                                }
                            )
                        }

                        composable("reset/{resetToken}") { backStackEntry ->
                            val resetToken = backStackEntry.arguments?.getString("resetToken") ?: ""
                            ResetPasswordScreen(
                                viewModel = resetPasswordViewModel,
                                resetToken = resetToken,
                                onSuccess = {
                                    navController.navigate("login") {
                                        popUpTo("login") { inclusive = true }
                                    }
                                },
                                onBack = {
                                    navController.popBackStack()
                                }
                            )
                        }

                        composable("main") {
                            MainScreen(
                                viewModel = mainViewModel,
                                refreshToken = refreshToken,
                                onLogout = {
                                    navController.navigate("login") {
                                        popUpTo(0) { inclusive = true }
                                    }
                                }
                            )
                        }
                    }
                }
            }
        }
    }
}