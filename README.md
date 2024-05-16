## Description: WebSocket Proxy for Server Instance Communication and Broadcasting

This project showcases the concept of a WebSocket proxy system designed to facilitate communication among multiple server instances. The primary WebSocket server listens for messages from any connected server instance and subsequently broadcasts these messages to all other connected server instances. This ensures seamless and real-time data synchronization across all server instances, enabling efficient inter-server communication and consistent data dissemination.

## Key Features:

- Main WebSocket Server: Acts as the central hub for receiving messages from any server instance.
- Broadcast Mechanism: Automatically broadcasts received messages to all connected server instances.


## Getting started
1. Start the main websocket server
   
   ```go run main.go websocket```
2. Start the server instance 1
   
   ```go run main.go server1```
3. Start the server instance 2
   
   ```go run main.go server2```
4. Call either one of the server instances /post API to broadcast message to all other connected server instances.


## Notes
Alternative solution is to use [NATS](https://nats.io/) which is much better and flexible.