CREATE TABLE comments (
    "id" uuid DEFAULT uuid_generate_v4(),
    "customer_id" uuid NOT NULL,
    "project_id" uuid NOT NULL,
    "issue_id" uuid NOT NULL,
    "author" uuid NOT NULL,     -- user identifier
    "content" text,
    "created_at" timestamp with time zone DEFAULT now(),
    "updated_at" timestamp with time zone DEFAULT now(),
    PRIMARY KEY (id, customer_id, project_id, issue_id)
    -- FOREIGN KEY (project_id, customer_id, issue_id) REFERENCES issues (id, customer_id, project_id) ON DELETE RESTRICT
    -- FOREIGN KEY (author) REFERENCES users (id, customer_id) ON DELETE RESTRICT
);
