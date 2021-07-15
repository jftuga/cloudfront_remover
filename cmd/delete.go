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
	"strings"
	"time"

	"github.com/jftuga/cloudfront_remover/cfOps"
	"github.com/spf13/cobra"
)

var deleteDistID string = ""

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a CloudFront Distribution",
	Run: func(cmd *cobra.Command, args []string) {
		deleteCFDistribution(deleteDistID)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&deleteDistID, "id", "i", "", "Delete CloudFront Distribution ID")
	deleteCmd.MarkFlagRequired("id")
}

func deleteCFDistribution(distributionId string) {
	if cfOps.DistIsEnabled(distributionId) {
		fmt.Printf("Disabling distribution: %s\n", distributionId)
		cfOps.DisableDistribution(distributionId)
		time.Sleep(5 * time.Second)
	}

	fmt.Printf("Deleting distribution: %s\n", distributionId)
	deleteResult := ""
	for {
		deleteResult = cfOps.DeleteDistribution(distributionId)
		if strings.Contains(deleteResult, "trying to delete has not been disabled") {
			fmt.Println(deleteResult)
			fmt.Println("Will try again in 10 seconds...")
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	if deleteResult == "" {
		fmt.Printf("Distribution deleted: %s\n", distributionId)
	}
}
