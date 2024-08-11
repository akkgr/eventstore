package eventstore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	DynamoDbClient *dynamodb.Client
	AggregateTable string
	EventsTable    string
	SnapshotsTable string
}

func NewDynamoDBClient(client *dynamodb.Client) *DynamoDBClient {
	return &DynamoDBClient{
		DynamoDbClient: client,
		AggregateTable: "Aggregates",
		EventsTable:    "Events",
		SnapshotsTable: "Snapshots",
	}
}

func (dbc *DynamoDBClient) AppendAggregate(aggregate Aggregate, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(aggregate)
	if err != nil {
		return err
	}

	_, err = dbc.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.AggregateTable), Item: item,
	})
	return err
}

func (dbc *DynamoDBClient) AppendSnapshot(snapshot Snapshot, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(snapshot)
	if err != nil {
		return err
	}

	_, err = dbc.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.SnapshotsTable), Item: item,
	})
	return err
}

func (dbc *DynamoDBClient) AppendEvent(event Event, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return err
	}

	_, err = dbc.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.EventsTable), Item: item,
	})
	return err
}

func (dbc *DynamoDBClient) GetAggregate(aggregateID string, ctx context.Context) (Aggregate, error) {
	aggregate := Aggregate{}

	id, err := attributevalue.Marshal(aggregateID)
	if err != nil {
		return aggregate, err
	}

	key := map[string]types.AttributeValue{"aggregateId": id}
	response, err := dbc.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(dbc.AggregateTable),
	})
	if err != nil {
		return aggregate, err
	}

	err = attributevalue.UnmarshalMap(response.Item, &aggregate)
	if err != nil {
		return aggregate, err
	}

	return aggregate, err
}

func (dbc *DynamoDBClient) GetEvents(aggregateID string, ctx context.Context) ([]Event, error) {
	var err error
	var response *dynamodb.QueryOutput
	var events []Event

	queryPaginator := dynamodb.NewQueryPaginator(dbc.DynamoDbClient, &dynamodb.QueryInput{
		TableName:              aws.String(dbc.EventsTable),
		KeyConditionExpression: aws.String("#pk = :pk"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "aggregateId",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: aggregateID},
		},
	})
	for queryPaginator.HasMorePages() {
		response, err = queryPaginator.NextPage(context.TODO())
		if err != nil {
			return events, err
		}
		var eventPage []Event
		err = attributevalue.UnmarshalListOfMaps(response.Items, &eventPage)
		if err != nil {
			return events, err
		}
		events = append(events, eventPage...)
	}

	return events, err
}
