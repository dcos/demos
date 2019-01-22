package io.dcos;

import java.util.Vector;

/**
 * Created by joergschad on 23.07.17.
 */
public class TransactionAggregate {
    public Vector<Transaction> transactionVector;
    public long startTimestamp;
    public long endTimestamp;
    public long amount;

    public TransactionAggregate() {
        this.transactionVector = new Vector<Transaction> ();

    }

    @Override
    public String toString() {

        String result = "TransactionAggregate {" +
                "startTimestamp=" + startTimestamp +
                ", endTimestamp=" + endTimestamp +
                ", totalAmount=" + amount +
                ':';

        for (Transaction transaction : transactionVector) {
            result += System.lineSeparator()
                    + transaction.toString ();
        }

        result += "}";

        return result;
    }

}
