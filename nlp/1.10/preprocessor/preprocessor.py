#!/usr/bin/env python
"""
Simple preprocessor
"""

import json
import re
from HTMLParser import HTMLParser
from kafka import KafkaConsumer, KafkaProducer
import argparse

class MLStripper(HTMLParser):
    """
    Parser class
    """
    def __init__(self):
        self.reset()
        self.fed = []
    def handle_data(self, d):
        self.fed.append(d)
    def get_data(self):
        return ''.join(self.fed)

def strip_tags(html):
    """
    Strip tags
    """
    tagstrip = MLStripper()
    tagstrip.feed(html)
    return tagstrip.get_data()

def remove_stopwords(record):
    """
    Remove stopwords from a string
    """
    words = record.split()
    # Core stopwords
    stopwords_core = [u'a', u'about', u'above',
                      u'after', u'again', u'against', u'all',
                      u'am', u'an', u'and', u'any', u'are',
                      u'arent', u'as', u'at', u'be', u'because',
                      u'been', u'before', u'being', u'below',
                      u'between', u'both', u'but', u'by',
                      u'can', 'cant', 'come', u'could', 'couldnt',
                      u'd', u'did', u'didn', u'do', u'does',
                      u'doesnt', u'doing', u'dont', u'down', u'during',
                      u'each', u'few', 'finally', u'for', u'from',
                      u'further', u'had', u'hadnt', u'has', u'hasnt',
                      u'have', u'havent', u'having', u'he', u'her',
                      u'here', u'hers', u'herself', u'him', u'himself',
                      u'his', u'how', u'i', u'if', u'in', u'into',
                      u'is', u'isnt', u'it', u'its', u'itself',
                      u'just', u'll', u'm', u'me', u'might', u'more',
                      u'most', u'must', u'my', u'myself', u'no',
                      u'nor', u'not', u'now', u'o', u'of', u'off',
                      u'on', u'once', u'only', u'or', u'other', u'our',
                      u'ours', u'ourselves', u'out', u'over', u'own',
                      u'r', u're', u's', 'said', u'same', u'she', u'should',
                      u'shouldnt', u'so', u'some', u'such',
                      u't', u'than', u'that', 'thats', u'the',
                      u'their', u'theirs', u'them', u'themselves',
                      u'then', u'there', u'these', u'they', u'this',
                      u'those', u'through', u'to', u'too', u'under',
                      u'until', u'up', u'very', u'was', u'wasnt',
                      u'we', u'were', u'werent', u'what', u'when',
                      u'where', u'which', u'while', u'who', u'whom',
                      u'why', u'will', u'with', u'wont', u'would',
                      u'y', u'you', u'your', u'yours', u'yourself',
                      u'yourselves']
    # Define custom words
    stopwords_custom = [u'possible', u'wondering', u'using', u'tried',
                        u'thanks', u'anyone', u'available',
                        u'cannot', u'without', u'offers', u'hoping',
                        u'provided', u'three', u'start', u'running',
                        u'looking', u'example', u'setting', u'fails',
                        u'everything', u'works', u'issue', u'stuck',
                        u'support', u'around', u'think', u'getting',
                        u'please', u'correct', u'following', u'installing',
                        u'single', u'second', u'first', u'existing',
                        u'inside', u'unable', u'creating', u'still',
                        u'error', u'following', u'given', u'appreciated'
                        u'wrong', u'however', u'every', u'always',
                        u'anything', u'showing', u'didnt', u'question',
                        u'ideas', u'confused', u'though', u'something',
                        u'nothing', u'recently', u'thank', u'recently',
                        u'problem', u'based', u'understanding', u'whats',
                        u'seems', u'someone', u'another', u'better',
                        u'solution', u'trying']
    stopwords = stopwords_core + stopwords_custom
    stopwords = [word.lower() for word in stopwords]
    # Remove stop words and words under x length
    text_out = [word.lower() for word in words if len(word) > 4 and word.lower() not in stopwords]
    return " ".join(text_out)

def main():
    """
    Preprocess from Kafka queue
    """

    parser = argparse.ArgumentParser(description="Preprocess Kafka data")
    parser.add_argument("-b" "--broker", dest="kafka_broker", help="Kafka broker address", required=True)
    parser.add_argument("-i", "--input-topic", dest="input_topic", help="Kafka input topic", default="ingest")
    parser.add_argument("-o", "--output-topic", dest="output_topic", help="Kafka output topic", default="pre-processed")


    args = parser.parse_args()

    consumer = KafkaConsumer(args.input_topic, value_deserializer=lambda m: json.loads(m.decode('ascii')),
                             auto_offset_reset='earliest',
                             bootstrap_servers=args.kafka_broker)
    producer = KafkaProducer(bootstrap_servers=args.kafka_broker)

    for message in consumer:
        question = message.value
        if 'body' in question:
            combined_text = question['body'] + question['title']
            # Remove code blocks
            cb_text = re.sub('<code>.*</code>', '', combined_text, flags=re.DOTALL)
            # Strip out any links
            links_removed = re.sub('<a href.*/a>', '', cb_text)
            # Remove newlines
            remove_newlines = re.sub(r'\n', ' ', links_removed)
            # Strip remaining tags and entities
            stripped_tags = strip_tags(remove_newlines)
            # Remove any paths
            paths_removed = re.sub(r'\w*\/\w*', '', stripped_tags)
            # Remove punctuation
            remove_punc = re.sub('[^A-Za-z0-9 ]+', '', paths_removed)
            # Remove any numeric strings ( ports etc. )
            remove_numeric = re.sub(r'\w*\d\w*', '', remove_punc)
	    # Remove stop words
            remove_stop = remove_stopwords(remove_numeric)
            producer.send(args.output_topic, remove_stop.encode('utf-8'))
        else:
            continue

if __name__ == '__main__':
    main()
