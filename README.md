# MQTT Publisher/Subscriber


## Publisher
The publisher code publishes temperature, humidity, and light data at specific intervals using the MQTT protocol. A Publisher object is created for each data type, and messages are sent to the MQTT server by calling the publishMessage function with these objects. The newTLSConfig function is also defined to establish a secure connection using TLS.

## Subscriber
This code listens to temperature, humidity, and light data using the MQTT protocol. The SubscriberClient struct contains subscription information and topics. The subscribe function subscribes to specific topics and prints the messages to the console. The newTLSConfig function is also defined to establish a secure connection using TLS.

