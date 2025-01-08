CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "username" varchar,
  "user_image_uri" varchar,
  "company_name" varchar,
  "company_image_uri" varchar
);