DROP TABLE IF EXISTS "wallets";
CREATE TABLE "wallets"(
    "user_id" SERIAL NOT NULL,
    "currency_id" TEXT NOT NULL,
    "currency_full_name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "public_key" TEXT NOT NULL,
    "private_key" TEXT NOT NULL,
    "last_balance" DOUBLE PRECISION NOT NULL,
    "last_balance_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    PRIMARY KEY(user_id,currency_id)
);

DROP TABLE IF EXISTS "transactions";
CREATE TABLE "transactions"(
    "tx_id" TEXT NOT NULL PRIMARY KEY,
    "room_id" TEXT,
    "currency_id" TEXT NOT NULL,
    "currency_full_name" TEXT NOT NULL,
    "from_address" TEXT NOT NULL,
    "to_address" TEXT NOT NULL,
    "transaction_at" TIMESTAMPTZ NOT NULL
);