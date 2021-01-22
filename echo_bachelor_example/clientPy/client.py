import grpc

# import the generated classes
import echo_pb2
import echo_pb2_grpc

# open gRPC channel
channel = grpc.insecure_channel("localhost:50051")

# create a stub (client)
stub = echo_pb2_grpc.EchoStub(channel)

# take the user input
print("Exit/exit to exit this program\n")
while True:
    user_input = str(input("What are you thinking? "))

    if user_input == "Exit" or user_input == "exit":
        break
    elif user_input == "":
        continue

    print(user_input)
    message = echo_pb2.Message(content=user_input)

    response = stub.Send(message)

    print("Message Sent!\n")

    print("Recived this from the server:" + str(response))
