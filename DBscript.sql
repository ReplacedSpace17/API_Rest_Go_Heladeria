-- Crear la base de datos Heladeria
CREATE DATABASE Heladeria;

-- Usar la base de datos Heladeria
USE Heladeria;

-- Crear la tabla sabores
CREATE TABLE sabores (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sabor VARCHAR(255),
    precio DECIMAL(10, 2)
);
