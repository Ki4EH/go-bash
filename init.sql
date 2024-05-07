CREATE TABLE IF NOT EXISTS "commands" (
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "script" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "output" TEXT NOT NULL
);