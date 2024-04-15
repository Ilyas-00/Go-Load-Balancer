# Go Load Balancer :computer::arrows_counterclockwise:

This project implements a load balancer in Go, allowing requests to be distributed among multiple servers in a balanced manner.

## Architecture :building_construction:

The project consists of several main components:

### 1. Simple Server (`SimpleServer`) :gear:

- **Features**:
  - Manages a single backend server. 
  - Checks the server's health status.

- **Methods**:
  - `Address() string`: Returns the server's address.
  - `IsAlive() bool`: Checks if the server is online.
  - `Serve(rw http.ResponseWriter, req *http.Request)`: Serves HTTP requests.

### 2. Load Balancer (`LoadBalancer`) :balance_scale:

- **Features**:
  - Distributes requests among multiple backend servers.
  - Utilizes a round-robin algorithm for load distribution.

- **Methods**:
  - `getNextAvailableServer() Server`: Returns the next available server.
  - `ServeProxy(rw http.ResponseWriter, req *http.Request)`: Serves requests by redirecting them to backend servers.

## Algorithms :bulb:

- **Round-Robin**:
  - Uses a round-robin load balancing algorithm to select backend servers fairly. Each request is directed to a different server in turn.

## Usage :rocket:

1. Clone the project from the GitHub repository:
 ```bash
   git clone https://github.com/Ilyas-00/Go-Load-Balancer.git
 ```
  -  Run the main program:
  ```bash
    go run main.go
  ```
  -  The load balancer will start listening for requests on the specified port.
