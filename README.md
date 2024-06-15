# SETUP PRODUCER_CONSUMER SCRIPT
- set up your .env file in the root folder using .env.sample format
- Have `make` installed on your pc or use the command to install make on your macbook `brew install make`
- run the command `make docker-up` to setup the docker images on docker and run the containers in detached mode.
- run `go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db -database="${DB_DSN}" up` to create the users table on the postgres database with DB_DSN as your connection string in the format `postgresql://dbuser\:password@localhost\:5432/dbname\?sslmode=disable`

## Start Consumer
- run `make start-consumer` to start up consumer

## Start Producer
- Edit `data.json` to suit the data structure below:
```
    {
        "action": "create",
        "data": {
            "name": "Luka Modric",
            "age": 39
        }
    }
```
- run `make start-producer` to start up producer
