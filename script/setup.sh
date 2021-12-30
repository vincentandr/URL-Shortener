#!/bin/bash

echo "Starting replica set initialize"
until mongo --host mongo1 --eval "print(\"waited for connection\")"
do
    sleep 2
done
echo "Connection finished"
echo "Creating replica set"
mongo --host mongo1 <<EOF
config = {
  	"_id" : "shop-mongo-set",
  	"members" : [
  		{
  			"_id" : 0,
  			"host" : "mongo1:27017"
  		},
  		{
  			"_id" : 1,
  			"host" : "mongo2:27017"
  		},
  		{
  			"_id" : 2,
  			"host" : "mongo3:27017"
  		}
  	]
  }
  rs.initiate(config, { force: true});
  rs.reconfig(config, { force: true });
  EOF
echo "Replica set created"