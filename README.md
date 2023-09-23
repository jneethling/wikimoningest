[![Go](https://github.com/jneethling/wikimoningest/actions/workflows/go.yml/badge.svg)](https://github.com/jneethling/wikimoningest/actions/workflows/go.yml)
# wikimoningest
A Golang Kafka ingester for the Hatnote Wikipedia monitor (https://github.com/hatnote/wikimon).

The purpose of this service is to act as a simple illustration of the following functions:
* Websocket client
* Apache Avro data serialisation (using an imported codec with a predefined schema)
* Apache Kafka producer

Please note for this toy project the schema supporting the codec is very permissive - in a real-world use case it would probably be much more restrictive, which would allow tighter control and more robust test cases.

## Todo
* Illustrate fan-out pattern with a Async producer.