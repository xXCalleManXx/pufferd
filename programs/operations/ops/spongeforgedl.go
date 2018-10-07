/*
 Copyright 2018 Padduck, LLC

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

package ops

import (
	"encoding/json"
	"errors"
	"github.com/pufferpanel/apufferi/common"
	"github.com/pufferpanel/pufferd/environments"
	"net/http"
	"os"
	"path"
)

const DOWNLOAD_API_URL = "https://dl-api.spongepowered.org/v1/org.spongepowered/spongeforge/downloads?type=stable"
const RECOMMENDED_API_URL = "https://dl-api.spongepowered.org/v1/org.spongepowered/spongeforge/downloads/recommended"
const FORGE_URL = "http://files.minecraftforge.net/maven/net/minecraftforge/forge/${minecraft}-${forge}/forge-${minecraft}-${forge}-installer.jar"

type SpongeForgeDl struct {
	ReleaseType string
}

type SpongeForgeDlOperationFactory struct {
}

func (of SpongeForgeDlOperationFactory) Key() string {
	return "spongeforgedl"
}

type download struct {
	Dependencies dependencies        `json:"dependencies"`
	Artifacts    map[string]artifact `json:"artifacts"`
}

type dependencies struct {
	Forge     string `json:"forge"`
	Minecraft string `json:"minecraft"`
}

type artifact struct {
	Url string `json:"url"`
}

func (op SpongeForgeDl) Run(env environments.Environment) error {

	var versionData download

	if op.ReleaseType == "latest" {
		client := &http.Client{}

		response, err := client.Get(DOWNLOAD_API_URL)
		if err != nil {
			return err
		}

		var all []download
		json.NewDecoder(response.Body).Decode(&all)
		response.Body.Close()

		versionData = all[0]
	} else {
		client := &http.Client{}

		response, err := client.Get(RECOMMENDED_API_URL)

		if err != nil {
			return err
		}

		json.NewDecoder(response.Body).Decode(&versionData)
		response.Body.Close()
	}

	if versionData.Artifacts == nil || len(versionData.Artifacts) == 0 {
		return errors.New("no artifacts found to download")
	}

	var versionMapping = make(map[string]interface{})
	versionMapping["forge"] = versionData.Dependencies.Forge
	versionMapping["minecraft"] = versionData.Dependencies.Minecraft

	err := downloadFile(common.ReplaceTokens(FORGE_URL, versionMapping), "forge-installer.jar", env)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(env.GetRootDirectory(), "mods"), 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = downloadFile(versionData.Artifacts[""].Url, path.Join("mods", "spongeforge.jar"), env)
	if err != nil {
		return err
	}

	return nil
}

func (of SpongeForgeDlOperationFactory) Create(op CreateOperation) Operation {
	releaseType := op.OperationArgs["releaseType"].(string)

	releaseType = common.ReplaceTokens(releaseType, op.DataMap)

	return SpongeForgeDl{ReleaseType: releaseType}
}