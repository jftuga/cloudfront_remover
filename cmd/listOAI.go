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
	"os"

	"github.com/jftuga/cloudfront_remover/cfOps"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listOAICmd represents the listOAI command
var listOAICmd = &cobra.Command{
	Use:   "listOAI",
	Short: "List CloudFront Origin Access Identities (OAIs)",
	Run: func(cmd *cobra.Command, args []string) {
		listOAIs()
	},
}

func init() {
	rootCmd.AddCommand(listOAICmd)
}

func listOAIs() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "ETAG", "COMMENT"})
	data := cfOps.GetOAIs()

	for _, entry := range data {
		table.Append(entry)
	}
	table.Render()
}
