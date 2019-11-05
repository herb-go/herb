CREATE TABLE columns(
    id SERIAL,
    f_nullable INTEGER,
    f_bigint BIGINT not null,
    f_bigserial BIGSERIAL not null,
    f_boolean BOOLEAN not null,
    f_bytea BYTEA not null,
    f_char CHAR(10) not null,
    f_varchar VARCHAR (10) not null,
    f_float8 FLOAT8 not null,
    f_integer INTEGER not null,
    f_json JSON not null,
    f_jsonb JSONB not null,
    f_real REAL not null,
    f_text TEXT not null,
    f_timestamptz timestamptz not null,
    f_timestamp Timestamp not null,
    f_smallint smallint not null,
    f_smallserial SMALLSERIAL not null,
    f_uuid UUID not null,
    f_xml XML not null,
    PRIMARY KEY(id)
) 