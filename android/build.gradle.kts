// android/build.gradle.kts
plugins {
    alias(libs.plugins.compose.compiler) apply false
    id("com.android.application") version "9.1.1" apply false
    id("com.google.protobuf") version "0.10.0" apply false
}