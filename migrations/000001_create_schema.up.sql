CREATE SCHEMA IF NOT EXISTS auth_schema;

CREATE DOMAIN auth_schema.uuid_key AS UUID
DEFAULT gen_random_uuid()
NOT NULL;

CREATE DOMAIN auth_schema.timestamp AS TIMESTAMP WITH TIME ZONE
DEFAULT (current_timestamp AT TIME ZONE 'UTC');