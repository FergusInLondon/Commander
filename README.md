


## Manual Testing

    $ go build
    $ ./Commander -d

### Bridge TCP connections to the Unix Socket (via Socat)

Socat is [pretty easy to use](https://coderwall.com/p/c3wyzq/forwarding-tcp-traffic-to-a-unix-socket) as well, and with one command you can treat the Unix Socket like any other TCP based URL.

    $ sudo apt-get install socat
    $ socat -d -d TCP4-LISTEN:8080,fork UNIX-CONNECT:/tmp/commander.sock
      2017/11/02 01:21:39 socat[9455] N listening on AF=2 0.0.0.0:8080

### Hit the endpoints

Go via Curl, or your favourite API tool - i.e Postman.

    $ curl localhost:8080/listing
      [{"name":"Echo","Description":"Echos back a message from the daemon"}]