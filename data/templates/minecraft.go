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
