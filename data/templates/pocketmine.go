/*
 Copyright 2017 Padduck, LLC

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

const Pocketmine = `{
  "pufferd": {
    "type": "pocketmine",
    "display": "PocketMine-MP",
    "install": {
      "commands": [
        {
          "type": "download",
          "files": "https://raw.githubusercontent.com/pmmp/php-build-scripts/master/installer.sh"
        },
        {
          "commands": [
            "chmod +x ./installer.sh",
            "./installer.sh"
          ],
          "type": "command"
        },
        {
          "type": "writefile",
          "text": "server-ip=${ip}\nserver-port=${port}\n",
          "target": "server.properties"
        }
      ]
    },
    "run": {
      "stop": "stop",
      "pre": [],
      "post": [],
      "arguments": [
      	"--no-wizard"
      ],
      "program": "./start.sh"
    },
    "environment": {
      "type": "tty"
    },
    "data": {
      "ip": {
        "value": "0.0.0.0",
        "required": true,
        "desc": "What IP to bind the server to",
        "display": "IP",
        "internal": false
      },
      "port": {
        "value": "19132",
        "required": true,
        "desc": "What port to bind the server to",
        "display": "Port",
        "internal": false,
        "type": "integer"
      }
    }
  }
}`
