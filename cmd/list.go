/*
Copyright © 2021 John Taylor

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List distributions and their OAIs",

	Run: func(cmd *cobra.Command, args []string) {
		listAllDistributions(args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listAllDistributions(args []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Comment", "1st alias", "Id", "ETag", "1st OAI"})
	data := getDistributionSummary()
	for _, entry := range data {
		table.Append(entry)
	}
	table.Render()
}

func getDistributionSummary() [][]string {
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
			}
		}

		item := []string{comment, alias, *obj.Id, getETag(*obj.Id), origin}
		data = append(data, item)
	}
	return data
}

func getETag(distributionId string) string {
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
