DROP TABLE IF EXISTS sessionid;

CREATE TABLE sessionid (
    name VARCHAR(255) NOT NULL PRIMARY KEY,
    sid VARCHAR(255) NOT NULL 
);

INSERT INTO sessionid 
    (name, sid)
VALUES
    ("John Doe", "007");    