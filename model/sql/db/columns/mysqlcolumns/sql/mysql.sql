CREATE TABLE columns(
    id INTEGER  AUTO_INCREMENT,
    f_nullable INTEGER,
    f_tinyint TINYINT not null,
    f_bit BIT not null,
    f_bool BOOL not null,
    f_smallint SMALLINT not null,
    f_mediumint MEDIUMINT not null,
    f_int int  not null,
    f_integer INTEGER not null,
    f_biginteger BIGINT not null,
    f_float FLOAT not null,
    f_double DOUBLE not null,
    f_doubleprecision DOUBLE PRECISION not null,
    f_datetime DATETIME not null,
    f_timestamp TIMESTAMP not null,
    f_char CHAR not null,
    f_varchar VARCHAR(255) not null,
    f_tinytext TINYTEXT not null,
    f_text TEXT not null,
    f_mediumtext MEDIUMTEXT not null,
    f_longtext LONGTEXT not null,
    f_binary BINARY not null,
    f_varbinary VARBINARY(255) not null,
    f_tinyblob TINYBLOB not null,
    f_blob BLOB not null,
    f_mediumblob MEDIUMBLOB not null,
    f_longblob LONGBLOB not null,
    PRIMARY KEY(id)
) 