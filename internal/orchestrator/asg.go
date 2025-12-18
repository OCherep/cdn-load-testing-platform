package orchestrator

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

type ASGController struct {
	asgName string
	client  *autoscaling.Client
}

func NewASGController() *ASGController {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	return &ASGController{
		asgName: os.Getenv("AGENT_ASG_NAME"),
		client:  autoscaling.NewFromConfig(cfg),
	}
}

func (a *ASGController) SetDesired(cap int32) error {
	_, err := a.client.SetDesiredCapacity(context.Background(),
		&autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: &a.asgName,
			DesiredCapacity:      &cap,
			HonorCooldown:        false,
		})
	return err
}

func ParseDelta(curr int32, delta int) int32 {
	v := int(curr) + delta
	if v < 0 {
		return 0
	}
	return int32(v)
}
