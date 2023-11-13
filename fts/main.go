package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/briandowns/spinner"
)

// ImageInfo representation
type ImageInfo struct {
	ImageDescription string `json:"ImageDescription,omitempty"`
	ImageName        string `json:"ImageName,omitempty"`
	ImageLocation    string `json:"ImageLocation,omitempty"`
	OwnerID          string `json:"OwnerID,omitempty"`
}

// AMI representation
type AMI struct {
	AmiID       string    `json:"AMI,omitempty"`
	Image       ImageInfo `json:"Image,omitempty"`
	InstanceIds []string  `json:"InstanceIds,omitempty"`
}

// AMIService handles AMI-related operations
type AMIService interface {
	GetAMIInfo(ctx context.Context, amiID string) (*ec2.DescribeImagesOutput, error)
	GatherAMIInfo(ctx context.Context, instances *ec2.DescribeInstancesOutput) ([]AMI, error)
}

// EC2Client handles EC2-related operations
type EC2Client struct {
	ec2Client *ec2.Client
}

// NewEC2Client creates a new EC2Client instance
func NewEC2Client(cfg aws.Config) *EC2Client {
	return &EC2Client{ec2Client: ec2.NewFromConfig(cfg)}
}

// GetInstances retrieves information about EC2 instances
func (c *EC2Client) GetInstances(ctx context.Context) (*ec2.DescribeInstancesOutput, error) {
	resp, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AMI implements AMIService for AMI-related operations
func (c *EC2Client) AMI() AMIService {
	return &AMIClient{ec2Client: c.ec2Client}
}

// AMIClient implements AMIService
type AMIClient struct {
	ec2Client *ec2.Client
}

// GetAMIInfo retrieves information about a specific AMI
func (c *AMIClient) GetAMIInfo(ctx context.Context, amiID string) (*ec2.DescribeImagesOutput, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []string{amiID},
	}

	result, err := c.ec2Client.DescribeImages(ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GatherAMIInfo gathers information about AMIs based on instances
func (c *AMIClient) GatherAMIInfo(ctx context.Context, instances *ec2.DescribeInstancesOutput) ([]AMI, error) {
	amiInfo := make(map[string]AMI)

	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			amiID := aws.ToString(instance.ImageId)

			if existingAMI, ok := amiInfo[amiID]; ok {
				existingAMI.InstanceIds = append(existingAMI.InstanceIds, *instance.InstanceId)
				amiInfo[amiID] = existingAMI
			} else {
				image, err := c.GetAMIInfo(ctx, amiID)
				if err != nil {
					return nil, err
				}

				imageInfo := ImageInfo{}
				if len(image.Images) > 0 {
					resultedImage := image.Images[0]
					imageInfo = ImageInfo{
						ImageDescription: aws.ToString(resultedImage.Description),
						ImageName:        aws.ToString(resultedImage.Name),
						ImageLocation:    aws.ToString(resultedImage.ImageLocation),
						OwnerID:          aws.ToString(resultedImage.OwnerId),
					}
				}
				amiInfo[amiID] = AMI{
					AmiID:       amiID,
					InstanceIds: []string{*instance.InstanceId},
					Image:       imageInfo,
				}
			}
		}
	}

	var amis []AMI
	for _, v := range amiInfo {
		amis = append(amis, v)
	}

	return amis, nil
}

// LoadAWSConfig loads AWS configuration
func LoadAWSConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, fmt.Errorf("error loading AWS configuration: %w", err)
	}
	return cfg, nil
}

// PrettyString returns a pretty-printed JSON representation of a value
func PrettyString(v interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

// runApp initializes and runs the application
func runApp(ctx context.Context) error {
	cfg, err := LoadAWSConfig(ctx)
	if err != nil {
		return err
	}

	ec2Client := NewEC2Client(cfg)
	instancesInfo, err := ec2Client.GetInstances(ctx)
	if err != nil {
		return fmt.Errorf("failed to get instances: %w", err)
	}

	if len(instancesInfo.Reservations) == 0 {
		return fmt.Errorf("no running instances were found")
	}

	amiInfo, err := ec2Client.AMI().GatherAMIInfo(ctx, instancesInfo)
	if err != nil {
		return fmt.Errorf("failed to gather AMI info: %w", err)
	}

	if len(amiInfo) == 0 {
		return fmt.Errorf("no Amazon Machine Images (AMIs) were found for the running instances")
	}

	amiInfoJSON, err := PrettyString(amiInfo)
	if err != nil {
		return fmt.Errorf("failed to pretty print AMI info: %w", err)
	}

	fmt.Println(amiInfoJSON)
	return nil
}

func main() {

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Prefix = "Gather information about all of the instances in the current region.: "

	s.Start() // Start the spinner
	if err := runApp(context.TODO()); err != nil {
		log.Fatal(err)
	}
	s.Stop() // Stop the spinner
}
