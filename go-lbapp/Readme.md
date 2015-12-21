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
create table Events(eventname VARCHAR(200), lat float, long float, creationtime time, creatorid VARCHAR(200));
```
