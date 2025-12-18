package orchestrator

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func StartAgents() error {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	c := ssm.NewFromConfig(cfg)

	_, err := c.SendCommand(context.TODO(), &ssm.SendCommandInput{
		DocumentName: aws.String("AWS-RunShellScript"),
		Targets: []types.Target{
			{
				Key:    aws.String("tag:Role"),
				Values: []string{"load-node"},
			},
		},
		Parameters: map[string][]string{
			"commands": {"docker restart agent"},
		},
	})
	return err
}

func BroadcastLimit(rps int) {
	cmd := fmt.Sprintf("export MAX_RPS=%d", rps)
	sendSSM(cmd)
}
