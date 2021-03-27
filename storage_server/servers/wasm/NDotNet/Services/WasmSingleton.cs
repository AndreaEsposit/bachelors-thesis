using Google.Protobuf;
using System;
using System.Collections.Generic;
using System.Linq;
using Wasmtime;

namespace DotNetServer.Services
{
    public class WasmSingleton
    {
        public WasmSingleton(string[] services, string wasmLocation)
        {
            using var engine = new Engine();
            using var store = new Store(engine);

            // pass access to data directory to this Wasm module
            WasiConfiguration wasiConfiguration = new WasiConfiguration();
            wasiConfiguration.WithPreopenedDirectory("./data", ".");

            // Create the WebAssembly-module
            using var module = Module.FromFile(engine, wasmLocation);

            // get instance of the wasm module ans specify the version of wasi
            using var host = new Host(store);
            host.DefineWasi("wasi_snapshot_preview1", wasiConfiguration);
            instance = host.Instantiate(module);

            // execure the _initialize function to given wasm access to the data folder
            ((dynamic)instance)._initialize();

            // export the only memory that we are using with this module 
            memory = instance.Memories.Where(m => m.Name == "memory").First();

            // export application functions from the WebAssembly module
            funcs = new Dictionary<String, Wasmtime.Externs.ExternFunction>();

            foreach (string func in services)
            {
                funcs[func] = instance.Functions.Where(f => f.Name == func).First();
            }
            Console.WriteLine("Wasm instance ready!");
        }

        // mutex lock to control wasm access
        private readonly object wasmLock = new object();

        // our wasm instance that gets made with lazy initialization 
        private Instance instance;
        private Wasmtime.Externs.ExternMemory memory;

        // keeps track of the functions used by the applicaion (no alloc, dealloc etc..)
        private Dictionary<String, Wasmtime.Externs.ExternFunction> funcs;

        public T callWasm<T>(String fn, IMessage message) where T : IMessage<T>, new()
        {
            // --- Copy the buffer to the module's memory 
            var bytes = message.ToByteArray();

            byte[] result;

            lock (wasmLock)
            {
                var ptr = ((dynamic)instance).new_alloc(bytes.Length);
                var len = bytes.Length;

                bytes.CopyTo(memory.Span[ptr..]);

                // you can technically define the function each and every time
                //var func = wasm.Functions.Where(f => f.Name == fn).First();

                // call the function you are intrested in
                var resPtr = funcs[fn].Invoke(ptr, len);

                // deallocate request protobuf message
                ((dynamic)instance).new_dealloc(ptr, len);

                // get Wasm response lenght
                var resultLen = ((dynamic)instance).get_response_len();

                // copy byte array fom WebAssembly to a result byte array
                var resultMemory = memory.Span[resPtr..(resPtr + resultLen)];
                result = new byte[resultLen];
                resultMemory.CopyTo(result);

                // deallocate repsonse protobud message form Wasmtime's memory
                ((dynamic)instance).new_dealloc(resPtr, resultLen);
            }
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

}
