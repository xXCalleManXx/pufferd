/*
 Copyright 2016 Padduck, LLC

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

package templates

const Spigot = `{
  "pufferd": {
    "type": "java",
    "install": {
      "commands": [
        {
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar",
          "type": "download"
        },
        {
          "commands": [
            "java -jar BuildTools.jar --rev ${version}"
          ],
          "type": "command"
        },
        {
          "source": "spigot-*.jar",
          "target": "server.jar",
          "type": "move"
        }
      ]
    },
    "run": {
      "stop": "stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"-Xmx${memory}",
      	"-jar",
      	"server.jar"
      ],
      "program": "java"
    },
    "data": {
      "version": "1.10"
    }
  }
}`

const CraftbukkitBySpigot = `{
  "pufferd": {
    "type": "java",
    "install": {
      "commands": [
        {
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar",
          "type": "download"
        },
        {
          "commands": [
            "java -jar BuildTools.jar --rev ${version}"
          ],
          "type": "command"
        },
        {
          "source": "spigot-*.jar",
          "target": "server.jar",
          "type": "move"
        }
      ]
    },
    "run": {
      "stop": "stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"-Xmx${memory}",
      	"-jar",
      	"server.jar"
      ],
      "program": "java"
    },
    "data": {
      "version": "1.10"
    }
  }
}`

const Vanilla = `{
  "pufferd": {
    "type": "java",
    "install": {
      "commands": [
        {
          "type": "download",
          "files": "https://s3.amazonaws.com/Minecraft.Download/versions/${version}/minecraft_server.${version}.jar"
        }
      ]
    },
    "run": {
      "stop": "stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"-Xmx${memory}",
      	"-jar",
      	"minecraft_server.${version}.jar"
      ],
      "program": "java"
    },
    "data": {
      "version": "1.10"
    }
  }
}`
