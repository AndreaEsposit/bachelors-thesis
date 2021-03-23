using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using GrpcServer;
using System.Text.Json;
using System.IO;
using System.Threading;
using System.Threading.Tasks;
using System;
using JsonSerializer = Utf8Json.JsonSerializer;

namespace DotNetServer.Services
{
    // I use this to store the lock (not ideal to have 1 lock for all file writes, but it is ok for benchmarking, since we use only one file)

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

        public readonly Mutex Mu = new Mutex();

        public JsonSerializerOptions Options = new JsonSerializerOptions
        {
            WriteIndented = true,
        };


    }
    public class Content
    {
        public int nseconds { get; set; }
        public long seconds { get; set; }
        public string value { get; set; }
    }

    public class StorageService : Storage.StorageBase
    {
        //private readonly ILogger<StorageService> _logger;

        private WasmSingleton wasmSingleton = WasmSingleton.Instance;


        public StorageService() //ILogger<StorageService> logger
        {
            // Set up the logger
            // _logger = logger;

            if (!wasmSingleton.instanceReady)
            {
                wasmSingleton.instanceReady = true;

                Console.WriteLine("I am ready for the benchmark!");
            }
        }

        public override Task<ReadResponse> Read(ReadRequest request, ServerCallContext context)
        {
            ReadResponse result = new ReadResponse();


            if (File.Exists($"./data/{request.FileName}.json"))
            {
                wasmSingleton.Mu.WaitOne(); // take lock
                Content content = JsonSerializer.Deserialize<Content>(File.ReadAllText($@"./data/{request.FileName}.json"));
                wasmSingleton.Mu.ReleaseMutex(); // release lock

                Timestamp time = new Timestamp
                {
                    Nanos = content.nseconds,
                    Seconds = content.seconds,
                };
                result.Timestamp = time;
                result.Value = content.value;
                result.Ok = 1;
                Console.WriteLine($"This is the value of message: {result.Ok}");
            }
            else
            {
                Timestamp time = new Timestamp
                {
                    Nanos = 0,
                    Seconds = 0,
                };
                result.Timestamp = time;
                result.Value = "";
                result.Ok = 0;
            }

            // Console.WriteLine($"This is the value of message: {resMessage.Value}");
            // Console.WriteLine($"This is the time of message: {resMessage.Timestamp}");
            return Task.FromResult(result);
        }

        public override Task<WriteResponse> Write(WriteRequest request, ServerCallContext context)
        {

            Content content = new()
            {
                seconds = request.Timestamp.Seconds,
                nseconds = request.Timestamp.Nanos,
                value = request.Value,
            };



            wasmSingleton.Mu.WaitOne(); // take lock
            File.WriteAllBytes($@"./data/{request.FileName}.json", JsonSerializer.Serialize(content));
            wasmSingleton.Mu.ReleaseMutex(); // release lock


            var result = new WriteResponse { Ok = 1 };
            //Console.WriteLine($"This is the status of message: {resMessage.Ok}");
            return Task.FromResult(result);
        }
    }
}
