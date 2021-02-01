import grpc

# import the generated classes
import echo_pb2
import echo_pb2_grpc

# open gRPC channel
channel = grpc.insecure_channel("152.94.1.100:50051")

# create a stub (client)
stub = echo_pb2_grpc.EchoStub(channel)

# take the user input
print("Exit/exit to exit this program\n")
while True:
    user_input = str(input("Message to send:  "))

    if user_input == "Exit" or user_input == "exit":
        break
    elif user_input == "":
        continue

    message = echo_pb2.echo_message(content=user_input)

    print("\n - Message Sent!")

    response = stub.Send(message)

    print(f" - Recived this from the server: '{response.content}'\n")
