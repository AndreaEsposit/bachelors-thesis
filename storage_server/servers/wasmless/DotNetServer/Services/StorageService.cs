using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using GrpcServer;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using System.IO;
using System.Threading.Tasks;


namespace DotNetServer.Services
{

    public class Content
    {
        public int nseconds { get; set; }
        public long seconds { get; set; }
        public string value { get; set; }
    }
    public class StorageService : Storage.StorageBase
    {
        private readonly ILogger<StorageService> _logger;

        public StorageService(ILogger<StorageService> logger)
        {
            _logger = logger;
        }


        public override Task<ReadResponse> Read(ReadRequest request, ServerCallContext context)
        {
            ReadResponse result = new ReadResponse();
            if (File.Exists($"./data/{request.FileName}.json"))
            {
                Content content = JsonConvert.DeserializeObject<Content>(File.ReadAllText($@"./data/{request.FileName}.json"));

                Timestamp time = new Timestamp
                {
                    Nanos = content.nseconds,
                    Seconds = content.seconds,
                };
                result.Timestamp = time;
                result.Value = content.value;
                result.Ok = 1;
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
            File.WriteAllText($@"./data/{request.FileName}.json", JsonConvert.SerializeObject(content));

            var result = new WriteResponse { Ok = 1 };
            //Console.WriteLine($"This is the status of message: {resMessage.Ok}");
            return Task.FromResult(result);
        }
    }
}
