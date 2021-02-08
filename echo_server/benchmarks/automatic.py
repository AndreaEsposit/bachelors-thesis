"""
Run this script in the folder where you want all the benchmarks results (Ex:"Wasm/Wasmless" dir).
Change R_Proto depending on where you run the code
GHZ needs to be changed depending on where ghz is (if it is installed with brew or go buil)

Program can be runned also like this: python3 autokatic.py name_result_file(ex: 1-Client) 200(number of messages) number_of_clients
@author: Andrea Esposito
"""

import json
import sys
import os
import subprocess
from time import sleep


THIS_FOLDER = os.path.dirname(os.path.abspath(__file__))
# RELATIVE_PROTOFILE_POSITION
R_PROTO = "../../proto/echo.proto"
GHZ = "../../../../ghz/cmd/ghz/ghz"



def run(cmd):
    p = subprocess.Popen(cmd)
    p.wait()


def runAllBenchMarks(number_of_benchmarks: int, clients: int, number_of_messages: int, benchmarks_name: str, port: str):
    for i in range(number_of_benchmarks):
        run([GHZ, "--insecure", "--proto", R_PROTO, "--call",
             "proto.Echo.Send", "-c", str(
                 clients), "-n", str(number_of_messages),
             "-d", "{\"content\":\"Random string\"}", "-o",
             "(" + str(i+1) + ")" + benchmarks_name + str(number_of_messages) + ".json", "-O", "pretty", port])
        sleep(0.5)  # Sleep half 1s


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
    ## Default Values
    numberOfClients = 1
    numberOfMessages = 2
    resultName = "AvgResult"

    args = len(sys.argv)
    if args > 1:
        if args == 3:
            resultName = sys.argv[1]
            numberOfMessages = sys.argv[2]
            numberOfClients = sys.argv[3]
        if args >=2:
            resultName = sys.argv[1]
            numberOfMessages = sys.argv[2]
        if args == 3:
            resultName = sys.argv[1]

    
    runAllBenchMarks(20, numberOfClients, numberOfMessages, os.path.basename(
        THIS_FOLDER), "152.94.1.102:50051")
    getUsefulData(resultName )
    cleanUp(os.path.basename(
        THIS_FOLDER))


if __name__ == "__main__":
    main()
