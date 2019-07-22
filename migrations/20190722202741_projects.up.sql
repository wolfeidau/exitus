CREATE TABLE projects (
    "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" citext NOT NULL UNIQUE,
    "description" text,
    "owner_id" uuid NOT NULL,
    "labels" text[] NOT NULL DEFAULT '{}'::text[],
    "created_at" timestamp with time zone DEFAULT now(),
    "updated_at" timestamp with time zone DEFAULT now()
);
