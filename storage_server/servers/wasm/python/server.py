from concurrent import futures
from pathlib import Path
import wasmtime
import grpc
import time
import numpy as np
from google.protobuf.timestamp_pb2 import Timestamp

# Import the generated classes
import storage_pb2_grpc
import storage_pb2

store = wasmtime.Store()


linker = wasmtime.Linker(store)
wasi_config = wasmtime.WasiConfig()
path = Path(__file__).parent / "./data"
wasi_config.preopen_dir(str(path), ".")
wasi = wasmtime.WasiInstance(store, "wasi_snapshot_preview1", wasi_config)
linker.define_wasi(wasi)

# Load and compile the WebAssembly-module
path = Path(__file__).parent / "../wasm_module/storage_application.wasm"
module_linking = wasmtime.Module.from_file(store.engine, path)

# Instantiate the module which only uses WASI
instance_linking = linker.instantiate(module_linking)

# Execute the _initialize function to give WASM access
# to the data folder
init = instance_linking.exports["_initialize"]
init()

# Export functions and memory from the WebAssembly module
alloc = instance_linking.exports["new_alloc"]
dealloc = instance_linking.exports["new_dealloc"]
get_len = instance_linking.exports["get_response_len"]
write = instance_linking.exports["store_data"]
read = instance_linking.exports["read_data"]
mem = instance_linking.exports["memory"]

# gRPC related variables
grpc_host = u'[::]'
grpc_port = u'50051'
grpc_address = u'{host}:{port}'.format(host=grpc_host, port=grpc_port)


# copy_mem handles the copy of serialized data to the
# Wasm's memory
def copy_memory(sdata: bytearray):
    print("in copy memory")

    print("managed to convert, this is the num: " + str((len(sdata))))
    ptr = alloc(len(sdata))
    print("allocated")

    # cast pointer to int32
    ptr32 = int(ptr)

    for i, v in enumerate(sdata):
        mem.data_ptr[ptr32 + i] = v

    return ptr32


# call_function handles all the calls the desired
def call_function(fn, bytes_as_string):
    # serialize message (not working)

    print("I came this far1")

    ptr = copy_memory(bytes_as_string)
    length = len(bytes_as_string)
    print("I came this far2")

    print("I got this ptr: " + str(ptr))
    res_ptr = fn(ptr, length)
    print("called fn")
    res_ptr32 = int(res_ptr)
    print("I came this far4")

    # deallocate request protobuf message
    dealloc(ptr, length)
    print("I came this far5")

    result_len = get_len()
    int_res_len = int(result_len)
    print("I came this far6")

    response = bytearray(int_res_len)

    print("We made the response: " + str(response))
    print("This is the ptr type " + str(type(res_ptr32)) +
          " this is the len type: " + str(type(int_res_len)))
    for i in range(int_res_len):
        response[i] = mem.data_ptr[res_ptr32+i]

    print("I came this far7")

    # deallocate response protobuf message
    dealloc(res_ptr32, int_res_len)
    print("I came this far8")

    #s = listToString(response)
    return response


# create a class to define the server functions, derived
# from storage_pb2_grpc.StorageServicer
class StorageServicer(storage_pb2_grpc.StorageServicer):

    def Read(self, request, context):

        b = request.SerializeToString()
        print("this is b:" + str(b) + " and this is the size: " + str(len(b)))
        response = call_function(read, b)

        return_message = storage_pb2.ReadResponse()

        return_message.ParseFromString(response)
        return return_message

    def Write(self, request, context):

        b = request.SerializeToString()
        print("this is b:" + str(b) + " and this is the size: " + str(len(b)))
        response = call_function(write, b)

        # time = Timestamp()
        # time.GetCurrentTime()
        return_message = storage_pb2.WriteResponse()

        return_message.ParseFromString(response)
        return return_message


# def listToString(s):
#     strl = ""
#     for ele in s:
#         strl += str(ele)
#     return strl


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

        # instead of server.wait_for_termination
        # since server.start() will not block,
        # a sleep-loop is added to keep alive
        try:
            while True:
                time.sleep(86400)
        except KeyboardInterrupt:
            server.stop(0)


if __name__ == "__main__":
    Server.run()
