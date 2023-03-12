#!/bin/bash
echo "sleeping for 10 seconds"
sleep 10

echo rs-init.sh time now: `date +"%T" `
mongosh --host post-mongo-1:27017 <<EOF
  var cfg = {
    "_id": "my-mongo-set",
    "version": 1,
    "members": [
      {
        "_id": 0,
        "host": "post-mongo-1:27017",
        "priority": 2
      },
      {
        "_id": 1,
        "host": "post-mongo-2:27017",
        "priority": 1
      }, 
      {
        "_id": 2,
        "host": "post-mongo-3:27017",
        "priority": 1
      }
    ]
  };
  rs.initiate(cfg);
EOF