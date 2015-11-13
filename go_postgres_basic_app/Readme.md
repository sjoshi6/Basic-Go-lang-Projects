## Basic Postgres connector for go-lang

1. Instal postgres
2. Start Postgres server
```
postgres -D /usr/local/var/postgres
```

3. Go-lang has a package called database/sql which is a generic interface for all sql db's.
4. Download and instance a db specific driver for each db type used (MySQL, postgres) etc.
