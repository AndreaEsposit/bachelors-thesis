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


def findFirstAndLastSec(path):
    with open(path, "r") as file:
        firstLine = True
        times = []
        for row in (csv.reader(file)):
            if firstLine:
                firstLine = False
                continue
            times.append(row[1])

        return times[0], times[-1]


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
                # Skip last second since we are likely to not end exactly at the end of the sec
                # Skip first second because of warmup
                firstSec = 0
                lastSec = 0
                # Opening csv file
                with open(real_path, "r") as file:
                    nFiles += 1

                    #print("This is first: " + str(firstSec) + " and this is last: " + str(lastSec))
                    firstLine = True
                    for row in (csv.reader(file)):
                        if firstLine:
                            firstLine = False
                            continue
                        if row[1] != firstSec and row[1] != lastSec:
                            nRequests += 1
                            sumLatency += int(row[0])/1000
                            if row[1] in times.keys():
                                times[row[1]] += 1
                            else:
                                times[row[1]] = 0

    sumRps = 0  # sum of Request per second
    for _, val in times.items():
        sumRps += val

    avgRps = (sumRps / len(times))  # Gets avg in ms
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


def cleanUp():
    obj = os.scandir(THIS_FOLDER)
    for entry in obj:
        if ".csv" in entry.name:
            os.remove(entry.name)


def main():
    getUsefulData(fileName="result")
    cleanUp()


if __name__ == "__main__":
    main()
