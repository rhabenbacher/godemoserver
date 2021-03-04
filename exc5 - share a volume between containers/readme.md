# Share a volume between containers
```yml
services:
  time-server:
    image: alpine
    command: sh -c "while true;do echo $$(date)>>/tmp/time/time.txt;echo 'stored new time';sleep 10;done"
    volumes:
      - "timedata:/tmp/time"
  time-reader:
    # If the time-server does not find the file it would terminate
    # This can be avoided with the use of depends_on
    depends_on:
    - time-server
    image: alpine
    command: sh -c "tail -f /time/time.txt"
    volumes:
      - "timedata:/time"
volumes:
  timedata:   
```