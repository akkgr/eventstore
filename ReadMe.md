aws dynamodb create-table \
    --table-name Events \
    --attribute-definitions AttributeName=aggregateId,AttributeType=S AttributeName=eventNumber,AttributeType=N \
    --key-schema AttributeName=aggregateId,KeyType=HASH AttributeName=eventNumber,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

aws dynamodb create-table \
    --table-name Aggregates \
    --attribute-definitions AttributeName=aggregateId,AttributeType=S \
    --key-schema AttributeName=aggregateId,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5