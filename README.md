
SETUP
1) Create a deploy key for the github repo.
2) Setup a git directory or use an existing one.
3) Setup webhook to point to the server.
4) Configure config correctly.

Linux Config Location: /etc/GoAutoDeploy/config.json

Windows Config Location ./cfg/config.json


Basic Example config (JSON)
--------------

{
    "serverPort": "9000",
    "deployments": [
        {
            "repoName": "test_repo",
            "type": "repository",
            "repoLocation": "/website",
            "repoBranch": "main",
            "secret": "MY_SECRET_KEY",
            "commands": [
                {
                    "name": "pm2",
                    "args": ["restart", "server"]
                }
            ]
        }
    ]
}