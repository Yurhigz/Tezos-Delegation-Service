DROP DATABASE IF EXISTS tzktdb;
CREATE DATABASE tzktdb;

\c tzktdb

-- Initialisation de la table delegations 
DROP TABLE IF EXISTS delegations CASCADE;
CREATE TABLE delegations (
    id SERIAL PRIMARY KEY,
    delegator VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    amount BIGINT NOT NULL,
    level BIGINT NOT NULL
);