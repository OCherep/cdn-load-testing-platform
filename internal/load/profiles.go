package load

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Profile struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Sessions   int    `json:"sessions"`
	MinRPS     int    `json:"min_rps"`
	MaxRPS     int    `json:"max_rps"`
	StickyMode bool   `json:"sticky_mode"`
}

func LoadProfileFromS3(bucket, key string) (Profile, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	s3c := s3.NewFromConfig(cfg)

	obj, err := s3c.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return Profile{}, err
	}
	defer obj.Body.Close()

	data, _ := io.ReadAll(obj.Body)

	var p Profile
	json.Unmarshal(data, &p)
	return p, nil
}
