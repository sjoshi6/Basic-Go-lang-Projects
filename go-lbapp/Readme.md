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
