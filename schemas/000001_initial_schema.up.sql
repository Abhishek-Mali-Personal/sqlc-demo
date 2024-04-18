CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS lookups
(
    lookup_id uuid DEFAULT uuid_generate_v4(),
    table_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    display_order integer NOT NULL,
    display_text character varying(255) COLLATE pg_catalog."default" NOT NULL,
    is_active boolean NOT NULL,
    parent_id uuid,
    internal_key character varying(100) COLLATE pg_catalog."default" NOT NULL,
    concurrency_key character varying(10) COLLATE pg_catalog."default" NOT NULL,
    create_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    create_user_id integer NOT NULL,
    update_date date,
    update_user_id integer,
    value_text character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT lookups_pkey PRIMARY KEY (lookup_id)
    );