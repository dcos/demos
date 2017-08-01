package io.dcos;

import org.apache.flink.api.common.functions.FlatMapFunction;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;
import org.apache.flink.streaming.api.windowing.time.Time;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.connectors.kafka.FlinkKafkaConsumer09;
import org.apache.flink.streaming.connectors.kafka.FlinkKafkaProducer09;
import org.apache.flink.streaming.util.serialization.SimpleStringSchema;
import org.apache.flink.api.java.tuple.Tuple2;
import org.apache.flink.api.java.utils.ParameterTool;
import org.apache.flink.util.Collector;
import java.util.Properties;
import java.util.ArrayList;
import java.util.Arrays;
import static java.lang.Character.toLowerCase;

public class FinancialTransactionJob {

    public static void main(String[] args) throws Exception {

        // Set parameters
        ParameterTool parameters = ParameterTool.fromArgs(args);
        String job = parameters.getRequired("job");
        int period = Integer.parseInt(parameters.getRequired("period"));
        String kafka_host = parameters.getRequired("kafka_host");

        // create execution environment
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        Properties properties = new Properties();
        properties.setProperty("bootstrap.servers", kafka_host);
        properties.setProperty("group.id", "flink_consumer");

        // Create input streams
        DataStream<String> nyc = env
                .addSource(new FlinkKafkaConsumer09<>("NYC", new SimpleStringSchema(), properties))
                .map(new MapFunction<String, String>() {

                    @Override
                    public String map(String value) throws Exception {
                        return "NYC " + value;
                    }
                });
        DataStream<String> moscow = env
                .addSource(new FlinkKafkaConsumer09<>("Moscow", new SimpleStringSchema(), properties))
                .map(new MapFunction<String, String>() {

                    @Override
                    public String map(String value) throws Exception {
                        return "Moscow " + value;
                    }
                });
        DataStream<String> tokyo = env
                .addSource(new FlinkKafkaConsumer09<>("Tokyo", new SimpleStringSchema(), properties))
                .map(new MapFunction<String, String>() {

                    @Override
                    public String map(String value) throws Exception {
                        return "Tokyo " + value;
                    }
                });
        DataStream<String> sf = env
                .addSource(new FlinkKafkaConsumer09<>("SF", new SimpleStringSchema(), properties))
                .map(new MapFunction<String, String>() {

                    @Override
                    public String map(String value) throws Exception {
                        return "SF " + value;
                    }
                });
        DataStream<String> merged = env
                .addSource(new FlinkKafkaConsumer09<>("London", new SimpleStringSchema(), properties))
                .map(new MapFunction<String, String>() {

                    @Override
                    public String map(String value) throws Exception {
                        return "London " + value;
                    }
                })
                .union(nyc, moscow, tokyo, sf);

        // Split string into Tuple2 of Location, Amount
        DataStream<Tuple2<String, Integer>> split_stream = merged
                .flatMap(new Splitter());

        DataStream<String> kafka_output;

        switch (job) {
            case "dollar_rate":
                // Rate in dollars over period
                DataStream<Tuple2<String, Integer>> dollar_rate = split_stream
                        .keyBy(0)
                        .timeWindow(Time.seconds(period))
                        .sum(1);
                // Convert tuples to strings
                kafka_output = dollar_rate
                        .flatMap(new Stringify());
                kafka_output.addSink(new FlinkKafkaProducer09<>("dollarrate", new SimpleStringSchema(), properties));
                dollar_rate.print();
            case "max_dollars":
                // Max transfer by location during period
                DataStream<Tuple2<String, Integer>> max_dollars = split_stream
                        .keyBy(0)
                        .timeWindow(Time.seconds(period))
                        .maxBy(1);
                kafka_output = max_dollars
                        .flatMap(new Stringify());
                kafka_output.addSink(new FlinkKafkaProducer09<>("maxdollars", new SimpleStringSchema(), properties));
                max_dollars.print();
            case "rate_count":
                // Number of transfers in period
                DataStream<Tuple2<String, Integer>> rate_count = split_stream
                        .flatMap(new Counter())
                        .keyBy(0)
                        .timeWindow(Time.seconds(period))
                        .sum(1);
                kafka_output = rate_count
                        .flatMap(new Stringify());
                kafka_output.addSink(new FlinkKafkaProducer09<>("ratecount", new SimpleStringSchema(), properties));
                rate_count.print();
                break;
            default:
                break;
        }

        env.execute();
    }

    // Data Types

    public static class Splitter implements FlatMapFunction<String, Tuple2<String, Integer>> {
        @Override
        public void flatMap(String value, Collector<Tuple2<String, Integer>> out) throws Exception {
            String[] tokens = value.split(" ");
            out.collect(new Tuple2<String, Integer>(tokens[0], Integer.parseInt(tokens[3])));
        }
    }

    public static class Counter implements FlatMapFunction<Tuple2<String, Integer>, Tuple2<String, Integer>> {
        @Override
        public void flatMap(Tuple2<String, Integer> input, Collector<Tuple2<String, Integer>> out) throws Exception {
            out.collect(new Tuple2<String, Integer>(input.f0, 1));
        }
    }

    public static class Stringify implements FlatMapFunction<Tuple2<String, Integer>, String> {
        @Override
        public void flatMap(Tuple2<String, Integer> input, Collector<String> out) throws Exception {
            out.collect(new String(input.f0 + " " + Integer.toString(input.f1)));
        }
    }

    private static final ArrayList<String> LOCATIONS = new ArrayList<String>(
            Arrays.asList("SF", "NYC", "London", "Moscow", "Tokyo"));

}