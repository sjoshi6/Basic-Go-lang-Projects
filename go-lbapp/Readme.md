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
create table Events(id SERIAL, event_name VARCHAR(200), lat float, lng float, creation_time timestamp, creator_id VARCHAR(200), start_time timestamp, end_time timestamp, max_mem int, min_mem int, friend_only boolean, gender CHAR(1), min_age int, max_age int);
```

### To create postgres user
```
db_lbapp=# CREATE USER postgres SUPERUSER;
```

### Curl command to insert create event
```
curl -X POST -d '{"eventname":"house warming party", "latitude":"100.20", "longitude":"111.70","creatorid":"apatwar", "start_time": "2015-12-30 10:00:00", "end_time":"2015-12-31 09:00:00", "max_mem":"30", "min_mem":"10", "friend_only":"True", "gender":"N", "min_age": "22", "max_age":"24"  }' http://localhost:8000/v1/create_event

curl -X POST -d '{"eventname":"chess ", "latitude":"100.8", "longitude":"111.20","creatorid":"sjoshi6", "start_time": "2015-12-30 10:00:00", "end_time":"2015-12-31 09:00:00", "max_mem":"30", "min_mem":"10", "friend_only":"False", "gender":"M", "min_age": "22", "max_age":"24"  }' http://localhost:8000/v1/create_event

curl -X POST -d '{"eventname":"Monopoly at my place", "latitude":"100.80", "longitude":"111.40","creatorid":"sdrangne", "start_time": "2015-12-31 23:00:00", "end_time":"2015-12-31 23:45:00", "max_mem":"4", "min_mem":"3", "friend_only":"False", "gender":"N", "min_age": "24", "max_age":"30"  }' http://localhost:8000/v1/create_event
```

### Search for NearbyEvents with Radius
```
curl -X GET -d '{"latitude":"100.8", "longitude":"111.2", "radius": "100"}' http://localhost:8000/v1/search_events
```

### Starting two golang API Servers on two ports
```
./go-lbapp 8000
./go-lbapp 8001
```
