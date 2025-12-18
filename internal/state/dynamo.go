package state

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoStore struct {
	Table  string
	Client *dynamodb.Client
}

func NewDynamoStoreFromEnv() *DynamoStore {
	table := os.Getenv("STATE_TABLE")
	if table == "" {
		panic("STATE_TABLE env not set")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	return &DynamoStore{
		Table:  table,
		Client: dynamodb.NewFromConfig(cfg),
	}
}

func (d *DynamoStore) GetTest(ctx context.Context, id string) (TestState, error) {
	out, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.Table),
		Key: map[string]types.AttributeValue{
			"test_id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return TestState{}, err
	}

	i := out.Item

	state := TestState{
		TestID:    i["test_id"].(*types.AttributeValueMemberS).Value,
		Status:    i["status"].(*types.AttributeValueMemberS).Value,
		StartedAt: atoi64(i["started_at"]),
		ExpiresAt: atoi64(i["expires_at"]),
		MinRPS:    atoi(i["min_rps"]),
		MaxRPS:    atoi(i["max_rps"]),
		ChaosConfig: ChaosConfig{
			Enabled:   i["chaos_enabled"].(*types.AttributeValueMemberBOOL).Value,
			LatencyMs: atoi(i["chaos_latency"]),
			ErrorRate: atoi(i["chaos_error"]),
		},
	}

	if v, ok := i["desired_rps"]; ok {
		r := atoi(v)
		state.DesiredRPS = &r
	}

	return state, nil
}

func atoi(v types.AttributeValue) int {
	n, _ := strconv.Atoi(v.(*types.AttributeValueMemberN).Value)
	return n
}

func atoi64(v types.AttributeValue) int64 {
	n, _ := strconv.ParseInt(v.(*types.AttributeValueMemberN).Value, 10, 64)
	return n
}

func (s *DynamoStore) MarkSLAViolation(
	ctx context.Context,
	testID string,
	reason string,
) error {

	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.table),
		Key: map[string]types.AttributeValue{
			"test_id": &types.AttributeValueMemberS{Value: testID},
		},
		UpdateExpression: aws.String(
			"SET sla_violated = :v, violation_at = :t, violation_msg = :m",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberBOOL{Value: true},
			":t": &types.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().Unix(), 10)},
			":m": &types.AttributeValueMemberS{Value: reason},
		},
	})

	return err
}
