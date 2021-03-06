// Copyright © 2018 Giuseppe Maxia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/datacharmer/dbdeployer/sandbox"
	"github.com/spf13/cobra"
)

func FindTemplate(requested string) (group, contents string) {
	for name, tvar := range sandbox.AllTemplates {
		for k, v := range tvar {
			if k == requested {
				contents = v.Contents
				group = name
				return
			}
		}
	}
	fmt.Printf("template '%s' not found\n", requested)
	os.Exit(1)
	return
}

func ShowTemplate(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Argument required: template name")
		os.Exit(1)
	}
	requested := args[0]
	_, contents := FindTemplate(requested)
	fmt.Println(contents)
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a given template",
	Long:  ``,
	Run:   ShowTemplate,
}

func init() {
	templatesCmd.AddCommand(showCmd)

	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
