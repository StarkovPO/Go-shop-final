package store

const (
	createUserTable = `CREATE TABLE IF NOT EXISTS "users" (
        "primary_id" SERIAL PRIMARY KEY,
        "id" varchar(36) UNIQUE,
        "login" varchar(255) UNIQUE,
        "password_hash" varchar(255),
        "created_at" timestamp NOT NULL
    )`

	createOrderTable = `CREATE TABLE IF NOT EXISTS "orders" (
        "primary_id" SERIAL PRIMARY KEY,
        "user_id" varchar(36),
        "id" BIGSERIAL UNIQUE,
        "status" varchar(255) NOT NULL,
        "accrual" integer,
        "uploaded_at" timestamp NOT NULL
    )`

	createBalanceTable = `CREATE TABLE IF NOT EXISTS "balance" (
        "primary_id" SERIAL PRIMARY KEY,
        "user_id" varchar(36),
        "current" float
    )`

	createWithdrawTable = `CREATE TABLE IF NOT EXISTS "withdrawn" (
        "primary_id" SERIAL PRIMARY KEY,
        "order_id" integer,
        "withdrawn" float,
        "user_id" varchar(36),
        "processed_at" timestamp NOT NULL
    )`

	createOrderForeignKey     = `ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");`
	createBalanceForeignKey   = `ALTER TABLE "balance" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");`
	createWithdrawForeignKey  = `ALTER TABLE "withdrawn" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");`
	createWithdrawForeignKey2 = `ALTER TABLE "withdrawn" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");`

	createLoginIndex   = `CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex ON public.users (login)`
	createOrderIDIndex = `CREATE UNIQUE INDEX IF NOT EXISTS orders_id_uindex ON public.orders (id)`

	createUser = `
        INSERT INTO users (id, login, password_hash, created_at)
        VALUES ($1, $2, $3, to_timestamp($4))
    `
	checkLogin = `
        SELECT EXISTS (SELECT 1 FROM users WHERE login = $1 LIMIT 1)
    `

	getUserPass = `SELECT password_hash FROM users WHERE login = $1 LIMIT 1`

	createOrder = `
        INSERT INTO orders (user_id, id, status, accrual, uploaded_at)
        VALUES ($1, $2, $3, $4, to_timestamp($5))
    `

	getUserFromOrders = `SELECT user_id FROM orders WHERE id = $1 LIMIT 1`

	getOrders = `SELECT id, status, accrual, uploaded_at FROM orders WHERE user_id = $1`

	getUserID = `SELECT id FROM users WHERE login = $1`
)
