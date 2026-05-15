package enzio.android.eatfit

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.data.remote.GrpcAuthService
import enzio.android.eatfit.data.remote.GrpcMealService
import enzio.android.eatfit.data.remote.GrpcProfileService
import enzio.android.eatfit.data.remote.GrpcDiaryService
import enzio.android.eatfit.domain.AuthRepository
import enzio.android.eatfit.domain.ProfileRepository
import enzio.android.eatfit.ui.screens.forgot.ForgotPasswordScreen
import enzio.android.eatfit.ui.screens.forgot.ForgotPasswordViewModel
import enzio.android.eatfit.ui.screens.login.LoginScreen
import enzio.android.eatfit.ui.screens.login.LoginViewModel
import enzio.android.eatfit.ui.screens.main.MainScreen
import enzio.android.eatfit.ui.screens.main.MainViewModel
import enzio.android.eatfit.ui.screens.meals.MealDetailScreen
import enzio.android.eatfit.ui.screens.meals.MealEditorScreen
import enzio.android.eatfit.ui.screens.meals.MealsViewModel
import enzio.android.eatfit.ui.screens.profile.ProfileViewModel
import enzio.android.eatfit.ui.screens.register.RegisterScreen
import enzio.android.eatfit.ui.screens.register.RegisterViewModel
import enzio.android.eatfit.ui.screens.reset.ResetPasswordScreen
import enzio.android.eatfit.ui.screens.reset.ResetPasswordViewModel
import enzio.android.eatfit.ui.screens.verify.VerifyEmailScreen
import enzio.android.eatfit.ui.screens.verify.VerifyEmailViewModel
import enzio.android.eatfit.ui.screens.settings.SettingsScreen
import enzio.android.eatfit.ui.screens.settings.SettingsViewModel
import enzio.android.eatfit.ui.screens.diary.DiaryScreen
import enzio.android.eatfit.ui.screens.diary.DiaryViewModel
import enzio.android.eatfit.ui.theme.EatfitTheme
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val sessionManager = SessionManager(applicationContext)
        val grpcAuthService = GrpcAuthService()
        val authRepository = AuthRepository(grpcAuthService, sessionManager)

        val grpcProfileService = GrpcProfileService()
        val profileRepository = ProfileRepository(grpcProfileService)

        val grpcMealService = GrpcMealService()
        val mealsViewModel = MealsViewModel(grpcMealService, sessionManager)

        val grpcDiaryService = GrpcDiaryService()
        val diaryViewModel = DiaryViewModel(grpcDiaryService, grpcMealService, grpcProfileService, sessionManager)

        val loginViewModel = LoginViewModel(authRepository)
        val registerViewModel = RegisterViewModel(authRepository)
        val verifyEmailViewModel = VerifyEmailViewModel(authRepository)
        val forgotPasswordViewModel = ForgotPasswordViewModel(authRepository)
        val resetPasswordViewModel = ResetPasswordViewModel(authRepository)
        val mainViewModel = MainViewModel(authRepository)
        val profileViewModel = ProfileViewModel(profileRepository, sessionManager)

        val isLoggedIn = runBlocking { sessionManager.isLoggedIn.first() }
        val refreshToken = runBlocking { sessionManager.refreshToken.first() ?: "" }
        val startDestination = if (isLoggedIn) "main" else "login"

        val settingsViewModel = SettingsViewModel(grpcProfileService, grpcAuthService, sessionManager)

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
                                    profileViewModel.loadProfile()
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
                                profileViewModel = profileViewModel,
                                mealsViewModel = mealsViewModel,
                                diaryViewModel = diaryViewModel,
                                onLogout = {
                                    navController.navigate("login") {
                                        popUpTo(0) { inclusive = true }
                                    }
                                },
                                onSettingsClick = {
                                    navController.navigate("settings")
                                },
                                onMealAdd = {
                                    navController.navigate("meal_editor")
                                },
                                onMealEdit = { mealId ->
                                    navController.navigate("meal_editor/$mealId")
                                },
                                onMealClick = { mealId ->
                                    navController.navigate("meal_detail/$mealId")
                                }
                            )
                        }

                        composable("settings") {
                            SettingsScreen(
                                viewModel = settingsViewModel,
                                onBack = {
                                    profileViewModel.loadProfile()
                                    navController.popBackStack()
                                }
                            )
                        }

                        composable("meal_editor") {
                            MealEditorScreen(
                                mealService = grpcMealService,
                                sessionManager = sessionManager,
                                onSaved = {
                                    mealsViewModel.loadMeals()
                                    navController.popBackStack()
                                },
                                onCancel = { navController.popBackStack() }
                            )
                        }
                        composable("meal_editor/{mealId}") { entry ->
                            val mealId = entry.arguments?.getString("mealId")
                            MealEditorScreen(
                                mealService = grpcMealService,
                                sessionManager = sessionManager,
                                mealId = mealId,
                                onSaved = {
                                    mealsViewModel.loadMeals()
                                    navController.popBackStack()
                                },
                                onCancel = { navController.popBackStack() }
                            )
                        }
                        composable("meal_detail/{mealId}") { entry ->
                            val mealId = entry.arguments?.getString("mealId") ?: ""
                            MealDetailScreen(
                                mealService = grpcMealService,
                                mealId = mealId,
                                onBack = { navController.popBackStack() }
                            )
                        }
                    }
                }
            }
        }
    }
}