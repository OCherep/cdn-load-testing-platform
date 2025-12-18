package state

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (s *Store) GetTest(testID string) (*TestState, error) {
	out, err := s.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &s.table,
		Key: map[string]types.AttributeValue{
			"test_id": &types.AttributeValueMemberS{Value: testID},
		},
	})
	if err != nil || out.Item == nil {
		return nil, err
	}

	return FromItem(out.Item), nil
}
