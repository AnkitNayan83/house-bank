-- Create "accounts" table
CREATE TABLE "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" VARCHAR NOT NULL,
    "balance" BIGINT NOT NULL,
    "currency" VARCHAR NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

-- Create "entries" table
CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" BIGINT NOT NULL,
    "amount" BIGINT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

-- Create "transfers" table
CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" BIGINT NOT NULL,
  "to_account_id" BIGINT NOT NULL,
  "amount" BIGINT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

-- Add Foreign Keys after creating all tables
ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

-- Create Indexes
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

-- Add comments
COMMENT ON COLUMN "entries"."amount" IS 'amount can be negative or positive';
COMMENT ON COLUMN "transfers"."amount" IS 'amount can only be positive';
