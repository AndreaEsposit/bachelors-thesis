from concurrent import futures
from pathlib import Path
import tempfile
import wasmtime
import grpc, time
import numpy as np 
import google.protobuf.message as proto

# Import the generated classes
import storage_pb2_grpc, storage_pb2

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
def copy_memory(sdata):
    ptr = alloc(np.int32(len(sdata)))

    # cast pointer to int32
    ptr32 = np.int32(ptr)

    for i, v in enumerate(sdata):
        mem.data_ptr[ptr32 + np.int32(i)] = v

    return ptr32

# call_function handles all the calls the desired
def call_function(fn, m):
    received_bytes = proto.Message.SerializeToString(m)

    ptr = copy_memory(received_bytes)
    length = np.int32(len(received_bytes))

    res_ptr = fn(ptr, length)
    res_ptr32 = np.int32(res_ptr)
    
    # deallocate request protobuf message
    dealloc(ptr, length)

    result_len = get_len()
    int_res_len = np.int32(result_len)

    response = bytearray(int_res_len)
    for i in response:
        response[i] = mem.data_ptr[res_ptr32, int_res_len]

    # deallocate response protobuf message
    dealloc(res_ptr32, int_res_len)
    return response



# create a class to define the server functions, derived
# from storage_pb2_grpc.StorageServicer
class StorageServicer(storage_pb2_grpc.StorageServicer):

    def Read(self, request, context):
        return_message = storage_pb2.ReadResponse(message='' % request.Filename)
        return return_message

    def Write(self, request, context):
        return_message = storage_pb2.WriteResponse(message='' % request.Filename)
        return return_message


class Server:
    # Initialize gRPC server
    @staticmethod
    def run():
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        storage_pb2_grpc.add_StorageServicer_to_server(StorageServicer(), server)
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


