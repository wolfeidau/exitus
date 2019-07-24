CREATE TABLE customers (
    "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" citext NOT NULL UNIQUE,
    "description" text,
    "labels" text[] NOT NULL DEFAULT '{}'::text[],
    "created_at" timestamp with time zone DEFAULT now(),
    "updated_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE projects (
    "id" uuid DEFAULT uuid_generate_v4(),
    "customer_id" uuid NOT NULL,
    "name" citext NOT NULL,
    "description" text,
    "labels" text[] NOT NULL DEFAULT '{}'::text[],
    "created_at" timestamp with time zone DEFAULT now(),
    "updated_at" timestamp with time zone DEFAULT now(),
    PRIMARY KEY (id, customer_id),
    UNIQUE ("customer_id", "name")
    -- CONSTRAINT customer_projects FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE RESTRICT
);
