#### Overall approach:

1. Create a new Redis client
2. Use the Ping() method to check if the Redis server is up and running
3. Create a new Fiber application
4. Create a new task using the POST method and push it to the Redis list
5. Get the first task and so on as first in first out using the GET method as a JSON response

#### How it works:

1. The POST /tasks route creates a new task and pushes it to the Redis list using the LPUSH command
2. The GET /tasks route pops the first item from the Resdis list as FIFO using the RPOP command

#### How to run:

go run main.go

#### What optimizations could be made?

1. Use of middleware to set authentication and authorization for the API
2. Handling multiple tasks at the same time using go routines
