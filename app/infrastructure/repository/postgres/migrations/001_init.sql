-- +migrate Up
CREATE TABLE "users" (
  "user_id" UUID PRIMARY KEY,
  "view" varchar(20) NOT NULL,
  "profile" jsonb NOT NULL DEFAULT '{}',
  "version" int NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT 'NOW()',
  "updated_at" timestamptz NOT NULL DEFAULT 'NOW()'
);

CREATE TABLE "idents" (
  "user_id" UUID NOT NULL,
  "ident" varchar(4096) NOT NULL,
  "ident_confirmed" boolean NOT NULL DEFAULT false,
  "kind" int NOT NULL DEFAULT 1,
  "password" varchar(1024),
  "version" int NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT 'NOW()',
  "updated_at" timestamptz NOT NULL DEFAULT 'NOW()',
  PRIMARY KEY ("ident", "kind")
);

ALTER TABLE "idents" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

CREATE INDEX "ident_idx" ON "idents" ("ident");

-- +migrate Down

DROP TABLE idents;
DROP TABLE profiles;