# gRPC/Kafka ImgToPDF Compose

This project is an **example** realization of three microservices..

### Main features
- gRPC for services communication 
- bidi-streaming gRPC protobuf
- Chi router
- Kafka to communicate with email notification service
- docker-compose to deploy client/server instances
***
## How it works:
Controller takes images uploaded via web UI and sends it to the gRPC server as complex data stream.
Server converts data stream of images into one .pdf file and sends it back in data stream.
Therefore Client service publishes message event in kafka while notification service reads it and makes notifications.



