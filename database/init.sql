DROP DATABASE IF EXISTS tzktdb;
CREATE DATABASE tzktdb;

\c tzktdb

-- Initialisation de la table delegations 
DROP TABLE IF EXISTS delegations CASCADE;
CREATE TABLE delegations (
    id SERIAL PRIMARY KEY,
    adress VARCHAR(50) NOT NULL UNIQUE,
    timestamp TIMESTAMP NOT NULL
    amount BIGINT NOT NULL
    blockheight BIGINT NOT NULL
);