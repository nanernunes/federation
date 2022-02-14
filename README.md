# Federation
Forward Messages Between Broker Topics

## Declaring Backends (staticly) by Environment Variables
```ini
PROTOCOL_SINGLENAME_PROPERTY=value
```

### Declaring a static AMQP Backend
```ini
AMQP_MYORIGIN_HOST=localhost
AMQP_MYORIGIN_PORT=5672
AMQP_MYORIGIN_VHOST=/
AMQP_MYORIGIN_USER=guest
AMQP_MYORIGIN_PASS=guest
```

### Declaring a static SNS Backend
```ini
SNS_MYTARGET_REGION=us-east-1
SNS_MYTARGET_AWS_ACCESS_KEY_ID=XXXXXXXXXXX
SNS_MYTARGET_AWS_SECRET_ACCESS_KEY=ZZZZZZZZZZZZZZZZZZZZZZ
```

## Declaring Federations
This declaration pushes all messages from the topic (MYORIGIN) `helloworld` towards the topic at (MYTARGET) `helloworld`
```ini
FEDERATION_HELLO_WORLD=MYTARGET_helloworld,MYORIGIN_helloworld
```

## How it works
With this first version, we only allowed mapping from AMQP (RabbitMQ) clients to SNS brokers. For each federation mapping we create a new AMQP queue with the suffix `federation` that consumes the messages (with an exchange of the same name) then forward them to the target.

## Next releases
- We plan to allow mappings in both directions **AMQP->SNS** and **SNS->AMQP**, as well as the support of new protocols (eg.: Kafka, Redis, etc)
- A web interface and API with the immutable static method and the dynamic ones we could create in real-time (stored in a database)