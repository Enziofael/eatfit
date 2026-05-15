package enzio.android.eatfit.data.remote

import enzio.android.eatfit.data.model.AuthResult
import enzio.android.eatfit.data.model.LoginRequest
import enzio.android.eatfit.data.model.RegisterRequest
import enzio.android.eatfit.proto.AuthServiceGrpc
import enzio.android.eatfit.proto.ForgotPasswordRequest
import enzio.android.eatfit.proto.LogoutRequest
import enzio.android.eatfit.proto.RegisterRequest as ProtoRegisterRequest
import enzio.android.eatfit.proto.ResetPasswordRequest
import enzio.android.eatfit.proto.VerifyEmailRequest
import enzio.android.eatfit.proto.LoginRequest as ProtoLoginRequest
import enzio.android.eatfit.proto.ChangeLoginRequest
import enzio.android.eatfit.proto.ChangePasswordRequest
import io.grpc.ManagedChannel
import io.grpc.okhttp.OkHttpChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.util.concurrent.TimeUnit


class GrpcAuthService(
    private val serverUrl: String = "192.168.0.10",
    private val port: Int = 50051
) {

    private val channel: ManagedChannel = OkHttpChannelBuilder
        .forAddress(serverUrl, port)
        .usePlaintext()
        .keepAliveTime(30, TimeUnit.SECONDS)
        .keepAliveTimeout(10, TimeUnit.SECONDS)
        .keepAliveWithoutCalls(true)
        .build()

    private val stub = AuthServiceGrpc.newBlockingStub(channel)

    suspend fun register(request: RegisterRequest): AuthResult = withContext(Dispatchers.IO) {
        try {
            val protoRequest = ProtoRegisterRequest.newBuilder()
                .setEmail(request.email)
                .setLogin(request.login)
                .setPassword(request.password)
                .setPasswordConfirmation(request.passwordConfirmation)
                .build()

            val response = stub.register(protoRequest)

            AuthResult(
                success = response.success,
                message = response.message,
                userId = response.userId
            )
        } catch (e: Exception) {
            AuthResult(success = false, message = "Error: ${e.message}")
        }
    }

    suspend fun verifyEmail(userId: String, code: String): AuthResult = withContext(Dispatchers.IO) {
        try {
            val request = VerifyEmailRequest.newBuilder()
                .setUserId(userId)
                .setVerificationCode(code)
                .build()

            val response = stub.verifyEmail(request)

            AuthResult(
                success = response.success,
                message = response.message
            )
        } catch (e: Exception) {
            AuthResult(success = false, message = "Error: ${e.message}")
        }
    }

    suspend fun login(request: LoginRequest): AuthResult = withContext(Dispatchers.IO) {
        try {
            val builder = ProtoLoginRequest.newBuilder()
                .setPassword(request.password)
                .setDeviceInfo("Android Emulator")

            if (request.loginIdentifier.contains('@')) {
                builder.email = request.loginIdentifier
            } else {
                builder.login = request.loginIdentifier
            }

            val response = stub
                .withDeadlineAfter(10, TimeUnit.SECONDS)
                .login(builder.build())

            AuthResult(
                success = true,
                accessToken = response.accessToken,
                refreshToken = response.refreshToken,
                userId = response.user.userId,
                email = response.user.email,
                login = response.user.login
            )
        } catch (e: Exception) {
            AuthResult(success = false, message = "Login failed: ${e.message}")
        }
    }

    suspend fun forgotPassword(loginIdentifier: String): AuthResult = withContext(Dispatchers.IO) {
        try {
            val request = ForgotPasswordRequest.newBuilder()
                .setLoginIdentifier(loginIdentifier)
                .build()

            val response = stub.forgotPassword(request)

            AuthResult(
                success = response.success,
                message = response.message,
                resetToken = response.resetToken
            )
        } catch (e: Exception) {
            AuthResult(success = false, message = "Error: ${e.message}")
        }
    }

    suspend fun resetPassword(resetToken: String, code: String, newPassword: String): AuthResult = withContext(Dispatchers.IO) {
        try {
            val request = ResetPasswordRequest.newBuilder()
                .setResetToken(resetToken)
                .setVerificationCode(code)
                .setNewPassword(newPassword)
                .setPasswordConfirmation(newPassword)
                .build()

            val response = stub.resetPassword(request)

            AuthResult(
                success = response.success,
                message = response.message
            )
        } catch (e: Exception) {
            AuthResult(success = false, message = "Error: ${e.message}")
        }
    }

    suspend fun logout(refreshToken: String) {
        withContext(Dispatchers.IO) {
            try {
                val request = LogoutRequest.newBuilder()
                    .setRefreshToken(refreshToken)
                    .build()
                stub
                    .withDeadlineAfter(5, java.util.concurrent.TimeUnit.SECONDS)
                    .logout(request)
            } catch (_: Exception) {
                // Игнорируем ошибки
            }
        }
    }

    suspend fun changeLogin(userId: String, newLogin: String, password: String): AuthResult = withContext(Dispatchers.IO) {
        try {
            val request = ChangeLoginRequest.newBuilder()
                .setUserId(userId)
                .setNewLogin(newLogin)
                .setPassword(password)
                .build()
            val response = stub.changeLogin(request)
            AuthResult(success = response.success, message = response.message)
        } catch (e: Exception) {
            AuthResult(success = false, message = e.message ?: "Error")
        }
    }

    suspend fun changePassword(userId: String, currentPassword: String, newPassword: String, confirmPassword: String): AuthResult = withContext(Dispatchers.IO) {
        try {
            val request = ChangePasswordRequest.newBuilder()
                .setUserId(userId)
                .setCurrentPassword(currentPassword)
                .setNewPassword(newPassword)
                .setPasswordConfirmation(confirmPassword)
                .build()
            val response = stub.changePassword(request)
            AuthResult(success = response.success, message = response.message)
        } catch (e: Exception) {
            AuthResult(success = false, message = e.message ?: "Error")
        }
    }
}