package dynamodbstore

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func newTable(dbc *DynamoDBClient, tablename string, tableDesc *dynamodb.CreateTableInput) (*types.TableDescription, error) {
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

func createEventsTable(dbc *DynamoDBClient) (*types.TableDescription, error) {
	ti := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("Id"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("Version"),
			AttributeType: types.ScalarAttributeTypeN,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("Id"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("Version"),
			KeyType:       types.KeyTypeRange,
		}},
		TableName: aws.String(dbc.eventsTable),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	return newTable(dbc, dbc.eventsTable, ti)
}

func createAggregatesTable(dbc *DynamoDBClient) (*types.TableDescription, error) {
	ti := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("Id"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("Id"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName: aws.String(dbc.aggregateTable),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	return newTable(dbc, dbc.aggregateTable, ti)
}

func CreateTables(dbc *DynamoDBClient) error {
	_, err := createEventsTable(dbc)
	if err != nil {
		return err
	}
	_, err = createAggregatesTable(dbc)
	if err != nil {
		return err
	}
	return nil
}
