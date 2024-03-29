import grpc
import time
import json
import os
import threading
from pathlib import Path
from concurrent import futures

import storage_pb2_grpc
import storage_pb2

# gRPC related variables
grpc_host = u'152.94.162.17'  # bbchain7
grpc_port = u'50051'
grpc_address = u'{host}:{port}'.format(host=grpc_host, port=grpc_port)


lock = threading.Lock()

# create a class to define the server functions, derived
# from storage_pb2_grpc.StorageServicer
class StorageServicer(storage_pb2_grpc.StorageServicer):

    def Read(self, request, context):
        filename = request.FileName

        path = Path(__file__).parent / "./data"
        filepath = str(path) + "/" + filename + ".json"

        if os.path.isfile(filepath):
            # take lock
            lock.acquire()
            with open(filepath) as f:
                data = json.load(f)

            val = data["value"]
            timestamp = {
                "seconds": data["seconds"],
                "nanos": data["nseconds"]
            }

            # release lock
            lock.release()

            return storage_pb2.ReadResponse(Value=val, Timestamp=timestamp, Ok=1)
        else:
            timestamp = {
                "seconds": 0,
                "nanos": 0
            }

            return storage_pb2.ReadResponse(Value="", Timestamp=timestamp, Ok=0)

    def Write(self, request, context):
        filename = request.FileName
        timestamp = request.Timestamp
        val = request.Value

        path = Path(__file__).parent / "./data"

        filepath = str(path)+"/" + filename + ".json"

        dataset = {
            "seconds": timestamp.seconds,
            "nseconds": timestamp.nanos,
            "value": val
        }

        # take lock
        lock.acquire()

        # write to JSON pretty
        with open(filepath, "w") as write_file:
            json.dump(dataset, write_file, indent=4)
        lock.release()  # release lock

        return storage_pb2.WriteResponse(Ok=1)


class Server:
    # Initialize gRPC server
    @ staticmethod
    def run():
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        storage_pb2_grpc.add_StorageServicer_to_server(
            StorageServicer(), server)
        server.add_insecure_port(grpc_address)
        server.start()
        print("Server is running at: " + grpc_address)

        # instead of server.wait_for_termination()
        # since server.start() will not block,
        # a sleep-loop is added to keep alive
        try:
            while True:
                time.sleep(86400)
        except KeyboardInterrupt:
            server.stop(0)


if __name__ == "__main__":
    Server.run()
