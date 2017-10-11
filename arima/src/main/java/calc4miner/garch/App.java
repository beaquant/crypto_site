package calc4miner.garch;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.util.ArrayList;
import java.util.Map;

import timeseries.TimeSeries;
import timeseries.models.Forecast;
import timeseries.models.arima.Arima;
import timeseries.models.arima.ArimaCoefficients;
import timeseries.models.arima.ArimaForecast;
import timeseries.models.arima.ArimaOrder;
import timeseries.models.arima.ArimaSimulation;

/**
 * Hello world!
 *
 */
public class App 
{
    public static void main( String[] args )
    {
        System.out.println( "Hello World!" );

        try {
            ArrayList<Double> vals = new ArrayList<Double>();
            ArrayList<Double> timeVals = new ArrayList<Double>();
            BufferedReader reader = new BufferedReader(new FileReader("ethereum_price.csv"));
            int count = 0;
            for(String line; (line = reader.readLine()) != null; count++) {
                timeVals.add((double)count);
                double val = Double.parseDouble(line.substring(line.indexOf(",") + 1));
                vals.add(val);
            }
            reader.close();
            double[] series = new double[vals.size()];
            double[] time = new double[vals.size()];
            for (int i = 0; i < vals.size(); i++) {
                series[i] = vals.get(i);
                time[i] = timeVals.get(i);
            }

            TimeSeries timeSeries = new TimeSeries(series);
            Arima.FittingStrategy fittingStrategy = Arima.FittingStrategy.CSSML;
            ArimaCoefficients coefficients = ArimaCoefficients.newBuilder()
                                                            .setMACoeffs(-0.6760904)
                                                            .setSeasonalMACoeffs(-0.5718134)
                                                            .setDifferences(1)
                                                            .setSeasonalDifferences(1)
                                                            .setSeasonalFrequency(12)
                                                            .build();
            // Arima model = Arima.model(timeSeries, coefficients, fittingStrategy);
            Arima model = Arima.model(timeSeries, ArimaOrder.order(90, 2, 2), fittingStrategy);
            Forecast forecast = ArimaForecast.forecast(model, 364);
            double[] forecast_arr = forecast.forecast().asArray();

            System.out.println(forecast_arr.length);

            File file = new File("price_predictions.csv");
            file.createNewFile();
            BufferedWriter writer = new BufferedWriter(new FileWriter(file));
            writer.write(String.format("%.20f\n", series[series.length - 1]));
            for (int i = 0; i < forecast_arr.length; i++) 
                writer.write(String.format("%.20f\n", forecast_arr[i]));
            writer.close();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
