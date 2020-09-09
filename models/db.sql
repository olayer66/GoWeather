-- Database: GoWeather

-- DROP DATABASE "GoWeather";

CREATE DATABASE "GoWeather"
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'Spanish_Spain.1252'
    LC_CTYPE = 'Spanish_Spain.1252'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

CREATE TABLE users (
	Id int NOT NULL,
	Name VARCHAR(50)NOT NULL,
	Surname VARCHAR(78)NOT NULL,
	Usuario VARCHAR(12)NOT NULL,
	Password VARCHAR(100)NOT NULL,
	ApiKey VARCHAR(100)NOT NULL
);
ALTER TABLE users ADD PRIMARY KEY (Id);

INSERT INTO users(Id,Name,Surname,Usuario,Password,ApiKey) VALUES
(1,'Paco','Porras','paco','1234','001'),
(2,'Ana','Armas','ana','1234','002'),
(3,'Maria','Piedra','maria','1234','003')