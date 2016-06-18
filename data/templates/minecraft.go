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

const Minecraft = `{
  "pufferd": {
    "type": "java",
    "install": {
      "commands": [
        {
          "type": "download",
          "files": "https://hub.spigotmc.org/BuildTools.jar"
        },
        {
          "type": "command",
          "commands": [
            "java -jar buildtools --rev ${version}",
            "mv spigot/spigot*.jar server.jar"
          ]
        }
      ],
      "windows": [
        {
          "type": "download",
          "files": "https://hub.spigotmc.org/BuildTools.jar"
        },
        {
          "type": "command",
          "commands": [
            "java -jar buildtools --rev ${version}",
            "move spigot/spigot*.jar server.jar"
          ]
        }
      ]
    },
    "run": {
      "stop": "/stop",
      "pre": [],
      "post": [],
      "arguments": "-Xmx1024M -jar server.jar"
    },
    "permissions": {
    }
  }
}`
