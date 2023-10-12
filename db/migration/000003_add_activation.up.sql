CREATE TABLE "activation_tokens" (
    "user_id" bigint NOT NULL UNIQUE,
    "activation_token" varchar NOT NULL, 
    "is_blocked" BOOLEAN NOT NULL DEFAULT false,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "activation_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
