/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-tools/pkg/generate/rbac"
)

func main() {
	o := &rbac.ManifestOptions{}
	cmd := &cobra.Command{
		Use:   "rbac",
		Short: "Generates RBAC manifests",
		Long: `Generate RBAC manifests from the RBAC annotations in Go source files.
Usage:
# rbac generate [--name manager] [--input-dir input_dir] [--output-dir output_dir]
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := rbac.Generate(o); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("RBAC manifests generated under '%s' directory\n", o.OutputDir)
		},
	}

	registerFlags(cmd, o)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func registerFlags(cmd *cobra.Command, o *rbac.ManifestOptions) {
	f := cmd.Flags()
	f.StringVar(&o.Name, "name", "manager", "Name to be used as prefix in identifier for manifests")
	f.StringVar(&o.InputDir, "input-dir", "./pkg", "input directory pointing to Go source files")
	f.StringVar(&o.OutputDir, "output-dir", "./config", "output directory where generated manifests will be saved.")
}
