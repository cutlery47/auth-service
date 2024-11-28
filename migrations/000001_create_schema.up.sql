CREATE SCHEMA IF NOT EXISTS auth_schema;

CREATE DOMAIN auth_schema.uuid_key as UUID
DEFAULT gen_random_uuid()
NOT NULL;

