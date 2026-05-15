package enzio.android.eatfit.data.model

data class LoginRequest(
    val loginIdentifier: String,
    val password: String
)

data class RegisterRequest(
    val email: String,
    val login: String,
    val password: String,
    val passwordConfirmation: String
)

data class AuthResult(
    val success: Boolean,
    val message: String = "",
    val accessToken: String = "",
    val refreshToken: String = "",
    val userId: String = "",
    val email: String = "",
    val login: String = "",
    val resetToken: String = ""
)