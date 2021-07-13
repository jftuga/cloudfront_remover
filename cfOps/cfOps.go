package cfOps

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"strconv"
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

		if comment == "terraform--codershowcase.com" {
			fmt.Println("debug")
		}

		origin := "N/A"
		if len(obj.Origins.Items) > 0 {
			if obj.Origins.Items[0].S3OriginConfig != nil {
				origin = *obj.Origins.Items[0].S3OriginConfig.OriginAccessIdentity
			}
		}
		item := []string{*obj.Id, GetETag(*obj.Id), strconv.FormatBool(*obj.Enabled), *obj.Status, alias, origin, comment}
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
		data = append(data,item)
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
	//fmt.Println(result)
	return result.DistributionConfig
}

func DistIsEnabled(distributionId string) bool {
	conf := GetDistConf(distributionId)
	return *conf.Enabled
}
