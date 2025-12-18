package state

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

func FromItem(item map[string]types.AttributeValue) *TestState {
	return &TestState{
		TestID:     item["test_id"].(*types.AttributeValueMemberS).Value,
		Status:     item["status"].(*types.AttributeValueMemberS).Value,
		ProfileKey: item["profile_key"].(*types.AttributeValueMemberS).Value,
		ChaosConfig: ChaosConfig{
			Enabled:    item["chaos_enabled"].(*types.AttributeValueMemberBOOL).Value,
			LatencyMs:  atoi(item["chaos_latency"]),
			ErrorRate:  atoi(item["chaos_error"]),
			BurstPause: item["chaos_burst"].(*types.AttributeValueMemberBOOL).Value,
		},
	}
}
