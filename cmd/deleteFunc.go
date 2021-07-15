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

var deleteFuncName string = ""
var deleteFuncETag string = ""

// deleteFuncCmd represents the deleteFunc command
var deleteFuncCmd = &cobra.Command{
	Use:   "deleteFunc",
	Short: "Delete a CloudFront Function",
	Run: func(cmd *cobra.Command, args []string) {
		deleteFunction(deleteFuncName, deleteFuncETag)
	},
}

func init() {
	rootCmd.AddCommand(deleteFuncCmd)
	deleteFuncCmd.Flags().StringVarP(&deleteFuncName, "name", "n", "", "Function name")
	deleteFuncCmd.Flags().StringVarP(&deleteFuncETag, "etag", "e", "", "Function ETag")
	deleteFuncCmd.MarkFlagRequired("name")
	deleteFuncCmd.MarkFlagRequired("etag")
}

func deleteFunction(name, etag string) {
	fmt.Printf("Deleting function: %s [%s]\n", name, etag)
	cfOps.DeleteFunc(name, etag)
}
