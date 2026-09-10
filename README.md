# Subglutee Project - User Service

REST + gRPC service สำหรับจัดการ user profile

- REST `:8082` - frontend/gateway (`/api/v1/users/me`)
- gRPC `:50052` - auth-service (`FindOrCreateUser` ตอน Google login)
- Postgres - เก็บ user

### Structure

```
routes/         map path -> controller
controllers/    เชื่อม HTTP กับ DTO, validate ข้อมูล แล้วเรียก usecase
usecases/       business logic
repositories/   db query
grpc/           gRPC controllers
dtos/           request/response struct
models/         db struct
config/         env loader + db pool
migrations/     sql migration
```

REST: `routes -> controllers -> usecases -> repositories -> Postgres`
gRPC: `grpc/ -> usecases -> repositories` (usecase ตัวเดียวกับ REST)

### Prerequisite

- Go 1.26
- Docker + Docker Compose
- `make` - `winget install ezwinports.make`
- tools:

```terminal
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/evilmartians/lefthook@latest
```

### Setup

```terminal
git clone https://github.com/polar-bear-cu/sgt-user-service.git
cd sgt-user-service
lefthook install
cp .env.example .env
go mod download
make compose-up
make migrate-up
```

### Useful Commands

Check `Makefile`

### Migrations

golang-migrate, ไฟล์คู่ `up`/`down` ใน `migrations/`

```terminal
make migrate-create name=add_something
make migrate-up
make migrate-down
```

แก้ schema = migration ใหม่เสมอ ห้ามแก้ไฟล์ที่ merge ไปแล้ว

### Dev tools

- pgweb: `localhost:8083` - ดู local db

#### ถ้าแก้ proto พร้อม service นี้

```terminal
cd ..
go work init ./sgt-proto ./sgt-user-service
```

หรือถ้ามี go.work แล้ว...

```terminal
go work use ./sgt-proto ./sgt-user-service
```
