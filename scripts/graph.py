import matplotlib.pyplot as plt
import csv

with open("/var/memory.txt", "r") as file:
    memory = []
    timestamps = []
    records = csv.reader(file)
    for record in records:
        timestamps.append(int(record[0]))
        memory.append(float(record[1]) / 1000.0)
    plt.plot(timestamps, memory)
    plt.savefig("graph.PNG")