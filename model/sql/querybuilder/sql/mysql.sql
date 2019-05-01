CREATE TABLE testtable1(
    id VARCHAR(255) not null,
    body text not null,    
    PRIMARY KEY(id)
) DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ENGINE=InnoDB;

CREATE TABLE testtable2(
    id2 VARCHAR(255) not null,
    body2 text not null,    
    PRIMARY KEY(id2)
) DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ENGINE=InnoDB;
