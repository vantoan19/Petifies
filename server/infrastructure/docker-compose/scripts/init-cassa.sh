#!/usr/bin/env bash

until printf "" 2>>/dev/null >>/dev/tcp/newfeed-cassandra/9042; do 
    sleep 5;
    echo "Waiting for cassandra to come up...";
done

echo "Creating keyspace and table..."
cqlsh newfeed-cassandra -u cassandra -p cassandra -e "CREATE KEYSPACE IF NOT EXISTS newfeed WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};"