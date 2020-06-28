SQL2Slack
=========
> a cron deamon that executes the specified sql query and forward its result to slack ()

Features
========
- Tiny & Portable.
- Works with multiple sql engine.
- Customize the slack message using javascript as well [underscore.js](https://underscorejs.org/).
- Cron like syntax for scheduling sql jobs.
- Uses [hcl language](https://github.com/hashicorp/hcl) for structured configurations.
- Ability to define a job per single file for future maintainability of large jobs.


Quick Overview
==============
```hcl
job tst {
    // slack-channel webhook url
    channel = "https://hooks.slack.com/services/xxxxxxxxxxxxxxxxxxxxx"

    // which sql driver do you use?
    driver = "mysql"

    // data source name (connection string)
    dsn = "root:root@tcp(127.0.0.1:3306)/spklvote"

    // the query this is a multiline example, you can just write the following
    // query = "select * from users"
    query = <<SQL
        SELECT users.* FROM users; -- WHERE created_at > DATE_SUB(NOW(), INTERVAL 2 HOUR)
    SQL

    // cron like syntax
    schedule = "* * * * *"

    // here you need to build the text message that will be sent 
    // to the slack channel.
    // ------- 
    // say(....): a function that will append the specified arguments into the main slack text message.
    // log(....): a logger function, for debugging purposes.
    // $rows: a variable that holds the output of the query execution.
    message = <<JS
        if ( $rows.length < 1 ) {
            return
        }

        say("there are (", $rows.length, ") new users!")
        say("users list is:")

        _.chain($rows).pluck('name').each(function(name){
            say("- ", name, " .")
        })
    JS
}
```

Integrating Slack
==================
1. Go [there](https://api.slack.com/apps).
2. Click on `Create New App`.
3. Choose `Incoming Webhooks` and activate it.
4. Scroll down to `Add New Webhook to Workspace` and follow the instructions.
5. Scroll down to the webhooks table, and copy the generated webhook url.`

Available SQL Drivers
=====================
| Driver | DSN |
---------| ------ |
| `mysql`| `usrname:password@tcp(server:port)/dbname?option1=value1&...`|
| `postgres`| `postgresql://username:password@server:port/dbname?option1=value1`|
| `sqlserver` | `sqlserver://username:password@host/instance?param1=value&param2=value` |
|             | `sqlserver://username:password@host:port?param1=value&param2=value`|
|             | `sqlserver://sa@localhost/SQLExpress?database=master&connection+timeout=30`|
| `mssql` | `server=localhost\\SQLExpress;user id=sa;database=master;app name=MyAppName`|
|         | `server=localhost;user id=sa;database=master;app name=MyAppName`|
|         | `odbc:server=localhost\\SQLExpress;user id=sa;database=master;app name=MyAppName` |
|         | `odbc:server=localhost;user id=sa;database=master;app name=MyAppName` |
| `clickhouse` |   `tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000` |

Installation
============
- [Binaries](/releases/)
- [Docker](https://hub.docker.com/r/alash3al/sql2slack)

Notes
=====
- by default `sql2slack` uses the current working directory as jobs files source, you can override that using `--jobs-dir` flag.
- each job file *must* have the `.s2s.hcl` suffix.
