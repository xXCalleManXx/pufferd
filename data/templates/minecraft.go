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
    "display": "Spigot - Minecraft",
    "install": {
      "commands": [
        {
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar",
          "type": "download"
        },
        {
          "commands": [
            "java -jar BuildTools.jar --rev ${version}",
            "echo eula=${eula} > eula.txt"
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
      	"-Xmx${memory}M",
      	"-jar",
      	"server.jar"
      ],
      "program": "java"
    },
    "environment": {
      "type": "standard"
    },
    "data": {
      "version": {
      	"value": "1.10",
      	"required": true,
      	"desc": "Version of Minecraft to install",
      	"display": "Version",
      	"internal": false
      },
      "memory": {
      	"value": "1024",
      	"required": true,
      	"desc": "How much memory in MB to allocate to the Java Heap",
      	"display": "Memory (MB)",
      	"internal": false
      },
      "ip": {
        "value": "0.0.0.0",
        "required": true,
        "desc": "What IP to bind the server to",
        "display": "IP",
        "internal": false
      },
      "port": {
        "value": "25565",
        "required": true,
        "desc": "What port to bind the server to",
        "display": "Port",
        "internal": false
      },
      "eula": {
        "value": "false",
        "required": true,
        "desc": "Do you agree to the Minecraft EULA?",
        "display": "EULA Agreement (true/false)",
        "internal": false
      }
    }
  }
}`

const CraftbukkitBySpigot = `{
  "pufferd": {
    "type": "java",
    "display": "CraftBukkit by Spigot - Minecraft",
    "install": {
      "commands": [
        {
          "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar",
          "type": "download"
        },
        {
          "commands": [
            "java -jar BuildTools.jar --rev ${version}",
            "echo eula=${eula} > eula.txt"
          ],
          "type": "command"
        },
        {
          "source": "craftbukkit-*.jar",
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
      	"-Xmx${memory}M",
      	"-jar",
      	"server.jar"
      ],
      "program": "java"
    },
    "environment": {
      "type": "standard"
    },
    "data": {
      "version": {
      	"value": "1.10",
      	"required": true,
      	"desc": "Version of Minecraft to install",
      	"display": "Version",
      	"internal": false
      },
      "memory": {
      	"value": "1024",
      	"required": true,
      	"desc": "How much memory in MB to allocate to the Java Heap",
      	"display": "Memory (MB)",
      	"internal": false
      },
      "ip": {
        "value": "0.0.0.0",
        "required": true,
        "desc": "What IP to bind the server to",
        "display": "IP",
        "internal": false
      },
      "port": {
        "value": "25565",
        "required": true,
        "desc": "What port to bind the server to",
        "display": "Port",
        "internal": false
      },
      "eula": {
        "value": "false",
        "required": true,
        "desc": "Do you agree to the Minecraft EULA?",
        "display": "EULA Agreement (true/false)",
        "internal": false
      }
    }
  }
}`

const Vanilla = `{
  "pufferd": {
    "type": "java",
    "display": "Vanilla - Minecraft",
    "install": {
      "commands": [
        {
          "type": "download",
          "files": "https://s3.amazonaws.com/Minecraft.Download/versions/${version}/minecraft_server.${version}.jar"
        },
        {
          "source": "minecraft_server.*.jar",
          "target": "server.jar",
          "type": "move"
        },
        {
          "commands": [
            "echo eula=${eula} > eula.txt"
          ],
          "type": "command"
        },
      ]
    },
    "run": {
      "stop": "stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"-Xmx${memory}M",
      	"-jar",
      	"server.jar"
      ],
      "program": "java"
    },
    "environment": {
      "type": "standard"
    },
    "data": {
      "version": {
      	"value": "1.10",
      	"required": true,
      	"desc": "Version of Minecraft to install",
      	"display": "Version",
      	"internal": false
      },
      "memory": {
      	"value": "1024",
      	"required": true,
      	"desc": "How much memory in MB to allocate to the Java Heap",
      	"display": "Memory (MB)",
      	"internal": false
      },
      "ip": {
        "value": "0.0.0.0",
        "required": true,
        "desc": "What IP to bind the server to",
        "display": "IP",
        "internal": false
      },
      "port": {
        "value": "25565",
        "required": true,
        "desc": "What port to bind the server to",
        "display": "Port",
        "internal": false
      },
      "eula": {
        "value": "false",
        "required": true,
        "desc": "Do you agree to the Minecraft EULA?",
        "display": "EULA Agreement (true/false)",
        "internal": false
      }
    }
  }
}`

const Forge = `{
  "pufferd": {
    "type": "java",
    "display": "MinecraftForge - Minecraft",
    "install": {
      "commands": [
        {
          "type": "download",
          "files": "http://files.minecraftforge.net/maven/net/minecraftforge/forge/1.10.2-12.18.1.2011/forge-1.10.2-12.18.1.2011-installer.jar"
        },
        {
          "source": "forge-*.jar",
          "target": "installer.jar",
          "type": "move"
        },
        {
          "commands": [
            "java -jar installer.jar --installServer",
            "echo eula=${eula} > eula.txt"
          ],
          "type": "command"
        },
        {
          "source": "forge-*-universal.jar",
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
      	"-Xmx${memory}M",
      	"-jar",
      	"server.jar"
      ],
      "program": "java"
    },
    "environment": {
      "type": "standard"
    },
    "data": {
      "memory": {
      	"value": "1024",
      	"required": true,
      	"desc": "How much memory in MB to allocate to the Java Heap",
      	"display": "Memory (MB)",
      	"internal": false
      },
      "ip": {
        "value": "0.0.0.0",
        "required": true,
        "desc": "What IP to bind the server to",
        "display": "IP",
        "internal": false
      },
      "port": {
        "value": "25565",
        "required": true,
        "desc": "What port to bind the server to",
        "display": "Port",
        "internal": false
      },
      "eula": {
        "value": "false",
        "required": true,
        "desc": "Do you agree to the Minecraft EULA?",
        "display": "EULA Agreement (true/false)",
        "internal": false
      }
    }
  }
}`

const Sponge = `{
  "pufferd": {
    "type": "java",
    "display": "SpongeForge - Minecraft",
    "install": {
      "commands": [
        {
          "type": "download",
          "files": [
          	"http://files.minecraftforge.net/maven/net/minecraftforge/forge/1.10.2-12.18.1.2011/forge-1.10.2-12.18.1.2011-installer.jar",
          	"http://files.minecraftforge.net/maven/org/spongepowered/spongeforge/1.8.9-1890-4.2.0-BETA-1653/spongeforge-1.8.9-1890-4.2.0-BETA-1653.jar"
          ]
        },
        {
          "source": "forge-*.jar",
          "target": "installer.jar",
          "type": "move"
        },
        {
          "target": "mods",
          "type": "mkdir"
        },
        {
          "source": "spongeforge-*.jar",
          "target": "mods",
          "type": "move"
        },
        {
          "commands": [
            "java -jar installer.jar --installServer",
            "echo eula=${eula} > eula.txt"
          ],
          "type": "command"
        },
        {
          "source": "forge-*-universal.jar",
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
      	"-Xmx${memory}M",
      	"-jar",
      	"server.jar"
      ],
      "program": "java"
    },
    "environment": {
      "type": "standard"
    },
    "data": {
      "memory": {
      	"value": "1024",
      	"required": true,
      	"desc": "How much memory in MB to allocate to the Java Heap",
      	"display": "Memory (MB)",
      	"internal": false
      },
      "ip": {
        "value": "0.0.0.0",
        "required": true,
        "desc": "What IP to bind the server to",
        "display": "IP",
        "internal": false
      },
      "port": {
        "value": "25565",
        "required": true,
        "desc": "What port to bind the server to",
        "display": "Port",
        "internal": false
      },
      "eula": {
        "value": "false",
        "required": true,
        "desc": "Do you agree to the Minecraft EULA?",
        "display": "EULA Agreement (true/false)",
        "internal": false
      }
    }
  }
}`
