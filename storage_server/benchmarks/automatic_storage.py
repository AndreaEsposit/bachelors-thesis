"""
Run this script in the folder where you want all the benchmarks results (Ex:"Wasm/Wasmless" dir).
Change R_Proto depending on where you run the code
GHZ needs to be changed depending on where ghz is (if it is installed with brew or go buil)

Program can be runned like this: python3 autokatic.py name_result_file(ex: 1-Client) 200(number of messages) number_of_clients mode(read/r or write/w)
@author: Andrea Esposito
"""

import json
import sys
import os
import subprocess
from time import sleep


THIS_FOLDER = os.path.dirname(os.path.abspath(__file__))
# RELATIVE_PROTOFILE_POSITION
R_PROTO = "../../proto/storage.proto"
GHZ = "../../../../ghz/cmd/ghz/ghz"

# Write Benchmark message:
read_file_name = "test"
read = "{\"FileName\":\"" + read_file_name + "\"}"

# Read Benchmark message:
# 10 bytes
write_conetent = "Testing!!!"

# 1Kb
_1kb = "Bruce Wayne was born to wealthy physician Dr. Thomas Wayne and his wife Martha, who were themselves members of the prestigious Wayne and Kane families of Gotham City, respectively. When he was three, Bruce's mother Martha was expecting a second child to be named Thomas Wayne, Jr. However, because of her intent to found a school for the underprivileged in Gotham, she was targeted by the manipulative Court of Owls, who arranged for her to have a car accident. She and Bruce survived, but the accident forced Martha into premature labor, and the baby was lost. While on vacation to forget about these events, the Wayne Family butler, Jarvis Pennyworth was killed by one of the Court of Owls' Talons. A letter he'd written to his son Alfred, warning him away from the beleaguered Wayne family, was never delivered. As such, Alfred - who had been an actor at the Globe Theatre at the time and a military medic before that, traveled to Gotham City to take up his father's place, serving the Waynes....."
_10kb = _1kb*10

write = "{\"FileName\":\"test\", \"Value\":\"" + write_conetent + \
    "\", \"Timestamp\":\"" + "2006-01-02T15:04:05.999999999Z" + "\"}"


def run(cmd):
    p = subprocess.Popen(cmd)
    p.wait()


def runAllBenchMarks(number_of_benchmarks: int, clients: int, number_of_messages: int, data: str, benchmarks_name: str, address: str, port: int, mode: str):
    for i in range(number_of_benchmarks):
        run([GHZ, "--insecure", "--proto", R_PROTO, "--call",
             "proto.Storage." + mode, "-c", str(
                 clients), "-n", str(number_of_messages),
             "-d", data, "-o",
             "(" + str(i + 1) + ")" + benchmarks_name + str(number_of_messages) + ".json", "-O", "pretty", address + ":" + str(port + i)])


def getUsefulData(result: str):
    new_data = {}

    new_data["benchmarks"] = []
    new_data["final_Avg."] = []

    sumAllLatency = 0
    sumAllThroughput = 0
    sumResponseTime = 0

    obj = os.scandir(THIS_FOLDER)
    for entry in obj:
        if entry.is_file():
            if entry.name.endswith(".json"):
                myfile = os.path.join(THIS_FOLDER, entry.name)
                # Opening JSON file
                with open(myfile, "r") as f:
                    # Return JSON object as a directory
                    data = json.load(f)

                    # Iterating though the json list
                    details = data["details"]

                    # Latency in ns
                    sum_latency = 0
                    for listing in details:
                        sum_latency += listing["latency"]
                    avg = sum_latency / len(details)

                    new_data["benchmarks"].append({
                        "file_name": entry.name.replace(".json", ""),
                        "avg_latency(ms)": float(avg * (10 ** -6)),
                        "avg_respons_time(ms)": data["average"] * (10 ** -6),
                        "requests/sec(throughput)": data["rps"],
                    })
                    sumAllLatency += float(avg * (10 ** -6))
                    sumResponseTime += data["average"] * (10 ** -6)
                    sumAllThroughput += data["rps"]

    number_bs = len(new_data["benchmarks"])
    new_data["final_Avg."].append({
        "latency(ms)": sumAllLatency / number_bs,
        "respons_time(ms)": sumResponseTime / number_bs,
        "requests/sec(throughput)": sumAllThroughput / number_bs,

    })

    newfile = os.path.join(THIS_FOLDER, result + ".json")
    with open(newfile, "w") as outfile:
        json.dump(new_data, outfile, indent=3, sort_keys=True)


def cleanUp(name: str):
    obj = os.scandir(THIS_FOLDER)
    for entry in obj:
        if name in entry.name:
            os.remove(entry.name)


def main():

    resultName = sys.argv[1]
    numberOfMessages = sys.argv[2]
    numberOfClients = sys.argv[3]
    testMode = sys.argv[4]

    if testMode == "r" or testMode == "read":
        data = read
        mode = "Read"
    elif testMode == "w" or testMode == "write":
        data = write
        mode = "Write"
    else:
        print("Write a valid mode! (read/r or write/w")
        return

    runAllBenchMarks(number_of_benchmarks=10, clients=numberOfClients, number_of_messages=numberOfMessages, data=data, benchmarks_name=os.path.basename(
        THIS_FOLDER), address="localhost", port=50051, mode=mode)
    getUsefulData(resultName)
    cleanUp(os.path.basename(THIS_FOLDER))


if __name__ == "__main__":
    main()
