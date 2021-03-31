using Grpc.Core;
using GrpcServer;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;

namespace DotNetServer.Services
{
    public class StorageService : Storage.StorageBase
    {
        private readonly WasmSingleton wasm;
        //private readonly ILogger<StorageService> _logger;

        public StorageService(WasmSingleton staticWasm)
        { //ILogger<StorageService> logger
            this.wasm = staticWasm;
            //_logger = logger;
        }


        public override Task<ReadResponse> Read(ReadRequest request, ServerCallContext context)
        {

            var result = wasm.callWasm<ReadResponse>("read_data", request);
            // Console.WriteLine($"This is the value of message: {resMessage.Value}");
            // Console.WriteLine($"This is the time of message: {resMessage.Timestamp}");
            return Task.FromResult(result);
        }

        public override Task<WriteResponse> Write(WriteRequest request, ServerCallContext context)
        {
            var result = wasm.callWasm<WriteResponse>("store_data", request);
            //Console.WriteLine($"This is the status of message: {resMessage.Ok}");
            return Task.FromResult(result);
        }
    }
}
