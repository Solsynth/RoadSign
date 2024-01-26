import asyncio
import websockets

async def handle_websocket(websocket, path):
    # This function will be called whenever a new WebSocket connection is established

    # Send a welcome message to the client
    await websocket.send("Welcome to the WebSocket server!")

    try:
        # Enter the main loop to handle incoming messages
        async for message in websocket:
            # Print the received message
            print(f"Received message: {message}")

            # Send a response back to the client
            response = f"Server received: {message}"
            await websocket.send(response)
    except websockets.exceptions.ConnectionClosedError:
        print("Connection closed by the client.")

# Create the WebSocket server
start_server = websockets.serve(handle_websocket, "localhost", 8765)

print("WebSocket server started at ws://localhost:8765")

# Run the server indefinitely
asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
