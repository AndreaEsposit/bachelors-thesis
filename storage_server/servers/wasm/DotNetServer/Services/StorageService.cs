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

                using var module = Module.FromFile(engine, "write.wasm");

                using var host = new Host(store);
                host.DefineWasi("wasi_snapshot_preview1", wasiConfiguration);
                var instance = host.Instantiate(module);
                ((dynamic)instance)._initialize();


                storageSingleton.wasm = instance;
                storageSingleton.fatto = true;


                Console.WriteLine("Instance ready");

            }


        }




        public override Task<ReadResponse> Read(ReadRequest request, ServerCallContext context)
        {
            storageSingleton.wasm.store_data();
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
