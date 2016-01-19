# go-keystore 

A key value store implementation in golang. Used to store only JSON values

## Postgres commands

### Create a postgres db
```
initdb db_lbapp
```

### Start the db
```
postgres -D db_lbapp
```

### Create a DB instead of postgres for storage node
```
createdb storagenode
```

### Register a new node with loadbalancer redis
```
curl -X POST -d '{"ip_address": "192.186.177.112"}' http://localhost:8000/v1/register
```

### Request for the next node's IP ready to serve request
```
curl -X GET http://localhost:8000/v1/nextNode
```


