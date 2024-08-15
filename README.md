# NATS JetStram using Go

__NATS__ is defined as "a connective technology built for ever increasingly hyper-connected world" - https://nats.io/about/

* It is simple
* Fast
* Scalable
* Felixible - PUb/Sub, Request-Reply, POint-to-Point)
* Safe - support TLS/SSL and token based authentication
* Real-time Data Streaming - allow to immediate processing and analysis of message as they are published.
* Adaptability - Can be deployed on-promisse or cloud or integrated

NATS Advantages
* Low Latency
* Throughput
* Simplicity and Overhead

While both NATS and Kafka have their strengths and weaknesses, NATS offers advantages in simplicity, performance, operational overhead, and flexibility for certain use cases, particularly in real-time applications and microservices architectures. 

## Tecnologies
* https://github.com/nats-io/nats.go
* Go 1.23
* Docker

## How to use

### set up NATS server instance
* You should have Docker installed
```
make up
```

### publish a message 

```
make pub MSG="your message"
```

### consumer message
```
make sub 
```

### turn off NATS server
```
make down
```


