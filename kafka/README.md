# go-eliona Kafka 
The go-eliona Kafka package provides handy methods to produce and read messages in Kafka topics. 

## Installation
To use the log package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/kafka"
```

Optionally you can define an environment variable named BROKERS which sets the Kafka bootstrap servers.
Not defined the packages takes the `kafka:29092` which is the default in eliona environment.

```bash
export BROKERS=10.10.100.1:29092,192.168.178.1:9092 # This is optionally, default is kafka:29092
```

## Usage

After installation, you can produce messages in Kafka topics.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/kafka"
```

For example, you can push a temperature object to the `climate` topic. You have to design a new
Producer and produce the temperature to the topic.

```go
type Temperature struct {
    Value int
    Unit  string
}
```

```go
producer := kafka.NewProducer()
defer producer.Close()
temperature := Temperature{Value: 24, Unit: "Celsius"}
_ = kafka.Produce(producer, "climate", temperature)
```

To read messages from a topic have to define a new consumer and subscribe a topic. After this, you can
read the temperatures through a channel.

```go
consumer := kafka.NewConsumer()
defer consumer.Close()
kafka.Subscribe(consumer, "climate")
temperatures := make(chan Temperature)
go kafka.Read(consumer, temperatures)

for temperature := range temperatures {
    fmt.Printf("Temperature is: %d %s", temperature.Value, temperature.Unit)
}
```
