# Jane

Jane is a chatops bot written in Go and is completely configuration driven. Contributions are welcome via pull requests. Jane was started as a way of getting DevOps tasks and feedback into Slack. There are a billion other bots, but we wanted to learn Go, so this was a fun way to learn it and meet our needs. The name "Jane" was chosen by @kcwinner because he is a big fan of the _Ender's Game_ books. The name is not meant to be gender specific and can be effectively changed when you set your bot up.


## Running

### Install
Soon we will create a better path to install Jane, but for now, you can compile and use the startup scripts in the startup directory.


### Configuration
The configuration of Jane is via a json config file. The configuration file is expected to be named 'jane.config' and will be looked for in this order:
* --config config.json - Pass in a configuration file location as a command line parameter
* ./jane.json - the location of the jane binary
* ~/jane.json - the home directory of the user
* /etc/jane.json - the global config


## Connectors
Connectors are what Jane uses to pull in information, interpret it, and issue out a response. The Routes specify where the results from the input should be written to or * for all. The Target can specify a channel in the case of Slack. 

For the connector configuration, when adding routes, you must specify the ID of the connector you want to route the response to.

Supported connectors:
* [bamboo](#bamboo-connector) - Atlassian Bamboo integration
* [cli](#cli-connector) - Command line interface
* [email](#email-connector) - Email
* [exec](#exec-connector) - Execution of commands with monitoring capability
* imageme - Pull back images or animated gifs
* jira - Atlassian Jira integration
* response - Text Responses
* rss - RSS Feed
* slack - Slack chat
* twilio - send SMS alerts
* website - Monitor return code of websites
* [webhook](#webhook-connector) - Listener for webhooks
* wolfram - Execute queries against Wolfram Alpha


### Bamboo Connector

This connector was written to integrate bamboo builds. It was written against the bamboo cloud/hosted solution, but should be compatible with installed versions. This connector listens for and displays builds and deploys in addition to letting you execute the builds (but not deploys because the bamboo api is dated).

#### Example:

```
{"Type": "bamboo", "ID": "bamboo server", "Active": true, "Debug": true,
   "Server": "<URL>.atlassian.net", "Login": "<JIRA USER>", "Pass": "<JIRA PASS>",
   "Routes": [
       {"Match": "*", "Connectors": "*", "Target": "#devops"}
   ]
}
```

#### Usage:
* Run builds: `bamboo build <build key>`
* Get build status: `bamboo status <environment or build key>`
* Make sure to sepecify where you want the build and deploy messages to end up in the routes

#### Fields:
* _Type_ - This specifies the type of connector, in this case, 'exec2'
* _ID_ - This should be a unique identifier for this connector
* _Active_ - This is a boolean value to set this connector to be activated
* _Debug_ - This is a boolean value to set if the connector shows debug information in the logs
* _Server_ - The bamboo server address
* _Login_ - The bamboo user to login with
* _Pass_ - The bamboo password to connect with
* _Users_ - List of users who can execute the commands in this connector [security](#security)
* _Routes_ - One or more [routes](#routes)


### Cli Connector

This connector runs Jane via the command line interface instead of as a daemon and is helpful for debugging, or just indulging your command line love.

#### Example:

```
{"Type": "cli", "ID": "term-bot", "Active": true, "Debug": false,
  "Routes": [
    {"Match": "*", "Connectors": "term-bot", "Target": ""}
  ]
}
```

#### Fields:
* _Type_ - This specifies the type of connector, in this case, 'exec2'
* _ID_ - This should be a unique identifier for this connector
* _Active_ - This is a boolean value to set this connector to be activated
* _Debug_ - This is a boolean value to set if the connector shows debug information in the logs
* _Routes_ - One or more [routes](#routes)


### Email Connector

*NOTE - This is experimental and untested*
This connector allows for the sending of emails. Point the connector to a valid SMTP server.

#### Example:

```
{"Type": "email", "ID": "EmailServer", "Active": false,
  "Server": "smtp-server.myserver.com", "Port": "465",
  "Login": "smtpuser", "Pass": "smtppassword",
  "From": "jane@myserver.com"
}
```

#### Usage:
* Make sure to sepecify the to address in the target for routes to send emails

#### Fields:
* _Type_ - This specifies the type of connector, in this case, 'exec2'
* _ID_ - This should be a unique identifier for this connector
* _Active_ - This is a boolean value to set this connector to be activated
* _Debug_ - This is a boolean value to set if the connector shows debug information in the logs
* _Server_ - The SMTP server
* _Port_ - The port for the SMTP server
* _Login_ - The SMTP login
* _Pass_ - The SMTP password


### Exec Connector

This connector provides a single means of making local and remote calls to Linux systems. You can allow these calls to be made by command, but also mark the calls with the RunCheck property to set Jane to check them. This combined with the interpreter for output, makes it a very capable monitoring platform.

#### Example:

```
{"Type": "exec", "ID": "Elastic Search", "Active": true,
  "Server": "elasticsearch1.somecompany.com", "Port": "22", "Login": "jane", "Pass": "abc123",
  "Commands": [
    {
        "Name": "Apt Check",
        "Match": "jane elasticsearch1 aptcheck",
        "Output": "```%stdout% ```",
        "Cmd": "/usr/lib/nagios/plugins/check_apt",
        "Args": "",
        "HideHelp": false,
        "Help": "jane elasticsearch1 aptcheck - To check our elasticsearch!",
        "RunCheck": true,
        "Interval": 1,
        "Remind": 15,
        "Green": "*OK*",
        "Yellow": "*WARNING*",
        "Red": "*CRITICAL*"
    }
  ],
  "Routes": [
    {"Match": "*", "Connectors": "slack", "Target": "#devops"}
  ]
}
```

#### Usage:
* To make local calls to the system, leave out the Server, Port, Login, Pass values.

#### Fields:
* _Type_ - This specifies the type of connector, in this case, 'exec2'
* _ID_ - This should be a unique identifier for this connector
* _Active_ - This is a boolean value to set this connector to be activated
* _Debug_ - This is a boolean value to set if the connector shows debug information in the logs
* _Server_ - The server address or IP to connect to
* _Port_ - The port number to connect to (Default: 22)
* _Login_ - The user to login with
* _Pass_ - The password to connect with
* _Users_ - List of users who can execute the commands in this connector [security](#security)
* _Commands_ - One or more commands to execute against the defined server
  * _Name_ - Readable name of check
  * _Match_ - Command [match](#matching)
  * _Output_ - Formatting for the output of the command, use `%stdout%` as the output
  * _Cmd_ - The command to execute (do not include arguments)
  * _Args_ - The arguments, space deliminated (you can access anything after the match above with %msg%)
  * _HideHelp_ - A boolean to show or hide the help when displaying help (Default: false)
  * _Help_ - Optional help text, otherwise it'll show the Match value
  * _RunCheck_ - A boolean that will have Jane periodically run this (Default: false)
  * _Interval_ - An integer that is the number of minutes between checks when RunCheck is true (Default: 1)
  * _Remind_ - An integer which is the number of units of Interval to wait before reminding of a non-Green status, with Zero being no reminders (Default: 0)
  * _Green_ - A [match](#matching) to identify what is in a green state
  * _Yellow_ - A [match](#matching) to identify what is in a yellow state
  * _Red_ - A [match](#matching) to identify what is in a red state
* _Routes_ - One or more [routes](#routes)


### Webhook Connector

This connector opens a port for Jane to receive webhook calls. Webhooks calls are matched against the command list matches. Json can be interpreted and used to substitute into the output string. 


#### Example:

```
{"Type": "webhook", "ID": "Integrations", "Active": true, "Debug": true,
  "Port": "8080",
  "Commands": [
    {
        "Name": "Loggly Alerts",
        "Match": "/loggly/alerts",
        "Process": false,
        "Output": "```{alert_name} - {search_link}```",
        "Red": "*alert*"
    },
    {
        "Name": "Git Commits",
        "Match": "/git/commit",
        "Process": true,
        "Output": "jane build stuff"
    },
    {
        "Name": "Messages",
        "Match": "/messages",
        "Process": false,
        "Output": "{?}"
    }
  ],
  "Routes": [
      {"Match": "*", "Connectors": "*", "Target": "#devops"}
  ]
}
```


#### Fields:
* _Type_ - This specifies the type of connector, in this case, 'exec2'
* _ID_ - This should be a unique identifier for this connector
* _Active_ - This is a boolean value to set this connector to be activated
* _Debug_ - This is a boolean value to set if the connector shows debug information in the logs
* _Port_ - The port number to listen to (should be above 1024 if not running as a privledged user)
* _Users_ - List of users who can execute the commands in this connector [security](#security)
* _Commands_ - One or more commands to match the incoming webhook
  * _Name_ - Name of the matching webhook check
  * _Match_ - Webhook URL [match](#matching) (this will always be after the server name and port)
  * _Process_ - This defines if the incoming message should be processed by the other connector commands (true) or just published out to the routes (false) (Default: false)
  * _Output_ - This is the formatting for the output. Use the (https://github.com/Jeffail/gabs#parsing-and-searching-json)[json parsing rules] or '{}' to output the entire json payload or '{?}' to output the query string.
  * _Green_ - A [match](#matching) to identify what is in a green state
  * _Yellow_ - A [match](#matching) to identify what is in a yellow state
  * _Red_ - A [match](#matching) to identify what is in a red state
* _Routes_ - One or more [routes](#routes)


## Core Concepts


### Aliases

With Jane, you can create aliases for commands.

#### Example:

```
"Aliases": [
  {"Match": "jane build app", "Output": "bamboo build KEY-PLAN"},
  {"Match": "jane monitor prod", "Output": "jane monitor prod1 && jane monitor prod2 && jane monitor prod3"}
]
```

#### Fields:
* _Match_ - This will do a straight string match (does not support wildcards)
* _Output_ - This is the command that will be run when matched


### Routes

Routes can exist for connectors that listen to or interpret commands. Routes can have more than one connector if you would like to send messages to more than one place. Jane also matches on the routes to filter which messages get sent.

#### Example:

```
"Routes": [
  {"Match": "*", "Connectors": "slack", "Target": "#devops"},
  {"Match": "*DANGER*", "Connectors": "slack", "Target": "@matt"}
]
```

#### Usage:
* Some connector publishers allow you to specify a Target, such as Slack which uses a target for a channel
* Match follows the Jane [match rules](#matching)

#### Fields:
* _Match_ - This will match the message or any message with "*" using the [match](#matching)
* _Connectors_ - The connector name (ID) or "*" to match all connectors
* _Target_ - The target which is connector specific or "*" for all


### Matching

Jane uses a consistent string matching method throughout.

#### Examples:

`*failure*` - Match anywhere in a string

`failure*` - Match at the beginning of a string

`*failure` - Match at the end of a string

`/fail(.+)/` - Regular expression matching


### Security

The way of securing who can execute actions via Jane is by setting an optional list of users who are allowed to run commands on connectors that implement commands.

#### Example:

```
"Users": "matt,ken,joe"
```

#### Usage:
* This only applies to connectors that implement commands which users can execute
* The list of users is comma delimited
* The user name is dependant on the connector type, you can run "jane whoami" to get your name


## Architecture

Jane makes heavy use of the Go thread and channel features. Each connector can implement one of the three phases of the Jane bot messaging - Listeners, Commands and Publishers

### Listeners

Listeners are implemented to be long running tasks that take input externally. Examples of this are listening to an RSS feed or connecting to the Slack API. When a listener gets an event it is interested in, it creates a new message and passes it through the command messaging channel where it can be further acted on. Each message will be processed by the commands except when the message's process flag is set to false, in which case it passes through to the Publishers.

### Commands

Commands are implemented to act upon messages. They can do any task based on the message, make changes to the message and pass the message through to the Publishers. It is worth noting that each message gets potentially processed by all Commands, but it is up to the Commands to decide if it should act or ignore. Once the command is complete, it will pass the message to the publisher messaging channel where it will be processed by the Publishers.

### Publishers

Publishers are a means of communicating back out to the world. A publisher will take the message handed to it, format it, and send it through its implemented publish method.


## Getting Involved

### Development Environment
* Get your toes wet with Go
* Setup your Go 1.5.3 environment
* Pull the project with 'go get github.com/projectjane/jane'
* Compile with 'go install jane.go'
* Use the sample _jane.json_ file checked in as a starting point
* Run your code and config with `go run jane.go --config ~/jane.json`
