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
          "type": "download",
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
        },
        {
          "type": "command",
          "commands": [
            "java -jar BuildTools.jar --rev ${version}",
            "mv spigot-${version}.jar server.jar",
            "rm -rf apache-maven* BuildTools.jar BuildTools.log.txt Bukkit Spigot BuildData CraftBukkit work craftbukkit-*.jar"
          ]
        }
      ],
      "windows": [
        {
          "type": "download",
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
        },
        {
          "type": "command",
          "commands": [
            "java -jar BuildTools.jar --rev ${version}",
            "move spigot/spigot*.jar server.jar"
          ]
        }
      ]
    },
    "run": {
      "stop": "/stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"-Xmx${memory}",
      	"-jar",
      	"server.jar"
      ]
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
          "type": "download",
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
        },
        {
          "type": "command",
          "commands": [
            "java -jar BuildTools.jar --rev ${version}",
            "mv craftbukkit-${version}.jar server.jar",
            "rm -rf apache-maven* BuildTools.jar BuildTools.log.txt Bukkit Spigot BuildData CraftBukkit work spigot-*.jar"
          ]
        }
      ],
      "windows": [
        {
          "type": "download",
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
        },
        {
          "type": "command",
          "commands": [
            "java -jar BuildTools.jar --rev ${version}",
            "move spigot/spigot*.jar server.jar"
          ]
        }
      ]
    },
    "run": {
      "stop": "/stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"-Xmx${memory}",
      	"-jar",
      	"server.jar"
      ]
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
      "stop": "/stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"-Xmx${memory}",
      	"-jar",
      	"minecraft_server.${version}.jar"
      ]
    },
    "data": {
      "version": "1.10"
    }
  }
}`