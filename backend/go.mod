module github.com/Enziofael/eatfit/backend

go 1.25.6

require (

	// Валидация email адресов
	github.com/badoux/checkmail v1.2.4

	// Redis клиент для кеширования и blacklist токенов
	github.com/go-redis/redis/v8 v8.11.5

	// JWT токены для аутентификации
	github.com/golang-jwt/jwt/v5 v5.3.1

	// Миграции базы данных
	github.com/golang-migrate/migrate/v4 v4.19.1

	// Работа с UUID
	github.com/google/uuid v1.6.0

	// PostgreSQL connection pool для pgx
	github.com/jackc/pgx/v4 v4.18.3

	// Драйвер PostgreSQL для Go (быстрый и надёжный)
	github.com/jackc/pgx/v5 v5.9.2

	// Загрузка конфигурации из .env файлов
	github.com/joho/godotenv v1.5.1

	// Отправка email через SMTP
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible

	// Генерация моков для тестирования
	github.com/stretchr/testify v1.11.1

	// Логирование (структурированное)
	go.uber.org/zap v1.28.0

	// Хеширование паролей (bcrypt)
	golang.org/x/crypto v0.51.0
	// gRPC фреймворк для высокопроизводительного RPC взаимодействия
	google.golang.org/grpc v1.81.1

	// Генерация Go кода из proto файлов
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
)
