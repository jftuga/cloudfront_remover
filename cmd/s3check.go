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
	"github.com/jftuga/wipeout_cloudfront/cfOps"
	"github.com/spf13/cobra"
)

var searchOAIid string = ""
var searchRegion string = ""

// s3checkCmd represents the s3check command
var s3checkCmd = &cobra.Command{
	Use:   "s3check",
	Short: "Search for OAI in S3 bucket permissions",
	Run: func(cmd *cobra.Command, args []string) {
		bucketOAISearch(searchOAIid, searchRegion)
	},
}

func init() {
	rootCmd.AddCommand(s3checkCmd)
	s3checkCmd.Flags().StringVarP(&searchOAIid, "id", "i", "", "CloudFront OAI")
	s3checkCmd.Flags().StringVarP(&searchRegion, "region", "r", "", "AWS Region")
}

func bucketOAISearch(oaiId, searchRegion string) {
	cfOps.OAIBucketSearch(oaiId, searchRegion)
}
