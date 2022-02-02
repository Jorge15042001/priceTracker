# Price tracker 

## Componentes

### Base de datos mySQL usando docker


```bash
docker run --name=mysql1 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql/mysql-server:8.0

docker exec -it mysql1 mysql -uroot -p123456

create database pricetracker;
```

### Servidor principal programado en GO

Se ejecuta en el puerto 3000

```bash
cd ./backend/servidorPrincipal
go run .
```


### Servicio extracción de datos programado en python

Se ejecuta en el puerto 5000

```bash
cd ./backend/servicioExtraccionDatos
python3 dataScrapper.py
```

### Correr web app

Entrar a [localhost:3000](http://localhost:3000) desde cualquier navegador
