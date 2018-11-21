# go-backpressure
Fun with TCP Flow Control


## Purpose

You can measure TCP Flow Control ("back pressure") functionality by starting a websocket server and pushing more data
to it than it can read at a time.

[TCP Backpressure](https://www.brianstorti.com/tcp-flow-control/) can be a tricky thing because it only applies to a 
connection between two hosts while in modern cloud applications, the client usually never talks to the application directly.
In-between there are loadbalancers, routers, proxies and more loadbalancers. 
Depending on the type of loadbalancer in front, the "slow down" message might never reach the client.

This app is meant for various testing scenarios such as

* local deployments
* single VM public IP deployments
* VM behind loadbalancer(s)
* [CloudFoundry App](https://www.cloudfoundry.org/) deployments


## How to use

The app supports both standard TCP and [websocket](https://www.html5rocks.com/en/tutorials/websockets/basics/) connections
using the [gorilla websocket implementation](https://github.com/gorilla/websocket).

* To start the server:
    * Set `protocol: tcp` or `protocol: websocket` in server/config.yml
    * Run command
    ```
    go build
    ./go-backpressure -server server/config.yml
    ```

* To start the client:
    * Set `protocol: tcp` or `protocol: websocket` in client/config.yml
    * Run command
    ```
    go build
    ./go-backpressure -client client/config.yml
    ```

* To use the webclient:
    * The websocket server also provides a minimal webclient
    * Set `protocol: websocket` in server/config.yml
    * Point your browser to [http://localhost:8080/client](http://localhost:8080/client)
    * To see console output, open developer tools
    
    
## Capturing Traffic

To observe when the server advertises a zero window, using [Wireshark](https://www.wireshark.org/#download) is recommended.
When running locally, use these settings:
 * Capture on Loopback lo0
 * Set filter to `tcp.port == 8080 and not websocket` to only see pure TCP traffic.
 * Open one TCP PDU (=protocol data unit, i.e. packet / datagram / segment)
    * Righ-click on "Calculated Window Size", select "Apply as column"
    
When running remote, you will likely put TLS in front using HAproxy or NGINX. Notice that since the traffic
will then be encrypted, Wireshark can no longer detect websocket traffic. All traffic will show as pure TCP and TLS.

Advertised window sizes will still be visible though.


## Deploying on Docker

The project provides a simple Dockerfile to allow you to deploy it on any server that runs docker.

To build, simply run

```
docker build -t go-backpressure .
```

To deploy the server on a remote machine, run

```
docker run -d -p 8080:8080 go-backpressure
```

## Deploying on CloudFoundry

The project also provides a minimal manifest for CF deployment.

After logging in to your org and space, simply run

```
cf push
```