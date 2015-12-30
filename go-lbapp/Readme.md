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
### Creating a table to store user details
```
create table Users(UserId VARCHAR(100) PRIMARY KEY,Password VARCHAR(200) NOT NULL, FirstName VARCHAR(200) NOT NULL, LastName VARCHAR(200) NOT NULL, Gender CHAR(1) NOT NULL, Age Int NOT NULL, PhoneNumber VARCHAR(20) NOT NULL);
```
### POST signup data & login
```
curl -X POST -d '{"userid":"sjoshi6", "password":"xxxxxx", "firstname": "xxxx", "lastname":"xxxx", "gender": "M", "age": "xx", "phonenumber":"xxxxxxxxxx"}' http://localhost:8000/v1/signup

curl -X POST -d '{"userid":"sjoshi6", "password":"saurabh8391"}' http://localhost:8000/v1/login
```

### Creating a table to store new events
```
create table Events(id SERIAL PRIMARY KEY, event_name VARCHAR(200) NOT NULL, lat float NOT NULL, lng float NOT NULL, creation_time timestamp NOT NULL, creator_id VARCHAR(200) NOT NULL, start_time timestamp NOT NULL, end_time timestamp, max_mem int, min_mem int, friend_only boolean NOT NULL, gender CHAR(1) NOT NULL, min_age int, max_age int, FOREIGN KEY (creator_id) REFERENCES Users);
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
