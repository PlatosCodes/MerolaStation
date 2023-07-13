CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trains" (
  "id" bigserial PRIMARY KEY,
  "model_number" varchar NOT NULL,
  "name" varchar NOT NULL,
  "value" bigint NOT NULL DEFAULT 0,
  "version" bigint NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "collection_trains" (
  "user_id" bigint NOT NULL,
  "train_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wishlist_trains" (
  "user_id" bigint NOT NULL,
  "train_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "trains" ("model_number");

CREATE INDEX ON "trains" ("name");

CREATE INDEX ON "collection_trains" ("user_id");

CREATE INDEX ON "wishlist_trains" ("user_id");

ALTER TABLE "collection_trains" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "collection_trains" ADD FOREIGN KEY ("train_id") REFERENCES "trains" ("id");

ALTER TABLE "wishlist_trains" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wishlist_trains" ADD FOREIGN KEY ("train_id") REFERENCES "trains" ("id");
