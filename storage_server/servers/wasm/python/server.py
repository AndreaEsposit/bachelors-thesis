from concurrent import futures
from pathlib import Path
import wasmtime
import grpc
import time
import threading

# Import the generated classes
import storage_pb2_grpc
import storage_pb2

functionsToImp = ["store_data", "read_data"]
wasmLocation = "../wasm_module/storage_application.wasm"
preOpenedFile = "./data"

# utilized to store exported Wasm functions
exportedWasmFunctions = {}

# gRPC related variables
grpc_host = u'152.94.162.16'  # 152.94.162.16 bbchain6
grpc_port = u'50051'
grpc_address = u'{host}:{port}'.format(host=grpc_host, port=grpc_port)

# lock
lock = threading.Lock()

# WasmInstantiate instatiates a Wasm module given a .wasm file location and a list
# of the functions that need to be exported


def WasmInstantiate(functions, wasmLocation, preOpenedDir):
    # Wasmtime Embedding
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
    exportedWasmFunctions["alloc"] = instance_linking.exports["new_alloc"]
    exportedWasmFunctions["dealloc"] = instance_linking.exports["new_dealloc"]
    exportedWasmFunctions["get_len"] = instance_linking.exports["get_response_len"]
    exportedWasmFunctions["memory"] = instance_linking.exports["memory"]

    for name in functions:
        exportedWasmFunctions[name] = instance_linking.exports[name]

    print("INSTANTIATED")


# copy_to_memory handles the copy of serialized data to the
# Wasm's memory
def copy_to_memory(sdata: bytearray):
    # allocate memory in wasm
    ptr = exportedWasmFunctions["alloc"](len(sdata))

    for i, v in enumerate(sdata):
        exportedWasmFunctions["memory"].data_ptr[ptr + i] = v

    return ptr


# call_wasm handles the actual wasm function calls, and takes care of all calls to alloc/dialloc in the wasm instance
def call_wasm(func, request, return_message):
    bytes_as_string = request.SerializeToString()

    lock.acquire()  # take lock
    ptr = copy_to_memory(bytes_as_string)
    length = len(bytes_as_string)

    result_ptr = exportedWasmFunctions[func](ptr, length)
    res_ptr_int = int(result_ptr)

    # deallocate request protobuf message
    exportedWasmFunctions["dealloc"](ptr, length)

    result_len = exportedWasmFunctions["get_len"]()
    # res_len_int = int(result_len)

    response = bytearray(result_len)

    for i in range(result_len):
        response[i] = exportedWasmFunctions["memory"].data_ptr[result_ptr+i]

    # deallocate response protobuf message
    exportedWasmFunctions["dealloc"](result_ptr, result_len)

    lock.release()  # release lock

    # parse response to a protobuf message
    return_message.ParseFromString(response)

    return return_message


# create a class to define the server functions, derived
# from storage_pb2_grpc.StorageServicer
class StorageServicer(storage_pb2_grpc.StorageServicer):

    def Read(self, request, context):
        return_message = call_wasm(
            "read_data", request, storage_pb2.ReadResponse())
        return return_message

    def Write(self, request, context):
        return_message = call_wasm(
            "store_data", request, storage_pb2.WriteResponse())
        return return_message


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
    WasmInstantiate(functionsToImp, wasmLocation, preOpenedFile)
    Server.run()
