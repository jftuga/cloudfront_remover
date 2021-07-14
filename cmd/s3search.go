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
	"github.com/jftuga/cloudfront_remover/cfOps"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var searchOAIid string = ""
var searchRegion string = ""

// s3searchCmd represents the s3search command
var s3searchCmd = &cobra.Command{
	Use:   "s3search",
	Short: "Search for an OAI in all S3 bucket policies",
	Run: func(cmd *cobra.Command, args []string) {
		bucketOAISearch(searchOAIid, searchRegion)
	},
}

func init() {
	rootCmd.AddCommand(s3searchCmd)
	s3searchCmd.Flags().StringVarP(&searchOAIid, "id", "i", "", "CloudFront OAI")
	s3searchCmd.Flags().StringVarP(&searchRegion, "region", "r", "", "AWS Region")
}

func outputResults(searchResults [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"BUCKET", "OAI"})
	for _, entry := range searchResults {
		table.Append(entry)
	}
	table.Render()
}

func bucketOAISearch(oaiId, searchRegion string) {
	allBuckets := cfOps.GetRegionBuckets(oaiId, searchRegion)
	var searchResults [][]string
	for _, bucket := range allBuckets {
		policy := cfOps.GetS3Policy(bucket, searchRegion)
		if strings.Contains(policy, oaiId) {
			item := []string{bucket, oaiId}
			searchResults = append(searchResults, item)
		}
	}
	if len(searchResults) > 0 {
		outputResults(searchResults)
	} else {
		fmt.Printf("OAI: %s was not found in the %d S3 Buckets that were searched.", oaiId, len(allBuckets))
	}
}
