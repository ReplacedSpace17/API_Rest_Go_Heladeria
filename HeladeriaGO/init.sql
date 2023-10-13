USE Heladeria;

-- Crear la tabla sabores
CREATE TABLE sabores (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sabor VARCHAR(255),
    precio DECIMAL(10, 2)
);