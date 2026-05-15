// android/app/build.gradle.kts
plugins {
    id("com.android.application")
    id("com.google.protobuf")
    alias(libs.plugins.compose.compiler)
}

android {
    namespace = "enzio.android.eatfit"
    compileSdk = 37

    defaultConfig {
        applicationId = "enzio.android.eatfit"
        minSdk = 30
        targetSdk = 34
        versionCode = 1
        versionName = "1.0"
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_21
        targetCompatibility = JavaVersion.VERSION_21
    }



    buildFeatures {
        compose = true
    }

    composeOptions {
    }
}

dependencies {
    // Compose BOM
    implementation(platform("androidx.compose:compose-bom:2026.05.00"))
    implementation("com.google.android.material:material:1.14.0")
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-graphics")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    implementation("androidx.compose.material:material-icons-extended")

    // Activity & Lifecycle
    implementation("androidx.activity:activity-compose:1.13.0")
    implementation("androidx.lifecycle:lifecycle-runtime-compose:2.10.0")
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.10.0")

    // Navigation
    implementation("androidx.navigation:navigation-compose:2.9.8")

    // gRPC
    implementation("io.grpc:grpc-okhttp:1.81.0")
    implementation("io.grpc:grpc-protobuf-lite:1.81.0")
    implementation("io.grpc:grpc-stub:1.81.0")
    implementation("com.google.protobuf:protobuf-javalite:4.34.1")
    compileOnly("javax.annotation:javax.annotation-api:1.3.2")

    // Coroutines
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.11.0")

    // DataStore
    implementation("androidx.datastore:datastore-preferences:1.2.1")
}

protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:3.25.3"
    }
    plugins {
        create("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:1.61.0"
        }
    }
    generateProtoTasks {
        all().forEach { task ->
            task.plugins {
                create("grpc") {
                    option("lite")
                }
            }
            task.builtins {
                create("java") {
                    option("lite")
                }
            }
        }
    }
}