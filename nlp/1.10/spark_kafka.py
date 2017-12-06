################################################################################################
# http://spark.apache.org/docs/2.0.2/ml-guide.html
# http://spark.apache.org/docs/2.0.2/sql-programming-guide.html
# http://spark.apache.org/docs/2.0.2/ml-features.html#tf-idf
# http://spark.apache.org/docs/2.0.2/api/python/pyspark.ml.html#pyspark.ml.feature.IDF
################################################################################################


from pyspark.ml.feature import HashingTF, IDF, Tokenizer, CountVectorizer
from pyspark.sql.types import *
from pyspark.sql.functions import *
from pyspark.ml.linalg import Vectors, SparseVector
from pyspark.ml.clustering import LDA, LDAModel, BisectingKMeans
from pyspark.sql.functions import monotonically_increasing_id
from pyspark.sql import SparkSession

import sys
import os
if os.path.exists('libs.zip'):
    sys.path.insert(0, 'libs.zip')

from kafka import KafkaProducer

producer = KafkaProducer(bootstrap_servers='broker.kafka.l4lb.thisdcos.directory:9092')

spark = SparkSession \
    .builder \
    .appName("Python Spark SQL basic example") \
    .getOrCreate()

################################################################################################

rawdata = spark.read.format("kafka").option("kafka.bootstrap.servers", "broker.kafka.l4lb.thisdcos.directory:9092") \
        .option("subscribe", "pre-processed") \
        .load()

rawdata = rawdata.selectExpr("CAST(value as STRING)")
rawdata = rawdata.withColumn("uid", monotonically_increasing_id())

def split_text(record):
    text  = record[0]
    uid   = record[1]
    words = text.split()
    return words

udf_splittext = udf(split_text , ArrayType(StringType()))
splittext = rawdata.withColumn("words", udf_splittext(struct([rawdata[x] for x in rawdata.columns])))

# Term Frequency Vectorization 
cv = CountVectorizer(inputCol="words", outputCol="rawFeatures", vocabSize = 1000)
cvmodel = cv.fit(splittext)
featurizedData = cvmodel.transform(splittext)

vocab = cvmodel.vocabulary
vocab_broadcast = spark.sparkContext.broadcast(vocab)

idf = IDF(inputCol="rawFeatures", outputCol="features")
idfModel = idf.fit(featurizedData)
rescaledData = idfModel.transform(featurizedData)

lda = LDA(maxIter=100, k=20, seed=123, optimizer="em", featuresCol="features")

ldamodel = lda.fit(rescaledData)
ldamodel.isDistributed()

ldatopics = ldamodel.describeTopics()

def map_termID_to_Word(termIndices):
    words = []
    for termID in termIndices:
        words.append(vocab_broadcast.value[termID])

    return words

udf_map_termID_to_Word = udf(map_termID_to_Word , ArrayType(StringType()))
ldatopics_mapped = ldatopics.withColumn("topic_desc", udf_map_termID_to_Word(ldatopics.termIndices))
ldaptopics_mapped_results = ldatopics_mapped.select(ldatopics_mapped.topic, ldatopics_mapped.topic_desc)

for x in ldaptopics_mapped_results.collect():
    producer.send('spark-output', " ".join(x[1]).encode('utf-8'))

ldaResults = ldamodel.transform(rescaledData)
