package store

const (
	createUserTable = `CREATE TABLE IF NOT EXISTS "Users" (
        "primary_id" integer PRIMARY KEY,
        "id" varchar(36) UNIQUE,
        "login" varchar(255) UNIQUE,
        "password_hash" varchar(60),
        "created_at" timestamp NOT NULL
    )`

	createOrderTable = `CREATE TABLE IF NOT EXISTS "Orders" (
        "primary_id" integer PRIMARY KEY,
        "user_id" varchar(36),
        "id" integer UNIQUE,
        "status" varchar(255) NOT NULL,
        "accrual" integer,
        "uploaded_at" timestamp NOT NULL
    )`

	createBalanceTable = `CREATE TABLE IF NOT EXISTS "Balance" (
        "primary_id" integer PRIMARY KEY,
        "user_id" varchar(36),
        "current" float
    )`

	createWithdrawTable = `CREATE TABLE IF NOT EXISTS "Withdrawn" (
        "primary_id" integer PRIMARY KEY,
        "order_id" integer,
        "withdrawn" float,
        "user_id" varchar(36),
        "processed_at" timestamp NOT NULL
    )`

	createOrderForeignKey     = `ALTER TABLE "Orders" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");`
	createBalanceForeignKey   = `ALTER TABLE "Balance" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");`
	createWithdrawForeignKey  = `ALTER TABLE "Withdrawn" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");`
	createWithdrawForeignKey2 = `ALTER TABLE "Withdrawn" ADD FOREIGN KEY ("order_id") REFERENCES "Orders" ("id");`
)
