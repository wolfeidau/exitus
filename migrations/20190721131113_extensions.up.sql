-- The citext module provides a case-insensitive character string type, citext. Essentially, it
-- internally calls lower when comparing values. Otherwise, it behaves almost exactly like text.
-- https://www.postgresql.org/docs/10/citext.html
CREATE EXTENSION IF NOT EXISTS "citext";

-- This module implements the hstore data type for storing sets of key/value pairs within a single
-- PostgreSQL value. This can be useful in various scenarios, such as rows with many attributes that
-- are rarely examined, or semi-structured data. Keys and values are simply text strings.
-- https://www.postgresql.org/docs/10/hstore.html
CREATE EXTENSION IF NOT EXISTS "hstore";

-- The pg_trgm module provides functions and operators for determining the similarity of alphanumeric
-- text based on trigram matching, as well as index operator classes that support fast searching for
-- similar strings.
-- https://www.postgresql.org/docs/10/pgtrgm.html
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- The uuid-ossp module provides functions to generate universally unique identifiers (UUIDs) using
-- one of several standard algorithms. There are also functions to produce certain special UUID constants.
-- https://www.postgresql.org/docs/10/uuid-ossp.html
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
