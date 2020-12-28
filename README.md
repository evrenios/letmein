# LETMEIN

![Default Slack Notifier](https://opensource.useinsider.com/letmein/images/letmein_image.png)

Letmein is a server-client coombination service where you can temporarly add a new IP to a pre-defined security groups for a custom time window.

When timer ends, added IP is automatically gets removed. 

When the program crash or exits, it also gracefully clean up temporary ip blocks so they never become stale. 

It also have a builtin slack notifier and a restricted ip blocks so you can prevent accidental ip additions/removals such as your corporate VPN or static home IP blocks.

![IP Restriction](https://opensource.useinsider.com/letmein/images/letmein_restricted.png)


You can download pre-built executables and run them immediatly or you can build it on your own to run on a k8 cluster or as a docker build.


| Variable | Example |
| ---- | -------- |
| LETMEIN_SECRET    | SUPERSECRETSECRET |
| SLACK_WEBHOOK    | https://hooks.slack.com/services/XXXXX/XXX |
| SLACK_CHANNEL    | ip-logs |
| SLACK_ICON_EMOJI    | :cookie: |
| SLACK_USERNAME    | WARDEN |
| BANLIST    | 192.168.2.1,127.0.0.1 |
| BANLIST    | us-east-1 |


## Installation
#### Server
* Mac 64-bit: https://opensource.useinsider.com/letmein/server/darwin_letmein_server
  ```
  $ curl https://opensource.useinsider.com/letmein/server/darwin_letmein_server > letmein_server && chmod +x letmein_server
  ```
* Linux 64-bit: https://opensource.useinsider.com/letmein/server/linux_letmein_server
* Windows 64-bit: https://opensource.useinsider.com/letmein/server/windows_letmein_server

#### Client
* Mac 64-bit: https://opensource.useinsider.com/letmein/server/darwin_letmein_client
  ```
  $ curl https://opensource.useinsider.com/letmein/server/darwin_letmein_client > letmein_client && chmod +x letmein_client
  ```
* Linux 64-bit: https://opensource.useinsider.com/letmein/server/linux_letmein_client
* Windows 64-bit: https://opensource.useinsider.com/letmein/server/windows_letmein_client