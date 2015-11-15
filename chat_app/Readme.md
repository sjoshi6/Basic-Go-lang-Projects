# A Chat Application written in pure Go

## Tech Details
Redis
Go
Go Channels
Go Routines

## Implementation Details
- Redis is used as the backend with redigo connector and redisurl for easy connectivity.
- Uses redis pubsub connections to subscribe on a redis "messages" queue
- On receiving the message the , go pushes the message onto a go channel
- Each user has 3 connections to redis
    1. A connection that manages user session start / end / active heartbeat
    2. A connection to publish messages from command line to "messages" queue on redis
    3. A connection to subscribe to "messages" queue
- A go function runs independently to fetch text from the subscribe GoChannel
- When a text is received the function classifies it as a command or a plain text
- Heartbeat sender is an independent go lang process that keeps connection alive

## Running Commands

Start Redis Server
```
redis-server
```

Go to the chat_app folder
```
cd chat_app
```

Compile the go code and fetch dependancies
```
go install
```

Move to the bin folder and run the compiled code
```
./chat_app <username>
```
