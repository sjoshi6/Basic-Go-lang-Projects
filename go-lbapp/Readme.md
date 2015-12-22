### Request
```
curl -X POST -d '{"userid":"sau123", "password":"test", "name":"sau"}' http://localhost:8000/v1/signup
```

### Creating a postgresql DB
```
initdb db_lbapp
```
### Start Postgres DB
```
postgres -D db_lbapp
```

### Initialize a DB inside the main postgresql DB
```
createdb db_lbapp
```

### Login to DB using PSQL
```
psql db_lbapp
```
### Creating a table to store passwords
```
create table userlogin (UserID VARCHAR(100), Password VARCHAR(200), Name VARCHAR(200));
```

### Creating a table to store new events
```
create table Events(id SERIAL, eventname VARCHAR(200), lat float, long float, creationtime timestamp, creatorid VARCHAR(200));
```

### Curl command to insert create event
```
curl -X POST -d '{"eventname":"house warming party", "latitude":"100.2", "longitude":"127.3", "creationtime": "2015-12-15 07:36:25", "creatorid":"sjoshi6" }' http://localhost:8000/v1/create_event
```
