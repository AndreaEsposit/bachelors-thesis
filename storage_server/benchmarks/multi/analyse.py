"""
analyse.py is used to analyse the cvs files that you get from the custom benchmark

Program can be runned like this: python3 analyse.py
@author: Andrea Esposito
"""
import json
import sys
import os
import csv


THIS_FOLDER = os.path.dirname(os.path.abspath(__file__))


def getUsefulData(fileName: str):
    new_data = {}

    new_data["Avg"] = []

    nFiles = 0

    # latency in microseconds
    sumLatency = 0
    nRequests = 0

    # time(UnixFormat) : times
    times = {}

    obj = os.scandir(THIS_FOLDER)
    for entry in obj:
        if entry.is_file():
            if entry.name.endswith(".csv"):
                real_path = os.path.join(THIS_FOLDER, entry.name)
                # Opening csv file
                with open(real_path, "r") as file:
                    nFiles += 1
                    reader = csv.reader(file)
                    # Skip first row in the file
                    firstline = True
                    for row in reader:
                        if firstline:
                            firstline = False
                            continue
                        nRequests += 1
                        sumLatency += int(row[0])/1000
                        if row[1] in times.keys():
                            times[row[1]] += 1
                        else:
                            times[row[1]] = 0

    sumRps = 0  # sum of Request per second
    for _, val in times.items():
        sumRps += val

    avgRps = (sumRps / len(times))  # Gest avg in ms
    avgLatency = sumLatency / nRequests

    # print(times)
    # print("This is sumRps " + str(sumRps))
    # print("This is avgRps " + str(avgRps))
    # print("This is nRequests " + str(nRequests))
    # print("This is avgLatency " + str(avgLatency))

    new_data["Avg"].append({
        "latency(ms)": avgLatency,
        "requests/sec(throughput)": avgRps,
    })

    newfile = os.path.join(THIS_FOLDER, fileName + ".json")
    with open(newfile, "w") as outfile:
        json.dump(new_data, outfile, indent=3, sort_keys=True)


def main():
    getUsefulData(fileName="result")


if __name__ == "__main__":
    main()
