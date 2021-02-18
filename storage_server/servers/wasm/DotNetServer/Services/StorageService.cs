using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using GrpcServer;
using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
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

        // this keeps truck of
        public bool instanceReady = false;

        // our wasm instance that gets made with lazy initialization 
        public Instance wasm;
        public Wasmtime.Externs.ExternMemory memory;
    }
    public class StorageService : Storage.StorageBase
    {
        private readonly ILogger<StorageService> _logger;

        private WasmSingleton wasmSingleton = WasmSingleton.Instance;


        public StorageService(ILogger<StorageService> logger)
        {

            if (!wasmSingleton.instanceReady)
            {
                _logger = logger;
                using var engine = new Engine();
                using var store = new Store(engine);


                WasiConfiguration wasiConfiguration = new WasiConfiguration();
                wasiConfiguration.WithPreopenedDirectory("./data", ".");

                using var module = Module.FromFile(engine, "../wasm_module/storage_application.wasm");

                using var host = new Host(store);
                host.DefineWasi("wasi_snapshot_preview1", wasiConfiguration);
                var instance = host.Instantiate(module);
                var memory = instance.Memories.Where(m => m.Name == "memory").First();
                ((dynamic)instance)._initialize();


                wasmSingleton.wasm = instance;
                wasmSingleton.memory = memory;
                wasmSingleton.instanceReady = true;

                Console.WriteLine("Instance ready");


            }

        }


        public override Task<ReadResponse> Read(ReadRequest request, ServerCallContext context)
        {

            var bytes = request.ToByteArray();

            var ptr = ((dynamic)wasmSingleton.wasm).new_alloc(bytes.Length);

            Console.WriteLine($"I got here and this is the ptr: {ptr}");

            bytes.CopyTo(wasmSingleton.memory.Span[ptr..]);

            var len = bytes.Length;


            var write = wasmSingleton.wasm.Functions.Where(f => f.Name == "read_data").First();

            //var resptr = ((dynamic)instance).store_data(ptr, len);

            var resptr = write.Invoke(ptr, len);

            ((dynamic)wasmSingleton.wasm).new_dealloc(ptr, len);

            var resultLen = ((dynamic)wasmSingleton.wasm).get_response_len();

            var result = new byte[resultLen];

            var memResult = wasmSingleton.memory.Span[resptr..(resptr + resultLen)];

            memResult.CopyTo(result);

            ReadResponse resMessage;

            resMessage = ReadResponse.Parser.ParseFrom(result);

            Console.WriteLine($"This is the value of message: {resMessage.Value}");
            Console.WriteLine($"This is the time of message: {resMessage.Timestamp}");

            ((dynamic)wasmSingleton.wasm).new_dealloc(resptr, resultLen);


            return Task.FromResult(resMessage);

        }

        public override Task<WriteResponse> Write(WriteRequest request, ServerCallContext context)
        {

            var bytes = request.ToByteArray();

            var ptr = ((dynamic)wasmSingleton.wasm).new_alloc(bytes.Length);

            Console.WriteLine($"I got here and this is the ptr: {ptr}");

            bytes.CopyTo(wasmSingleton.memory.Span[ptr..]);

            var len = bytes.Length;

            // Console.WriteLine("I got here");

            var write = wasmSingleton.wasm.Functions.Where(f => f.Name == "store_data").First();

            //var resptr = ((dynamic)instance).store_data(ptr, len);

            var resptr = write.Invoke(ptr, len);

            ((dynamic)wasmSingleton.wasm).new_dealloc(ptr, len);

            var resultLen = ((dynamic)wasmSingleton.wasm).get_response_len();

            var result = new byte[resultLen];

            var memResult = wasmSingleton.memory.Span[resptr..(resptr + resultLen)];

            memResult.CopyTo(result);

            WriteResponse resMessage;

            resMessage = WriteResponse.Parser.ParseFrom(result);

            Console.WriteLine($"This is the status of message: {resMessage.Ok}");

            ((dynamic)wasmSingleton.wasm).new_dealloc(resptr, resultLen);

            return Task.FromResult(resMessage);
        }
    }
}
