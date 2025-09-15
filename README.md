# Tezos Delegation Service

## Objectives

This project is a Golang service that retrieves new [delegations](https://opentezos.com/node-baking/baking/delegating/) made on the Tezos protocol and exposes them through a public REST API.  


## How it works

1. The service **polls the TzKT API** every 10 seconds to fetch new delegations.  
   - If up to date, it switches to a 15-minute polling interval.  

2. All retrieved data is **stored in a PostgreSQL database** for persistence.  

3. A public **HTTP API** exposes the stored delegations through:  
```

GET /xtz/delegations

```

4. The endpoint supports **two optional query parameters**:  
- `year` → year in format `YYYY`  
- `level` → block height  

---

## API Example

**Request:**
```

GET http://localhost:3000/xtz/delegations?timestamp=2018&level=251000

````

**Response:**
```json
{
  "Delegations": [
    {
      "timestamp": "2018-MM-ddThh:mm:ssZ",
      "amount": x,
      "delegator": <randomDelegatorAddress>,
      "level": 251000
    }
  ]
}
```

---

## Setup & Installation

### Requirements

* Go 1.21+
* Docker & Docker Compose

### Run locally

1. Unzip the folder:
2. Get inside the main folder : 
    ```bash
    cd ./kiln-projects 
    ```
3. Build the container : 
    ```bash
    docker-compose up -d 
    ```
4. Install the go packages :
    ```bash
    go mod tidy
    ```
    
5. Run the service:

   ```bash
   go run main.go
   ```
6. API available at:

   ```
   http://localhost:3000/xtz/delegations
   ```


## PostgreSQL Tips

### 1. Connect to the database from Docker

Open a shell inside the PostgreSQL container:

```bash
docker-compose exec -it -u postgres postgres /bin/bash
```

Then connect to PostgreSQL:

```bash
psql -U postgres
```

### 2. Useful `psql` commands

* Connect to tzkt database:

  ```sql
  \c tzktdb
  ```

* Show delegations structure:

  ```sql
  \d delegations
  ```

* Get informations about delegations:

  ```sql
  SELECT * FROM delegations ORDER BY timestamp DESC LIMIT 100;
  ```

---

## Testing

Unit and integration tests are provided to validate the main components:

- Polling: valid JSON, invalid JSON, and API unreachable cases.
- Database: insertion and retrieval of delegations (with a test database).
- REST API: endpoint tested with query parameters.

Run tests with:
```bash
go test ./tests -v

```

--- 

## Possible Improvements

* Replace polling with **TzKT WebSocket API** for real-time data.
* Add a first request to the DB to get the latest level in case of application restart to prevent from reidexing all the DB again.
* Add another query parameter such as `limit` to allow users to decide how many documents they want.
* Segment more the project by moving the `StartPoller`function in another folder.
* Add a Circuit Breaker to prevent API shortage.
* Rate limiting and retry logic could be considered to reduce the impact on TZKT API and improve resilience.
