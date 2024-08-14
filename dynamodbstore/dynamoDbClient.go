package dynamodbstore

import (
	"context"
	"strconv"

	"github.com/akkgr/eventstore/eventstore"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	store          *dynamodb.Client
	aggregateTable string
	eventsTable    string
	SnapshotsTable string
}

func NewDynamoDBClient(ctx context.Context, local bool) *DynamoDBClient {
	var cfg aws.Config
	var client *dynamodb.Client
	var err error

	if local {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("localhost"),
			config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
					Source: "Mock credentials used above for local instance",
				},
			}),
		)

		if err != nil {
			panic(err)
		}

		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://localhost:8000")
		})
	} else {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			panic(err)
		}
		client = dynamodb.NewFromConfig(cfg)
	}

	return &DynamoDBClient{
		store:          client,
		aggregateTable: "Aggregates",
		eventsTable:    "Events",
		SnapshotsTable: "Snapshots",
	}
}

func (dbc *DynamoDBClient) AppendAggregate(aggregate eventstore.Aggregate, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(aggregate)
	if err != nil {
		return err
	}

	version := strconv.Itoa(aggregate.LastEvent.Version - 1)

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.aggregateTable), Item: item,
		ConditionExpression: aws.String("version <> :version"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":version": &types.AttributeValueMemberN{Value: version},
		},
	})
	return err
}

func (dbc *DynamoDBClient) AppendSnapshot(snapshot eventstore.Snapshot, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(snapshot)
	if err != nil {
		return err
	}

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.SnapshotsTable), Item: item,
	})
	return err
}

func (dbc *DynamoDBClient) AppendEvent(event eventstore.Event, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return err
	}

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.eventsTable), Item: item,
	})
	return err
}

func (dbc *DynamoDBClient) GetAggregate(id eventstore.AggregateId, ctx context.Context) (eventstore.Aggregate, error) {
	aggregate := eventstore.Aggregate{}

	marshaledId, err := attributevalue.Marshal(id)
	if err != nil {
		return aggregate, err
	}

	key := map[string]types.AttributeValue{"aggregateId": marshaledId}
	response, err := dbc.store.GetItem(ctx, &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(dbc.aggregateTable),
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

func (dbc *DynamoDBClient) GetEvents(id eventstore.AggregateId, ctx context.Context) ([]eventstore.Event, error) {
	var err error
	var response *dynamodb.QueryOutput
	var events []eventstore.Event

	queryPaginator := dynamodb.NewQueryPaginator(dbc.store, &dynamodb.QueryInput{
		TableName:              aws.String(dbc.eventsTable),
		KeyConditionExpression: aws.String("#pk = :pk"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "aggregateId",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: id},
		},
	})
	for queryPaginator.HasMorePages() {
		response, err = queryPaginator.NextPage(context.TODO())
		if err != nil {
			return events, err
		}
		var eventPage []eventstore.Event
		err = attributevalue.UnmarshalListOfMaps(response.Items, &eventPage)
		if err != nil {
			return events, err
		}
		events = append(events, eventPage...)
	}

	return events, err
}
