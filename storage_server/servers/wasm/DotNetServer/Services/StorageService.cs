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
    public sealed class StorageSingleton
    {
        private static StorageSingleton _instance; //tells us where we actaully saved the the object

        private StorageSingleton() { }
        public static StorageSingleton Instance

        {
            get
            {
                if (_instance == null)
                {
                    _instance = new StorageSingleton();

                }

                return _instance;
            }
        }

        public bool fatto = false;
        public dynamic wasm;
    }
    public class StorageService : Storage.StorageBase
    {
        private readonly ILogger<StorageService> _logger;

        private StorageSingleton storageSingleton = StorageSingleton.Instance;


        public StorageService(ILogger<StorageService> logger)
        {

            if (!storageSingleton.fatto)
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


                storageSingleton.wasm = instance;
                storageSingleton.fatto = true;

                Console.WriteLine("Instance ready");


                // write test 

                // var message = new WriteRequest();

                // message.Value = "This is incredibly top secret data!";
                // message.Timestamp = DateTime.UtcNow.ToTimestamp();
                // message.FileName = "Sparta";

                // var bytes = message.ToByteArray();

                // var ptr = ((dynamic)instance).new_alloc(bytes.Length);

                // Console.WriteLine($"I got here and this is the ptr: {ptr}");

                // bytes.CopyTo(memory.Span[ptr..]);

                // var len = bytes.Length;

                // // Console.WriteLine("I got here");

                // var write = instance.Functions.Where(f => f.Name == "store_data").First();

                // //var resptr = ((dynamic)instance).store_data(ptr, len);

                // var resptr = write.Invoke(ptr, len);

                // ((dynamic)instance).new_dealloc(ptr, len);

                // var resultLen = ((dynamic)instance).get_response_len();

                // var result = new byte[resultLen];

                // var memResult = memory.Span[resptr..(resptr + resultLen)];

                // memResult.CopyTo(result);

                // WriteResponse resMessage;

                // resMessage = WriteResponse.Parser.ParseFrom(result);

                // Console.WriteLine($"This is the status of message: {resMessage.Ok}");

                // ((dynamic)instance).new_dealloc(resptr, resultLen);



                // read test 

                var message = new ReadRequest();
                message.FileName = "Sparta";



                var bytes = message.ToByteArray();

                var ptr = ((dynamic)instance).new_alloc(bytes.Length);

                Console.WriteLine($"I got here and this is the ptr: {ptr}");

                bytes.CopyTo(memory.Span[ptr..]);

                var len = bytes.Length;

                // Console.WriteLine("I got here");

                var read = instance.Functions.Where(f => f.Name == "read_data").First();

                //var resptr = ((dynamic)instance).store_data(ptr, len);

                var resptr = read.Invoke(ptr, len);

                ((dynamic)instance).new_dealloc(ptr, len);

                var resultLen = ((dynamic)instance).get_response_len();

                var result = new byte[resultLen];

                var memResult = memory.Span[resptr..(resptr + resultLen)];

                memResult.CopyTo(result);

                ReadResponse resMessage;

                resMessage = ReadResponse.Parser.ParseFrom(result);

                Console.WriteLine($"This is the value of message: {resMessage.Value}");
                Console.WriteLine($"This is the time of message: {resMessage.Timestamp}");

                ((dynamic)instance).new_dealloc(resptr, resultLen);


            }

        }


        public override Task<ReadResponse> Read(ReadRequest request, ServerCallContext context)
        {
            //storageSingleton.wasm.store_data();
            ReadResponse output = new ReadResponse();

            if (request.FileName == "ImportantData")
            {
                output.Value = "This is incredibly top secret data!";
                output.Ok = 1;
                output.Timestamp = DateTime.UtcNow.ToTimestamp();
            }
            else
            {
                output.Value = "¨Not very important data";
                output.Ok = 1;
                output.Timestamp = DateTime.UtcNow.ToTimestamp();
            }

            return Task.FromResult(output);

        }

        public override Task<WriteResponse> Write(WriteRequest request, ServerCallContext context)
        {
            return base.Write(request, context);
        }
    }
}
