package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EC2ClientTestSuite struct {
	suite.Suite
	stubber   *testtools.AwsmStubber
	ec2Client *EC2Client
}

func (suite *EC2ClientTestSuite) SetupTest() {
	suite.stubber = testtools.NewStubber()
	suite.ec2Client = NewEC2Client(*suite.stubber.SdkConfig)
}

func (suite *EC2ClientTestSuite) TearDownTest() {
	testtools.ExitTest(suite.stubber, suite.T())
}

func TestEC2ClientSuite(t *testing.T) {
	suite.Run(t, new(EC2ClientTestSuite))
}

func (suite *EC2ClientTestSuite) TestGetInstances() {
	suite.stubber.Add(testtools.Stub{
		OperationName: "DescribeInstances",
		Input:         &ec2.DescribeInstancesInput{},
		Output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							ImageId:          aws.String("fakeImageID"),
							InstanceId:       aws.String("fakeInstanceID"),
							InstanceType:     types.InstanceTypeA12xlarge,
							PrivateDnsName:   aws.String("fakePrivateDNS"),
							PrivateIpAddress: aws.String("fakePrivateIP"),
						},
					},
				},
			},
		},
	})

	// Act
	response, err := suite.ec2Client.GetInstances(context.TODO())

	// Assert
	assert.NoError(suite.T(), err, "GetInstances should not return an error")
	require.NotNil(suite.T(), response, "GetInstances response should not be nil")
	assert.Len(suite.T(), response.Reservations, 1, "Expected one reservation")
	assert.Len(suite.T(), response.Reservations[0].Instances, 1, "Expected one instance in the reservation")

	// Additional assertions for specific fields
	instance := response.Reservations[0].Instances[0]
	assert.Equal(suite.T(), "fakeImageID", aws.ToString(instance.ImageId), "Unexpected ImageId")
	assert.Equal(suite.T(), "fakeInstanceID", aws.ToString(instance.InstanceId), "Unexpected InstanceId")
	assert.Equal(suite.T(), types.InstanceTypeA12xlarge, instance.InstanceType, "Unexpected InstanceType")
	assert.Equal(suite.T(), "fakePrivateDNS", aws.ToString(instance.PrivateDnsName), "Unexpected PrivateDnsName")
	assert.Equal(suite.T(), "fakePrivateIP", aws.ToString(instance.PrivateIpAddress), "Unexpected PrivateIpAddress")
}

func (suite *EC2ClientTestSuite) TestGetAMIInfo() {
	amiID := "fakeAmiID"

	suite.stubber.Add(testtools.Stub{
		OperationName: "DescribeImages",
		Input:         &ec2.DescribeImagesInput{ImageIds: []string{amiID}},
		Output: &ec2.DescribeImagesOutput{
			Images: []types.Image{
				{
					Description:   aws.String("fakeImageDescription"),
					Name:          aws.String("fakeImageName"),
					ImageLocation: aws.String("fakeImageLocation"),
					OwnerId:       aws.String("fakeOwnerID"),
				},
			},
		},
	})

	// Act
	response, err := suite.ec2Client.AMI().GetAMIInfo(context.TODO(), amiID)

	// Assert
	assert.NoError(suite.T(), err, "GetAMIInfo should not return an error")
	require.NotNil(suite.T(), response, "GetAMIInfo response should not be nil")
	assert.Len(suite.T(), response.Images, 1, "Expected one image in the response")

	// Additional assertions for specific fields
	image := response.Images[0]
	assert.Equal(suite.T(), "fakeImageDescription", aws.ToString(image.Description), "Unexpected ImageDescription")
	assert.Equal(suite.T(), "fakeImageName", aws.ToString(image.Name), "Unexpected ImageName")
	assert.Equal(suite.T(), "fakeImageLocation", aws.ToString(image.ImageLocation), "Unexpected ImageLocation")
	assert.Equal(suite.T(), "fakeOwnerID", aws.ToString(image.OwnerId), "Unexpected OwnerID")
}

func (suite *EC2ClientTestSuite) TestGatherAMIInfo() {
	// Add a stub for DescribeInstances operation
	suite.stubber.Add(testtools.Stub{
		OperationName: "DescribeInstances",
		Input:         &ec2.DescribeInstancesInput{},
		Output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							ImageId:          aws.String("fakeAmiID"),
							InstanceId:       aws.String("fakeInstanceID"),
							InstanceType:     types.InstanceTypeA12xlarge,
							PrivateDnsName:   aws.String("fakePrivateDNS"),
							PrivateIpAddress: aws.String("fakePrivateIP"),
						},
					},
				},
			},
		},
	})

	// Add a stub for DescribeImages operation
	amiID := "fakeAmiID"
	suite.stubber.Add(testtools.Stub{
		OperationName: "DescribeImages",
		Input:         &ec2.DescribeImagesInput{ImageIds: []string{amiID}},
		Output: &ec2.DescribeImagesOutput{
			Images: []types.Image{
				{
					Description:   aws.String("fakeImageDescription"),
					Name:          aws.String("fakeImageName"),
					ImageLocation: aws.String("fakeImageLocation"),
					OwnerId:       aws.String("fakeOwnerID"),
				},
			},
		},
	})

	// Act
	instancesInfo, err := suite.ec2Client.GetInstances(context.TODO())
	require.NoError(suite.T(), err, "GetInstances should not return an error")

	response, err := suite.ec2Client.AMI().GatherAMIInfo(context.TODO(), instancesInfo)
	require.NoError(suite.T(), err, "GatherAMIInfo should not return an error")
	// Mocked AMIInfo for GetAMIInfo
	suite.TestGetAMIInfo()

	// Assertions
	assert.Len(suite.T(), response, 1, "Expected one AMI in the response")

	ami := response[0]
	assert.Equal(suite.T(), "fakeAmiID", ami.AmiID, "Unexpected AMI ID")
	assert.Len(suite.T(), ami.InstanceIds, 1, "Expected one instance ID in the AMI")

	image := ami.Image
	assert.Equal(suite.T(), "fakeImageDescription", image.ImageDescription, "Unexpected ImageDescription")
	assert.Equal(suite.T(), "fakeImageName", image.ImageName, "Unexpected ImageName")
	assert.Equal(suite.T(), "fakeImageLocation", image.ImageLocation, "Unexpected ImageLocation")
	assert.Equal(suite.T(), "fakeOwnerID", image.OwnerID, "Unexpected OwnerID")

}
