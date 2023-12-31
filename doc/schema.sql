-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-08-01T13:34:48.329Z

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" CITEXT UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "hashed_password" bytea NOT NULL,
  "activated" BOOLEAN NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT (0001-01-01 00:00:00),
  "version" bigint NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trains" (
  "id" bigserial PRIMARY KEY,
  "model_number" varchar NOT NULL,
  "name" varchar NOT NULL,
  "value" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "version" bigint NOT NULL DEFAULT 1,
  "last_edited_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "collection_trains" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "train_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wishlist_trains" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "train_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trade_offers" (
  "id" bigserial PRIMARY KEY,
  "offered_train" bigint NOT NULL,
  "offered_train_owner" bigint NOT NULL,
  "requested_train" bigint NOT NULL,
  "requested_train_owner" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trade_entry" (
  "id" bigserial PRIMARY KEY,
  "offered_train" bigint NOT NULL,
  "offered_train_owner" bigint NOT NULL,
  "requested_train" bigint NOT NULL,
  "requested_train_owner" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" bool NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "trains" ("model_number");

CREATE INDEX ON "trains" ("name");

CREATE INDEX ON "collection_trains" ("user_id");

CREATE INDEX ON "wishlist_trains" ("user_id");

CREATE INDEX ON "trade_offers" ("offered_train_owner");

CREATE INDEX ON "trade_offers" ("requested_train_owner");

CREATE INDEX ON "trade_entry" ("offered_train_owner");

CREATE INDEX ON "trade_entry" ("requested_train");

CREATE INDEX ON "trade_entry" ("offered_train");

CREATE INDEX ON "trade_entry" ("requested_train_owner");

CREATE INDEX ON "sessions" ("id");

ALTER TABLE "collection_trains" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "collection_trains" ADD FOREIGN KEY ("train_id") REFERENCES "trains" ("id");

ALTER TABLE "wishlist_trains" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wishlist_trains" ADD FOREIGN KEY ("train_id") REFERENCES "trains" ("id");

ALTER TABLE "trade_offers" ADD FOREIGN KEY ("offered_train") REFERENCES "trains" ("id");

ALTER TABLE "trade_offers" ADD FOREIGN KEY ("offered_train_owner") REFERENCES "users" ("id");

ALTER TABLE "trade_offers" ADD FOREIGN KEY ("requested_train") REFERENCES "trains" ("id");

ALTER TABLE "trade_offers" ADD FOREIGN KEY ("requested_train_owner") REFERENCES "users" ("id");

ALTER TABLE "trade_entry" ADD FOREIGN KEY ("offered_train") REFERENCES "trains" ("id");

ALTER TABLE "trade_entry" ADD FOREIGN KEY ("offered_train_owner") REFERENCES "users" ("id");

ALTER TABLE "trade_entry" ADD FOREIGN KEY ("requested_train") REFERENCES "trains" ("id");

ALTER TABLE "trade_entry" ADD FOREIGN KEY ("requested_train_owner") REFERENCES "users" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
