# gRPC ImgToPDF Compose Service

This project is an **experiment** realization.

### Main features
- gRPC as transport protocol
- bidi-streaming gRPC protobuf
- Chi router
- docker-compose to deploy client/server instances
***
## How it works:
Controller takes images uploaded via web UI and sends it to the gRPC server as complex data stream.
Server converts data stream of images into one .pdf file and sends it back in data stream.



