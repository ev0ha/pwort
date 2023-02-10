# pwort
Simple CLI password manager written in Go. You can create users and apps with passwords, display,<br>
update and delete them. See below for detailed usage.

# Disclamer

This is for my own educational purposes, so it goes without saying: Not intended for serious use!<br>
I coded this within a week while learning some Go and SQLite: There are no tests, style inconsistencies,<br>
possibly bugs. The master password is encrypted using bcrypt. I recently made some changes and plan to<br>
make some more in the future (the update command is not completely finished), so I decided to put the code on Github.

# Usage

- Create an user or an app

    **create** -u *USERNAME* -a *APP* [-unsafe]<br>
    Short command: **cr**<br>
    Create a new user with a new master password by only passing an username without specifying an app.<br>
    Users must be unique and a new one will not be created if you specify the same username.<br>
    By default the password must be at least 12 characters long, have lower- and uppercase letters and<br>
    at least one numeric and special character, respectively. Additionally pass the **unsafe** flag if<br>
    you don't want these restrictions. Create a new app by additionally passing an app, same **unsafe** rule applies.<br>
    Examples:<br>
    create -u Foo<br>
    create -u Foo -a Bar<br>
    cr -u Foo

- Display whether an user exists or the password of an app

    **show** -u *USERNAME* -a *APP*<br>
    Short command: **sh**<br>
    By only passing an user, checks if the user exists. If you pass an app, will ask for the master password<br>
    of the specified user and display the app password afterwards.<br>
    Examples:<br>
    show -u Foo<br>
    show -u Foo -a Bar<br>
    sh -u Foo

- Updating an user or an app

    **update** -u *USERNAME* -a *APP* [-unsafe]<br>
    Short command: **up**<br>
    If you want to update an user, only specify the user. For both, user and app, you will be asked whether you want<br>
    to change the name and/or the password. Pass **unsafe** flag if you don't want safety restrictions.<br>
    Examples:<br>
    update -u Foo<br>
    update -u Foo -a Bar<br>
    up -u Foo

- Deleting an user or an app

    **delete** -u *USERNAME* -a *APP*
    Short command: **dl**<br>
    Deletes either an user (and all their apps!) or just a specific app.<br>
    Examples:<br>
    delete -u Foo<br>
    delete -u Foo -a Bar<br>
    dl -u Foo

