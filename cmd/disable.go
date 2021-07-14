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
	"github.com/spf13/cobra"
)

var disableDistID string = ""

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a CloudFront Distribution",
	Run: func(cmd *cobra.Command, args []string) {
		disableCFDistribution(disableDistID)
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)
	disableCmd.Flags().StringVarP(&disableDistID, "id", "i", "", "Disable CloudFront Distribution ID")
	disableCmd.MarkFlagRequired("id")
}

func disableCFDistribution(distributionId string) {
	fmt.Printf("Disabling distribution: %s\n", distributionId)
	cfOps.DisableDistribution(distributionId)
}
