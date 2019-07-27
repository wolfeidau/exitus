CREATE TABLE issues (
    "id" uuid DEFAULT uuid_generate_v4(),
    "customer_id" uuid NOT NULL,
    "project_id" uuid NOT NULL,
    "reporter" uuid NOT NULL,   -- user identifier
    "assignee" uuid NULL,       -- user identifier
    "subject" citext NOT NULL,
    "content" text,
    "state" text NOT NULL,
    "severity" text,
    "category" text,
    "labels" text[] NOT NULL DEFAULT '{}'::text[],
    "created_at" timestamp with time zone DEFAULT now(),
    "updated_at" timestamp with time zone DEFAULT now(),
    PRIMARY KEY (id, customer_id, project_id)
    -- FOREIGN KEY (project_id, customer_id) REFERENCES projects (id, customer_id) ON DELETE RESTRICT
    -- FOREIGN KEY (reporter) REFERENCES users (id, customer_id) ON DELETE RESTRICT
    -- FOREIGN KEY (assignee) REFERENCES users (id, customer_id) ON DELETE RESTRICT
);
