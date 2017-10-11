csv = read.csv(file="data/ethereum_price.csv", sep=",")
csv_data = csv[,2]
csv = read.csv(file="data/price_predictions.csv", sep=",")
pred_data=csv[,1]
x1  <- seq(1, length(csv_data), 1)
x2  <- seq(length(csv_data) + 1, length(csv_data) + length(pred_data), 1)
plot(x1,csv_data,type="l",col="red")
lines(x2,pred_data,type="l",col="green")
# plot(x2,pred_data,type="l",col="green")


csv = read.csv(file="data/ethereum_stats.csv", sep=",")
csv_data = csv[,2]
csv = read.csv(file="data/stats_predictions.csv", sep=",")
pred_data=csv[,1]
x1  <- seq(1, length(csv_data), 1)
x2  <- seq(length(csv_data) + 1, length(csv_data) + length(pred_data), 1)
plot(x1,csv_data,type="l",col="red")
lines(x2,pred_data,type="l",col="green")
# plot(x2,pred_data,type="l",col="green")
