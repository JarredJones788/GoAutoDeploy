
SETUP
1) Create a deploy key for the github repo.
2) Setup a git directory or use an existing one.
3) Setup a webhook to point to the server: http://SERVER-IP:CONFIG_PORT/github/push
4) setup config correctly. Basic Example config can be found in /cfg

Config Location On Server: /etc/autodeploy/config.json

CLI Service Commands
- autodeploy init (First time use, creates a systemd service)
- autodeploy start 
- autodeploy stop
- autodeploy restart
- autodeploy status
- autodeploy remove

GO Packages Used
- https://github.com/kardianos/osext
