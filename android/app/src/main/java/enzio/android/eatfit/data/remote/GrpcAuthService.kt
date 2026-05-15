package enzio.android.eatfit.data.remote

import enzio.android.eatfit.data.model.AuthResult
import enzio.android.eatfit.data.model.LoginRequest
import enzio.android.eatfit.data.model.RegisterRequest
import enzio.android.eatfit.proto.eatfit.v1.AuthServiceGrpc
import enzio.android.eatfit.proto.eatfit.v1.ForgotPasswordRequest
import enzio.android.eatfit.proto.eatfit.v1.LogoutRequest
import enzio.android.eatfit.proto.eatfit.v1.ResetPasswordRequest
import enzio.android.eatfit.proto.eatfit.v1.VerifyEmailRequest
import io.grpc.ManagedChannel
import io.grpc.okhttp.OkHttpChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.util.concurrent.TimeUnit

class GrpcAuthService(
    private val serverUrl: String = "10.0.2.2",
    private val port: Int = 50051
) {

    private val channel: ManagedChannel = OkHttpChannelBuilder
        .forAddress(serverUrl, port)
        .usePlaintext()
        .keepAliveTime(30, TimeUnit.SECONDS)
        .build()

    private val stub = AuthServiceGrpc.newBlockingStub(channel)

    suspend fun register(request: RegisterRequest): AuthResult = withContext(Dispatchers.IO) {
        try {
            val protoRequest = eatfit.v1.RegisterRequest.newBuilder()
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
            AuthResult(success = false, message = e.message ?: "Connection error")
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
            AuthResult(success = false, message = e.message ?: "Error")
        }
    }

    suspend fun login(request: LoginRequest): AuthResult = withContext(Dispatchers.IO) {
        try {
            val builder = eatfit.v1.LoginRequest.newBuilder()
                .setPassword(request.password)
                .setDeviceInfo("Android")

            if (request.loginIdentifier.contains('@')) {
                builder.email = request.loginIdentifier
            } else {
                builder.login = request.loginIdentifier
            }

            val response = stub.login(builder.build())

            AuthResult(
                success = true,
                accessToken = response.accessToken,
                refreshToken = response.refreshToken,
                userId = response.user.userId,
                email = response.user.email,
                login = response.user.login
            )
        } catch (e: Exception) {
            AuthResult(success = false, message = e.message ?: "Invalid credentials")
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
            AuthResult(success = false, message = e.message ?: "Error")
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
            AuthResult(success = false, message = e.message ?: "Error")
        }
    }

    suspend fun logout(refreshToken: String) {
        try {
            val request = LogoutRequest.newBuilder()
                .setRefreshToken(refreshToken)
                .build()
            stub.logout(request)
        } catch (_: Exception) { }
    }
}