# Price tracker 

## Componentes

### Base de datos

```bash
docker run --name=mysql1 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql/mysql-server:8.0

docker exec -it mysql1 mysql -uroot -p123456

create database pricetracker;
```

### Servidor principal

```bash
cd ./backend/servidorPrincipal
go run .
```


### Servicio extracción de datos

```bash
cd ./backend/servicioExtraccionDatos
python3 dataScrapper.py
```

