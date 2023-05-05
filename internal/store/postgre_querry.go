package store

const (
	createUserTable = `CREATE TABLE IF NOT EXISTS "users" (
        "primary_id" integer PRIMARY KEY,
        "id" varchar(36) UNIQUE,
        "login" varchar(255) UNIQUE,
        "password_hash" varchar(60),
        "created_at" timestamp NOT NULL
    )`

	createOrderTable = `CREATE TABLE IF NOT EXISTS "orders" (
        "primary_id" integer PRIMARY KEY,
        "user_id" varchar(36),
        "id" integer UNIQUE,
        "status" varchar(255) NOT NULL,
        "accrual" integer,
        "uploaded_at" timestamp NOT NULL
    )`

	createBalanceTable = `CREATE TABLE IF NOT EXISTS "balance" (
        "primary_id" integer PRIMARY KEY,
        "user_id" varchar(36),
        "current" float
    )`

	createWithdrawTable = `CREATE TABLE IF NOT EXISTS "withdrawn" (
        "primary_id" integer PRIMARY KEY,
        "order_id" integer,
        "withdrawn" float,
        "user_id" varchar(36),
        "processed_at" timestamp NOT NULL
    )`

	createOrderForeignKey     = `ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");`
	createBalanceForeignKey   = `ALTER TABLE "balance" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");`
	createWithdrawForeignKey  = `ALTER TABLE "withdrawn" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");`
	createWithdrawForeignKey2 = `ALTER TABLE "withdrawn" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");`
)
