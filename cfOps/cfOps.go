package cfOps

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"strconv"
	"strings"
)

//This is a simple package to execute a few CloudFront Operations

// getDistributionData - For all distributions, return a 2D string array with each entry
// containing: ID, ETAG, 1ST ALIAS, 1ST OAI, COMMENT
func GetDistributionData() [][]string {
	var data [][]string
	svc := cloudfront.New(session.New())
	input := &cloudfront.ListDistributionsInput{}

	result, err := svc.ListDistributions(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeInvalidArgument:
				fmt.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return [][]string{}
	}

	for _, obj := range result.DistributionList.Items {
		alias := "N/A"
		if len(obj.Aliases.Items) > 0 {
			alias = *obj.Aliases.Items[0]
		}

		comment := "N/A"
		if len(*obj.Comment) > 0 {
			comment = *obj.Comment
		}

		origin := "N/A"
		if len(obj.Origins.Items) > 0 {
			if obj.Origins.Items[0].S3OriginConfig != nil {
				origin = *obj.Origins.Items[0].S3OriginConfig.OriginAccessIdentity
				origin = strings.Replace(origin,"origin-access-identity/cloudfront/", "", 1)
			}
		}
		item := []string{*obj.Id, GetETag(*obj.Id), strconv.FormatBool(*obj.Enabled), *obj.Status, GetACMCert(*obj.Id), alias, origin, comment}
		data = append(data, item)
	}
	return data
}

func GetOAIs() [][]string {
	var data [][]string
	svc := cloudfront.New(session.New())
	input := &cloudfront.ListCloudFrontOriginAccessIdentitiesInput{}

	result, err := svc.ListCloudFrontOriginAccessIdentities(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeInvalidArgument:
				fmt.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return [][]string{}
	}

	for _, obj := range result.CloudFrontOriginAccessIdentityList.Items {
		comment := "N/A"
		if len(*obj.Comment) > 0 {
			comment = *obj.Comment
		}
		item := []string{*obj.Id, GetOAIETag(*obj.Id), comment}
		data = append(data, item)
	}
	return data
}

func DeleteFunc(name, etag string) string {
	svc := cloudfront.New(session.New())
	input := cloudfront.DeleteFunctionInput{}
	input.Name = &name
	input.IfMatch = &etag

	_, err := svc.DeleteFunction(&input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// FIXME: these case stmts are wrong
			case cloudfront.ErrCodeInvalidArgument:
				fmt.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err.Error()
	}
	return ""
}

func GetFuncs() [][]string {
	var data [][]string
	svc := cloudfront.New(session.New())
	input := &cloudfront.ListFunctionsInput{}

	result, err := svc.ListFunctions(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// FIXME: these case stmts are wrong
			case cloudfront.ErrCodeInvalidArgument:
				fmt.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return [][]string{}
	}

	for _, obj := range result.FunctionList.Items {
		name := "N/A"
		if len(*obj.Name) > 0 {
			name = *obj.Name
		}
		stage := "N/A"
		if obj.FunctionMetadata.Stage != nil && len(*obj.FunctionMetadata.Stage) > 0 {
			stage = *obj.FunctionMetadata.Stage
		}

		var input cloudfront.DescribeFunctionInput
		input.Name = &name
		input.Stage = &stage

		result, err := svc.DescribeFunction(&input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case cloudfront.ErrCodeNoSuchCloudFrontOriginAccessIdentity:
					fmt.Println(cloudfront.ErrCodeNoSuchCloudFrontOriginAccessIdentity, aerr.Error())
				case cloudfront.ErrCodeAccessDenied:
					fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return [][]string{}
		}

		item := []string{name, stage, *result.ETag, result.FunctionSummary.FunctionMetadata.CreatedTime.String(), *result.FunctionSummary.FunctionConfig.Comment }
		data = append(data, item)
	}
	return data
}

func GetOAIETag(oaiId string) string {
	svc := cloudfront.New(session.New())
	input := &cloudfront.GetCloudFrontOriginAccessIdentityInput{}
	input.SetId(oaiId)

	result, err := svc.GetCloudFrontOriginAccessIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeNoSuchCloudFrontOriginAccessIdentity:
				fmt.Println(cloudfront.ErrCodeNoSuchCloudFrontOriginAccessIdentity, aerr.Error())
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "N/A"
	}
	return *result.ETag
}

// GetETag - given a CF Dist ID, return its ETag value
func GetETag(distributionId string) string {
	svc := cloudfront.New(session.New())
	input := &cloudfront.GetDistributionInput{}
	input.SetId(distributionId)

	result, err := svc.GetDistribution(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeNoSuchDistribution:
				fmt.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "N/A"
	}
	return *result.ETag
}

// GetACMCert - Check for a Viewer Certificate
func GetACMCert(distributionId string) string {
	svc := cloudfront.New(session.New())
	input := &cloudfront.GetDistributionInput{}
	input.SetId(distributionId)

	result, err := svc.GetDistribution(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeNoSuchDistribution:
				fmt.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "N/A"
	}
	cert := false
	if result.Distribution.DistributionConfig.ViewerCertificate != nil {
		cert = true
	}
	return strconv.FormatBool(cert)
}

func DeleteOAI(oaiId string) {
	svc := cloudfront.New(session.New())
	input := cloudfront.DeleteCloudFrontOriginAccessIdentityInput{}
	input.SetId(oaiId)
	etag := GetOAIETag(oaiId)
	input.IfMatch = &etag

	_, err := svc.DeleteCloudFrontOriginAccessIdentity(&input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			case cloudfront.ErrCodeInvalidIfMatchVersion:
				fmt.Println(cloudfront.ErrCodeInvalidIfMatchVersion, aerr.Error())
			case cloudfront.ErrCodeNoSuchCloudFrontOriginAccessIdentity:
				fmt.Println(cloudfront.ErrCodeNoSuchCloudFrontOriginAccessIdentity, aerr.Error())
			case cloudfront.ErrCodePreconditionFailed:
				fmt.Println(cloudfront.ErrCodePreconditionFailed, aerr.Error())
			case cloudfront.ErrCodeOriginAccessIdentityInUse:
				fmt.Println(cloudfront.ErrCodeOriginAccessIdentityInUse, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
}

func DeleteDistribution(distributionId string) string {
	input := cloudfront.DeleteDistributionInput{}
	input.Id = &distributionId
	etag := GetETag(distributionId)
	input.IfMatch = &etag

	svc := cloudfront.New(session.New())
	_, err := svc.DeleteDistribution(&input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			case cloudfront.ErrCodeDistributionNotDisabled:
				fmt.Println(cloudfront.ErrCodeDistributionNotDisabled, aerr.Error())
			case cloudfront.ErrCodeInvalidIfMatchVersion:
				fmt.Println(cloudfront.ErrCodeInvalidIfMatchVersion, aerr.Error())
			case cloudfront.ErrCodeNoSuchDistribution:
				fmt.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
			case cloudfront.ErrCodePreconditionFailed:
				fmt.Println(cloudfront.ErrCodePreconditionFailed, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err.Error()
	}
	return ""
}

func DisableDistribution(distributionId string) {
	svc := cloudfront.New(session.New())
	conf := GetDistConf(distributionId)
	input := &cloudfront.UpdateDistributionInput{}
	input.DistributionConfig = conf
	input.SetId(distributionId)
	value := false
	input.DistributionConfig.Enabled = &value
	etag := GetETag(distributionId)
	input.IfMatch = &etag

	_, err := svc.UpdateDistribution(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			case cloudfront.ErrCodeCNAMEAlreadyExists:
				fmt.Println(cloudfront.ErrCodeCNAMEAlreadyExists, aerr.Error())
			case cloudfront.ErrCodeIllegalUpdate:
				fmt.Println(cloudfront.ErrCodeIllegalUpdate, aerr.Error())
			case cloudfront.ErrCodeInvalidIfMatchVersion:
				fmt.Println(cloudfront.ErrCodeInvalidIfMatchVersion, aerr.Error())
			case cloudfront.ErrCodeMissingBody:
				fmt.Println(cloudfront.ErrCodeMissingBody, aerr.Error())
			case cloudfront.ErrCodeNoSuchDistribution:
				fmt.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
			case cloudfront.ErrCodePreconditionFailed:
				fmt.Println(cloudfront.ErrCodePreconditionFailed, aerr.Error())
			case cloudfront.ErrCodeTooManyDistributionCNAMEs:
				fmt.Println(cloudfront.ErrCodeTooManyDistributionCNAMEs, aerr.Error())
			case cloudfront.ErrCodeInvalidDefaultRootObject:
				fmt.Println(cloudfront.ErrCodeInvalidDefaultRootObject, aerr.Error())
			case cloudfront.ErrCodeInvalidRelativePath:
				fmt.Println(cloudfront.ErrCodeInvalidRelativePath, aerr.Error())
			case cloudfront.ErrCodeInvalidErrorCode:
				fmt.Println(cloudfront.ErrCodeInvalidErrorCode, aerr.Error())
			case cloudfront.ErrCodeInvalidResponseCode:
				fmt.Println(cloudfront.ErrCodeInvalidResponseCode, aerr.Error())
			case cloudfront.ErrCodeInvalidArgument:
				fmt.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			case cloudfront.ErrCodeInvalidOriginAccessIdentity:
				fmt.Println(cloudfront.ErrCodeInvalidOriginAccessIdentity, aerr.Error())
			case cloudfront.ErrCodeTooManyTrustedSigners:
				fmt.Println(cloudfront.ErrCodeTooManyTrustedSigners, aerr.Error())
			case cloudfront.ErrCodeTrustedSignerDoesNotExist:
				fmt.Println(cloudfront.ErrCodeTrustedSignerDoesNotExist, aerr.Error())
			case cloudfront.ErrCodeInvalidViewerCertificate:
				fmt.Println(cloudfront.ErrCodeInvalidViewerCertificate, aerr.Error())
			case cloudfront.ErrCodeInvalidMinimumProtocolVersion:
				fmt.Println(cloudfront.ErrCodeInvalidMinimumProtocolVersion, aerr.Error())
			case cloudfront.ErrCodeInvalidRequiredProtocol:
				fmt.Println(cloudfront.ErrCodeInvalidRequiredProtocol, aerr.Error())
			case cloudfront.ErrCodeNoSuchOrigin:
				fmt.Println(cloudfront.ErrCodeNoSuchOrigin, aerr.Error())
			case cloudfront.ErrCodeTooManyOrigins:
				fmt.Println(cloudfront.ErrCodeTooManyOrigins, aerr.Error())
			case cloudfront.ErrCodeTooManyCacheBehaviors:
				fmt.Println(cloudfront.ErrCodeTooManyCacheBehaviors, aerr.Error())
			case cloudfront.ErrCodeTooManyCookieNamesInWhiteList:
				fmt.Println(cloudfront.ErrCodeTooManyCookieNamesInWhiteList, aerr.Error())
			case cloudfront.ErrCodeInvalidForwardCookies:
				fmt.Println(cloudfront.ErrCodeInvalidForwardCookies, aerr.Error())
			case cloudfront.ErrCodeTooManyHeadersInForwardedValues:
				fmt.Println(cloudfront.ErrCodeTooManyHeadersInForwardedValues, aerr.Error())
			case cloudfront.ErrCodeInvalidHeadersForS3Origin:
				fmt.Println(cloudfront.ErrCodeInvalidHeadersForS3Origin, aerr.Error())
			case cloudfront.ErrCodeInconsistentQuantities:
				fmt.Println(cloudfront.ErrCodeInconsistentQuantities, aerr.Error())
			case cloudfront.ErrCodeTooManyCertificates:
				fmt.Println(cloudfront.ErrCodeTooManyCertificates, aerr.Error())
			case cloudfront.ErrCodeInvalidLocationCode:
				fmt.Println(cloudfront.ErrCodeInvalidLocationCode, aerr.Error())
			case cloudfront.ErrCodeInvalidGeoRestrictionParameter:
				fmt.Println(cloudfront.ErrCodeInvalidGeoRestrictionParameter, aerr.Error())
			case cloudfront.ErrCodeInvalidTTLOrder:
				fmt.Println(cloudfront.ErrCodeInvalidTTLOrder, aerr.Error())
			case cloudfront.ErrCodeInvalidWebACLId:
				fmt.Println(cloudfront.ErrCodeInvalidWebACLId, aerr.Error())
			case cloudfront.ErrCodeTooManyOriginCustomHeaders:
				fmt.Println(cloudfront.ErrCodeTooManyOriginCustomHeaders, aerr.Error())
			case cloudfront.ErrCodeTooManyQueryStringParameters:
				fmt.Println(cloudfront.ErrCodeTooManyQueryStringParameters, aerr.Error())
			case cloudfront.ErrCodeInvalidQueryStringParameters:
				fmt.Println(cloudfront.ErrCodeInvalidQueryStringParameters, aerr.Error())
			case cloudfront.ErrCodeTooManyDistributionsWithLambdaAssociations:
				fmt.Println(cloudfront.ErrCodeTooManyDistributionsWithLambdaAssociations, aerr.Error())
			case cloudfront.ErrCodeTooManyLambdaFunctionAssociations:
				fmt.Println(cloudfront.ErrCodeTooManyLambdaFunctionAssociations, aerr.Error())
			case cloudfront.ErrCodeInvalidLambdaFunctionAssociation:
				fmt.Println(cloudfront.ErrCodeInvalidLambdaFunctionAssociation, aerr.Error())
			case cloudfront.ErrCodeInvalidOriginReadTimeout:
				fmt.Println(cloudfront.ErrCodeInvalidOriginReadTimeout, aerr.Error())
			case cloudfront.ErrCodeInvalidOriginKeepaliveTimeout:
				fmt.Println(cloudfront.ErrCodeInvalidOriginKeepaliveTimeout, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
}

// GetDistConf- return a Distribution config for the given ID
func GetDistConf(distributionId string) *cloudfront.DistributionConfig {
	svc := cloudfront.New(session.New())
	input := &cloudfront.GetDistributionConfigInput{}
	input.SetId(distributionId)

	result, err := svc.GetDistributionConfig(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeNoSuchDistribution:
				fmt.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}
	return result.DistributionConfig
}

func DistIsEnabled(distributionId string) bool {
	conf := GetDistConf(distributionId)
	return *conf.Enabled
}

func GetRegionBuckets(searchRegion string) []string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(searchRegion),
	})
	svc := s3.New(sess)
	input := &s3.ListBucketsInput{}

	result, err := svc.ListBuckets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return []string{}
	}

	var regionBuckets []string
	for _, bucket := range result.Buckets {
		regionBuckets = append(regionBuckets, *bucket.Name)
	}
	return regionBuckets
}

func FindBucketRegion(bucketName, searchRegion string) string {
	sess := session.Must(session.NewSession())
	bucketRegion, err := s3manager.GetBucketRegion(context.Background(), sess, bucketName, searchRegion)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			fmt.Printf("unable to find bucket %s's region not found\n", bucketName)
		}
		fmt.Println(err.Error())
		return ""
	}
	return bucketRegion
}

func GetS3Policy(bucketName, searchRegion string) (string, string) {
	bucketRegion := FindBucketRegion(bucketName, searchRegion)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucketRegion),
	})

	svc := s3.New(sess)
	input := &s3.GetBucketPolicyInput{
		Bucket: aws.String(bucketName),
	}

	result, err := svc.GetBucketPolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				if !strings.Contains(aerr.Error(), "NoSuchBucketPolicy") {
					fmt.Println(aerr.Error())
				}
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", ""
	}
	return *result.Policy, bucketRegion
}
