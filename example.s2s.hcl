job tst {
    channel = "https://hooks.slack.com/services/xxxxxxxxxxxxxxxxxxxxxxxxxxxx"

    driver = "mysql"

    dsn = "root:root@tcp(127.0.0.1:3306)/spklvote"

    query = <<SQL
        SELECT users.* FROM users; -- WHERE created_at > DATE_SUB(NOW(), INTERVAL 2 HOUR)
    SQL

    schedule = "* * * * *"

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