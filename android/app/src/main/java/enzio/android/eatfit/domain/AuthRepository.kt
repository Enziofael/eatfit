package enzio.android.eatfit.domain

import enzio.android.eatfit.data.local.SessionManager
import enzio.android.eatfit.data.model.AuthResult
import enzio.android.eatfit.data.model.LoginRequest
import enzio.android.eatfit.data.model.RegisterRequest
import enzio.android.eatfit.data.remote.GrpcAuthService

class AuthRepository(
    private val grpcService: GrpcAuthService,
    private val sessionManager: SessionManager
) {
    suspend fun register(request: RegisterRequest): AuthResult {
        return grpcService.register(request)
    }

    suspend fun verifyEmail(userId: String, code: String): AuthResult {
        return grpcService.verifyEmail(userId, code)
    }

    suspend fun login(request: LoginRequest): AuthResult {
        val result = grpcService.login(request)
        if (result.success) {
            sessionManager.saveSession(
                result.accessToken,
                result.refreshToken,
                result.userId,
                result.email,
                result.login
            )
        }
        return result
    }

    suspend fun forgotPassword(loginIdentifier: String): AuthResult {
        return grpcService.forgotPassword(loginIdentifier)
    }

    suspend fun resetPassword(resetToken: String, code: String, newPassword: String): AuthResult {
        return grpcService.resetPassword(resetToken, code, newPassword)
    }

    suspend fun logout(refreshToken: String) {
        try {
            grpcService.logout(refreshToken)
        } finally {
            sessionManager.clearSession()
        }
    }
}