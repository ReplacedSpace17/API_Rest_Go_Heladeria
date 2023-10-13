# Utiliza una imagen base de Golang
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el código fuente de tu aplicación al contenedor
COPY . .

# Instala el controlador de MySQL para Go
RUN go get -u github.com/go-sql-driver/mysql

# Compila tu aplicación Go
RUN go build -o main

# Expon el puerto en el que se ejecutará tu aplicación
EXPOSE 8080

# Ejecuta tu aplicación Go al inicio del contenedor
CMD ["/app/main"]
