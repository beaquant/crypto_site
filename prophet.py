#!/usr/bin/python3

import sys
import datetime
import pandas as pd
from fbprophet import Prophet
import matplotlib.pyplot as plt

df = pd.DataFrame({"ds" : [], "y": []})
ds_arr = []
y_arr = []
for line in sys.stdin:
    parts = line.strip().split(",")
    ds_arr.append(datetime.datetime.fromtimestamp(float(parts[0])))
    y_arr.append(float(parts[1]))
    # print(parts[1] + " " + parts[2])
    # print(datetime.datetime.fromtimestamp(float(parts[1])).strftime('%Y-%m-%dT%H:%M:%SZ') + " " + parts[2])

df = pd.DataFrame({"ds" : ds_arr, "y": y_arr})
m = Prophet()
m.fit(df)

future = m.make_future_dataframe(periods=0)
forecast = m.predict(future)
# print(forecast[['ds', 'yhat']])
plt.plot(forecast['yhat'], color="#aa0000")
plt.plot(y_arr, color="#00aa00")
plt.show()
forecast.to_csv("data/ethereum_predictions.csv", columns = ["ds", "yhat"])
print(forecast["yhat"])

# print(df)
