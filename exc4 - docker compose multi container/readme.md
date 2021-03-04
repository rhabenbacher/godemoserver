# Run two containers with docker-compose

## docker-compose.yml
```yml
services:
  frontend:
    image: goserver:0.1
    # expose the frontend to host on port 8080
    ports: 
      - "8080:8000"
     # start goserver with command frontend    
    command: ["./goserver","frontend"]
    # environment variables to connect to api server
    environment:
      - API_HOST=apiserver
      - API_PORT=3000
    depends_on:
      - apiserver   
  apiserver:
    image: goserver:0.1
    # expose the port internally on port 3000
    expose:
      - "3000"
    # start goserver with command rest  
    command: ["./goserver","rest"]
```