package state

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Store) PutTest(t TestState) error {
	_, err := s.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &s.table,
		Item: map[string]types.AttributeValue{
			"test_id":        &types.AttributeValueMemberS{Value: t.TestID},
			"status":         &types.AttributeValueMemberS{Value: t.Status},
			"profile_key":    &types.AttributeValueMemberS{Value: t.ProfileKey},
			"nodes":          &types.AttributeValueMemberN{Value: strconv.Itoa(t.Nodes)},
			"sessions":       &types.AttributeValueMemberN{Value: strconv.Itoa(t.Sessions)},
			"started_at":     &types.AttributeValueMemberN{Value: strconv.FormatInt(t.StartedAt, 10)},
			"expires_at":     &types.AttributeValueMemberN{Value: strconv.FormatInt(t.ExpiresAt, 10)},
			"min_rps":        &types.AttributeValueMemberN{Value: strconv.Itoa(t.MinRPS)},
			"max_rps":        &types.AttributeValueMemberN{Value: strconv.Itoa(t.MaxRPS)},
			"canary_percent": &types.AttributeValueMemberN{Value: strconv.Itoa(t.CanaryPercent)},
			"chaos_enabled":  &types.AttributeValueMemberBOOL{Value: t.ChaosConfig.Enabled},
			"chaos_latency":  &types.AttributeValueMemberN{Value: strconv.Itoa(t.ChaosConfig.LatencyMs)},
			"chaos_error":    &types.AttributeValueMemberN{Value: strconv.Itoa(t.ChaosConfig.ErrorRate)},
			"chaos_burst":    &types.AttributeValueMemberBOOL{Value: t.ChaosConfig.BurstPause},
		},
		if t.DesiredRPS != nil {
			item["desired_rps"] = &types.AttributeValueMemberN{
				Value: strconv.Itoa(*t.DesiredRPS),
			}
		}
	})
	return err
}
