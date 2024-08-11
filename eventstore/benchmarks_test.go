package eventstore_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/akkgr/eventstore/eventstore"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func BenchmarkStoreLoadEvents(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
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

	client := dynamodb.NewFromConfig(cfg)
	dbc := eventstore.NewDynamoDBClient(client)
	es := eventstore.NewEventStore(dbc, dbc)

	for i := 0; i < b.N; i++ {
		_, err := es.LoadEvents("123", context.Background())
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkStoreAppendEvents(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
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

	client := dynamodb.NewFromConfig(cfg)
	dbc := eventstore.NewDynamoDBClient(client)
	es := eventstore.NewEventStore(dbc, dbc)

	for i := 0; i < b.N; i++ {
		_ = es.Append(eventstore.Event{
			AggregateID:   "123",
			EventNumber:   i + 1,
			AggregateType: "Customer",
			EventName:     "CustomerUpdated",
			Created:       time.Now(),
			Data:          json.RawMessage(`{"name": "test test"}`),
		}, context.Background())
	}
}
