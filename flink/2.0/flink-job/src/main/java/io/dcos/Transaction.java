package io.dcos;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

/**
 * Created by joergschad on 22.07.17.
 */
public class Transaction {
    private long timestamp;
    private int origin;
    private int target;
    private int amount;

    public int getOrigin() {
        return origin;
    }

    public void setOrigin(int origin) {
        this.origin = origin;
    }

    public int getTarget() {
        return target;
    }

    public void setTarget(int target) {
        this.target = target;
    }

    public int getAmount() {
        return amount;
    }

    public void setAmount(int amount) {
        this.amount = amount;
    }

    public Transaction() { }

    public Transaction(long timestamp, int origin, int target, int amount) {
        this.timestamp = timestamp;
        this.origin = origin;
        this.target = target;
        this.amount = amount;
    }


    public Transaction(String transaction) {


        String[] parts = transaction.split(";");

        System.out.println ("Generating: "+ transaction );


        try {
            this.timestamp = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'").parse(parts[0]).getTime();
            this.origin = Integer.parseInt(parts[1]);
            this.target = Integer.parseInt(parts[2]);
            this.amount = Integer.parseInt(parts[3]);

        } catch (Exception e) {
            e.printStackTrace();
        }
        ;
        System.out.println ("Generated Transaction: "+ this.toString () );
    }

    public long getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }



    @Override
    public String toString() {
        return "Transaction{" +
                "timestamp=" + timestamp +
                ", origin=" + origin +
                ", target='" + target + '\'' +
                ", amount=" + amount +
                '}';
    }
}