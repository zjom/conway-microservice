# Conway Game of Life Microservice

This is a microservice that computes the next iteration of conway's game of life.

## Usage

### Running the server
To run the server, run:
```bash
go run main.go
```

To check that the server is running, you can send a GET request to `http://localhost:8080/health`.

### Protocol
This endpoint receives a JSON object with the following structure:
```json
{
    "board": [
        [0, 1, 0],
        [0, 1, 0],
        [0, 1, 0]
    ]
}
```

The board is a 2D array of integers, where 0 represents a dead cell and 1 represents a live cell.

The response will be a JSON object with the following structure:
```json
{
    "board": [
        [0, 0, 0],
        [1, 1, 1],
        [0, 0, 0]
    ]
}
```

The server will respond with errors if:
  1. the board is empty
  2. the rows are not of the same length
  3. a cell is not 0 or 1

