// Copyright © 2018 Urvil Patel <patelurvil38@gmail.com>
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

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urvil38/kubepaas/cmd"
	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/util"
)

func main() {
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	config.KubeConfig.ProjectRoot = projectRoot

	config.KubeConfig.KubepaasRoot = filepath.Join(projectRoot, "kubepaas")

	config.CLIConf = config.NewCLIConfig()

	if util.ConfigFileExists() {
		err := config.CLIConf.Read()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	cmd.Execute()
}
