[![Go](https://github.com/jneethling/wikimoningest/actions/workflows/go.yml/badge.svg)](https://github.com/jneethling/wikimoningest/actions/workflows/go.yml)
# wikimoningest
A Golang Kafka ingester for the Hatnote Wikipedia monitor (https://github.com/hatnote/wikimon).

The purpose of this service is to act as a simple illustration of the following functions:
* Websocket client
* Externally imported, schema-based Apache Avro data serialisation
* Apache Kafka producer

## Todo
* Improve test coverage (in particular for the producer)
* Illustrate fan-out pattern with a Async producer