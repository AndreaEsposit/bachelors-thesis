using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using GrpcServer;
using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Threading;
using Wasmtime;

namespace DotNetServer.Services
{

    // Singleton pattern to create only one instance of the Wasm Instance
    public sealed class WasmSingleton
    {
        private static WasmSingleton _instance; // tells us where we actaully saved the the object

        private WasmSingleton() { }
        public static WasmSingleton Instance

        {
            get
            {
                if (_instance == null)
                {
                    _instance = new WasmSingleton();

                }

                return _instance;
            }
        }

        // this keeps truck of the state of the instace
        public bool instanceReady = false;

        // mutex lock to control wasm access
        private readonly Mutex mu = new Mutex();

        // our wasm instance that gets made with lazy initialization 
        public Instance wasm;
        public Wasmtime.Externs.ExternMemory memory;

        // keeps track of the functions used by the applicaion (no alloc, dealloc etc..)
        public Dictionary<String, Wasmtime.Externs.ExternFunction> funcs;

        public T callWasm<T>(String fn, IMessage message) where T : IMessage<T>, new()
        {
            // --- Copy the buffer to the module's memory 
            var bytes = message.ToByteArray();

            mu.WaitOne(); // acquire lock

            var ptr = ((dynamic)wasm).new_alloc(bytes.Length);
            var len = bytes.Length;

            //Console.WriteLine($"I got here and this is the ptr: {ptr}");

            bytes.CopyTo(memory.Span[ptr..]);

            // you can technically define the function each and every time
            //var func = wasm.Functions.Where(f => f.Name == fn).First();

            // call the function you are intrested in
            var resPtr = funcs[fn].Invoke(ptr, len);

            // deallocate request protobuf message
            ((dynamic)wasm).new_dealloc(ptr, len);

            // get Wasm response lenght
            var resultLen = ((dynamic)wasm).get_response_len();

            // copy byte array fom WebAssembly to a result byte array
            var resultMemory = memory.Span[resPtr..(resPtr + resultLen)];
            var result = new byte[resultLen];
            resultMemory.CopyTo(result);

            // deallocate repsonse protobud message form Wasmtime's memory
            ((dynamic)wasm).new_dealloc(resPtr, resultLen);

            mu.ReleaseMutex();

            message = parse<T>(result);
            return message is T value ? value : default(T);
        }

        // generic parser
        private IMessage parse<T>(byte[] buf) where T : IMessage<T>, new()
        {
            MessageParser<T> parser = new MessageParser<T>(() => new T());
            return parser.ParseFrom(buf);
        }
    }
    public class StorageService : Storage.StorageBase
    {
        //private readonly ILogger<StorageService> _logger;
        private WasmSingleton wasmSingleton = WasmSingleton.Instance;


        public StorageService() //ILogger<StorageService> logger
        {

            if (!wasmSingleton.instanceReady)
            {
                // Set up the logger
                //_logger = logger;
                // initialize the WebAssembly module
                using var engine = new Engine();
                using var store = new Store(engine);


                // pass access to data directory to this Wasm module
                WasiConfiguration wasiConfiguration = new WasiConfiguration();
                wasiConfiguration.WithPreopenedDirectory("./data", ".");

                // Create the WebAssembly-module
                using var module = Module.FromFile(engine, "../wasm_module/storage_application.wasm");

                // get instance of the wasm module ans specify the version of wasi
                using var host = new Host(store);
                host.DefineWasi("wasi_snapshot_preview1", wasiConfiguration);
                var instance = host.Instantiate(module);

                // execure the _initialize function to given wasm access to the data folder
                ((dynamic)instance)._initialize();

                // export the only memory that we are using with this module 
                var memory = instance.Memories.Where(m => m.Name == "memory").First();

                wasmSingleton.wasm = instance;
                wasmSingleton.memory = memory;

                // export application functions from the WebAssembly module
                var funcs = new Dictionary<String, Wasmtime.Externs.ExternFunction>();
                funcs["write"] = wasmSingleton.wasm.Functions.Where(f => f.Name == "store_data").First();
                funcs["read"] = wasmSingleton.wasm.Functions.Where(f => f.Name == "read_data").First();


                wasmSingleton.funcs = funcs;
                wasmSingleton.instanceReady = true;

                Console.WriteLine("Instance ready");

            }

        }


        public override Task<ReadResponse> Read(ReadRequest request, ServerCallContext context)
        {
            var result = wasmSingleton.callWasm<ReadResponse>("read", request);
            // Console.WriteLine($"This is the value of message: {resMessage.Value}");
            // Console.WriteLine($"This is the time of message: {resMessage.Timestamp}");
            return Task.FromResult(result);
        }

        public override Task<WriteResponse> Write(WriteRequest request, ServerCallContext context)
        {
            var result = wasmSingleton.callWasm<WriteResponse>("write", request);
            //Console.WriteLine($"This is the status of message: {resMessage.Ok}");
            return Task.FromResult(result);
        }
    }
}
