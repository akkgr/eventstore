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
	lastEventTable string
	eventsTable    string
}

func NewDynamoDBClient(c context.Context, local bool) *DynamoDBClient {
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
		cfg, err := config.LoadDefaultConfig(c)
		if err != nil {
			panic(err)
		}
		client = dynamodb.NewFromConfig(cfg)
	}

	return &DynamoDBClient{
		store:          client,
		lastEventTable: "LastEvent",
		eventsTable:    "Events",
	}
}

func (dbc *DynamoDBClient) AppendLastEvent(e eventstore.Event, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(e)
	if err != nil {
		return err
	}

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.lastEventTable), Item: item,
	})

	return err
}

func (dbc *DynamoDBClient) UpdateLastEvent(e eventstore.Event, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(e)
	if err != nil {
		return err
	}

	version := strconv.Itoa(e.Version - 1)

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.lastEventTable), Item: item,
		ConditionExpression: aws.String("Version = :version"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":version": &types.AttributeValueMemberN{Value: version},
		},
	})

	return err
}

func (dbc *DynamoDBClient) AppendEvent(e eventstore.Event, ctx context.Context) error {
	item, err := attributevalue.MarshalMap(e)
	if err != nil {
		return err
	}

	_, err = dbc.store.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbc.eventsTable), Item: item,
	})

	return err
}

func (dbc *DynamoDBClient) GetLastEvent(id string, c context.Context) (eventstore.Event, error) {
	a := eventstore.Event{}

	marshaledId, err := attributevalue.Marshal(id)
	if err != nil {
		return a, err
	}

	key := map[string]types.AttributeValue{"Id": marshaledId}
	response, err := dbc.store.GetItem(c, &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(dbc.lastEventTable),
	})
	if err != nil {
		return a, err
	}

	item := eventstore.Event{}
	err = attributevalue.UnmarshalMap(response.Item, &item)
	if err != nil {
		return a, err
	}
	return item, err
}

func (dbc *DynamoDBClient) GetEvents(id string, v int, c context.Context) ([]eventstore.Event, error) {
	var err error
	var response *dynamodb.QueryOutput
	var events []eventstore.Event

	queryPaginator := dynamodb.NewQueryPaginator(dbc.store, &dynamodb.QueryInput{
		TableName:              aws.String(dbc.eventsTable),
		KeyConditionExpression: aws.String("#pk = :pk and #sk > :sk"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "Id",
			"#sk": "Version",
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
		var eventPage []eventstore.Event
		err = attributevalue.UnmarshalListOfMaps(response.Items, &eventPage)
		if err != nil {
			return events, err
		}
		events = append(events, eventPage...)
	}

	return events, err
}
