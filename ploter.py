"""Simple script for plotting the contents of results.csv
"""

import matplotlib.pyplot as plt
import pandas as pd
import numpy as np

plt.figure(dpi=500)

df = pd.read_csv("results.csv")

x = df["total_csvs"]
y = df["memory_consumption"] / 1e6

plt.plot(x, y, "yo")

plt.xlabel("# Total CSVs uploaded")
plt.ylabel("Memory usage (GB)")
plt.ylim(ymin=0)

plt.legend()
plt.grid(axis="x", color="0.95")

plt.savefig("results.jpg")
