job tst {
    channel = "https://enbs7oy1sjuy8.x.pipedream.net"

    driver = "mysql"

    dsn = "root:root@tcp(127.0.0.1:3306)/spklvote"

    query = <<SQL
        SELECT users.* FROM users
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