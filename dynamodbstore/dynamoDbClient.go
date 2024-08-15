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
	snapshotsTable string
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
		snapshotsTable: "Snapshots",
	}
}

func (dbc *DynamoDBClient) AppendAggregate(a eventstore.Aggregate, ctx context.Context) error {
	item := aggregateFrom(a)
	dbitem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	version := strconv.Itoa(item.LastEvent.Version - 1)

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.aggregateTable), Item: dbitem,
		ConditionExpression: aws.String("version <> :version"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":version": &types.AttributeValueMemberN{Value: version},
		},
	})
	return err
}

func (dbc *DynamoDBClient) AppendSnapshot(s eventstore.Snapshot, ctx context.Context) error {
	item := snapshotFrom(s)
	dbitem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.snapshotsTable), Item: dbitem,
	})
	return err
}

func (dbc *DynamoDBClient) AppendEvent(e eventstore.Event, ctx context.Context) error {
	item := eventFrom(e)
	dbitem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.eventsTable), Item: dbitem,
	})
	return err
}

func (dbc *DynamoDBClient) GetAggregate(id eventstore.AggregateId, ctx context.Context) (eventstore.Aggregate, error) {
	a := eventstore.Aggregate{}

	marshaledId, err := attributevalue.Marshal(id)
	if err != nil {
		return a, err
	}

	key := map[string]types.AttributeValue{"id": marshaledId}
	response, err := dbc.store.GetItem(ctx, &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(dbc.aggregateTable),
	})
	if err != nil {
		return a, err
	}

	item := aggregate{}
	err = attributevalue.UnmarshalMap(response.Item, &item)
	if err != nil {
		return a, err
	}
	a = aggregateTo(item)
	return a, err
}

func (dbc *DynamoDBClient) GetSnapshot(id eventstore.AggregateId, ctx context.Context) (eventstore.Snapshot, error) {
	s := eventstore.Snapshot{}

	marshaledId, err := attributevalue.Marshal(id)
	if err != nil {
		return s, err
	}

	key := map[string]types.AttributeValue{"id": marshaledId}
	response, err := dbc.store.GetItem(ctx, &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(dbc.snapshotsTable),
	})
	if err != nil {
		return s, err
	}

	item := snapshot{}
	err = attributevalue.UnmarshalMap(response.Item, &item)
	if err != nil {
		return s, err
	}
	s = snapshotTo(item)
	return s, err
}

func (dbc *DynamoDBClient) GetEvents(id eventstore.AggregateId, v eventstore.EventNumber, ctx context.Context) ([]eventstore.Event, error) {
	var err error
	var response *dynamodb.QueryOutput
	var events []eventstore.Event

	queryPaginator := dynamodb.NewQueryPaginator(dbc.store, &dynamodb.QueryInput{
		TableName:              aws.String(dbc.eventsTable),
		KeyConditionExpression: aws.String("#pk = :pk and #sk > :sk"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "id",
			"#sk": "version",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: id},
			":sk": &types.AttributeValueMemberN{Value: strconv.Itoa(int(v))},
		},
	})
	for queryPaginator.HasMorePages() {
		response, err = queryPaginator.NextPage(context.TODO())
		if err != nil {
			return events, err
		}
		var eventPage []event
		err = attributevalue.UnmarshalListOfMaps(response.Items, &eventPage)
		if err != nil {
			return events, err
		}
		for _, e := range eventPage {
			events = append(events, eventTo(e))
		}
	}

	return events, err
}
