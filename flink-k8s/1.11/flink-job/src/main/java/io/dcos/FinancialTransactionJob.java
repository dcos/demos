package io.dcos;

import org.apache.flink.api.common.functions.*;
import org.apache.flink.streaming.api.TimeCharacteristic;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.datastream.SingleOutputStreamOperator;
import org.apache.flink.streaming.api.functions.timestamps.AscendingTimestampExtractor;
import org.apache.flink.streaming.api.windowing.assigners.SlidingEventTimeWindows;
import org.apache.flink.streaming.api.windowing.time.Time;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.connectors.kafka.FlinkKafkaConsumer09;
import org.apache.flink.streaming.connectors.kafka.FlinkKafkaProducer09;
import org.apache.flink.streaming.util.serialization.SimpleStringSchema;
import org.apache.flink.api.java.utils.ParameterTool;
import java.util.Properties;

import static java.lang.Math.max;
import static java.lang.Math.min;


public class FinancialTransactionJob {

    public static void main(String[] args) throws Exception {

        // Set parameters
        ParameterTool parameters = ParameterTool.fromArgs (args);
        String inputTopic = parameters.get ("inputTopic", "transactions");
        String outputTopic = parameters.get ("outputTopic", "fraud");
        String kafka_host = parameters.get ("kafka_host", "broker.kafka.l4lb.thisdcos.directory:9092");

        // create execution environment
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment ();
        env.setStreamTimeCharacteristic (TimeCharacteristic.EventTime);


        Properties properties = new Properties ();
        properties.setProperty ("bootstrap.servers", kafka_host);
        properties.setProperty ("group.id", "flink_consumer");


        DataStream<Transaction> transactionsStream = env
                .addSource (new FlinkKafkaConsumer09<> (inputTopic, new SimpleStringSchema (), properties))
                // Map from String from Kafka Stream to Transaction.
                .map (new MapFunction<String, Transaction> () {
                    @Override
                    public Transaction map(String value) throws Exception {
                        return new Transaction (value);
                    }
                });

        transactionsStream.print ();

        // Extract timestamp information to support 'event-time' processing
        SingleOutputStreamOperator<Transaction> timestampedStream = transactionsStream.assignTimestampsAndWatermarks (
                new AscendingTimestampExtractor<Transaction> () {
                    @Override
                    public long extractAscendingTimestamp(Transaction element) {
                        return element.getTimestamp ();
                    }
                });
        timestampedStream.print ();

        DataStream<TransactionAggregate> rate_count = timestampedStream
                .keyBy ("origin", "target")
                // Sum over ten minute
                .window (SlidingEventTimeWindows.of (Time.minutes (10), Time.minutes (2)  ))
                // Fold into Aggregate.
                .fold (new TransactionAggregate (), new FoldFunction<Transaction,TransactionAggregate> () {
                    @Override
                     public TransactionAggregate fold( TransactionAggregate transactionAggregate, Transaction transaction) {
                         transactionAggregate.transactionVector.add (transaction);
                         transactionAggregate.amount += transaction.getAmount ();
                         transactionAggregate.startTimestamp = min(transactionAggregate.startTimestamp, transaction.getTimestamp ());
                         transactionAggregate.endTimestamp = max (transactionAggregate.endTimestamp, transaction.getTimestamp ());

                         return transactionAggregate;
            }
        });

        DataStream<String> kafka_output = rate_count
                .filter (new FilterFunction<TransactionAggregate> () {
                             @Override
                             public boolean filter(TransactionAggregate transaction) throws Exception {
                                 // Output if summed amount greater than 10000
                                 return (transaction.amount > 10000) ;
                }})
                .map (new MapFunction<TransactionAggregate, String> () {
                    @Override
                    public String map(TransactionAggregate transactionAggregate) throws Exception {
                            return transactionAggregate.toString ();
                    }
                });

        rate_count.print ();


        kafka_output.addSink (new FlinkKafkaProducer09<> (outputTopic, new SimpleStringSchema (), properties));

        env.execute ();
    }
}
