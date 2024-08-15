package dynamodbstore

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func createTable(dbc *DynamoDBClient, tablename string, tableDesc *dynamodb.CreateTableInput) (*types.TableDescription, error) {
	table, err := dbc.store.CreateTable(context.TODO(), tableDesc)
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", tablename, err)
		return nil, err
	}
	waiter := dynamodb.NewTableExistsWaiter(dbc.store)
	err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(dbc.eventsTable)}, 5*time.Minute)
	if err != nil {
		log.Printf("Wait for table exists failed. Here's why: %v\n", err)
	}
	return table.TableDescription, nil
}

func CreateEventsTable(dbc *DynamoDBClient) (*types.TableDescription, error) {
	ti := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("version"),
			AttributeType: types.ScalarAttributeTypeN,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("version"),
			KeyType:       types.KeyTypeRange,
		}},
		TableName: aws.String(dbc.eventsTable),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	return createTable(dbc, dbc.eventsTable, ti)
}

func CreateAggregatesTable(dbc *DynamoDBClient) (*types.TableDescription, error) {
	ti := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName: aws.String(dbc.aggregateTable),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	return createTable(dbc, dbc.aggregateTable, ti)
}

func CreateSnapshotsTable(dbc *DynamoDBClient) (*types.TableDescription, error) {
	ti := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName: aws.String(dbc.snapshotsTable),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	return createTable(dbc, dbc.snapshotsTable, ti)
}
