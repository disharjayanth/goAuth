DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY ,
    name VARCHAR(255) NOT NULL UNIQUE ,
    password VARCHAR(255) NOT NULL 
);

INSERT INTO users
    (name, password)
VALUES 
    ("John Doe", "1234");