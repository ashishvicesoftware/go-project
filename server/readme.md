# Vice Software Go API Boilerplate

This boilerplate was created to provided an opinonated starting point for projects that share our development values and priorities.

- [Vice Software Go API Boilerplate](#vice-software-go-api-boilerplate)
- [Getting Started](#getting-started)
	- [Installing Postgres](#installing-postgres)
		- [Starting the Database](#starting-the-database)
		- [Running PSQL](#running-psql)
	- [Setting Up the Database](#setting-up-the-database)
		- [Creating User](#creating-user)
		- [Creating the DB](#creating-the-db)
		- [Initializing the DB Schema](#initializing-the-db-schema)
		- [Seeding the Database](#seeding-the-database)
- [Installing Depedencies, Building and Running the App](#installing-depedencies-building-and-running-the-app)
- [Changing default configurations](#changing-default-configurations)
- [Our Values and Priorities](#our-values-and-priorities)

# Getting Started

1. Install Go
2. [Install PostGres Version 10](#installing-postgres)
3. [Setup the database](#setting-up-the-database)

## Installing Go
https://golang.org/doc/install#install

## Installing Postgres

On mac we recommend using [Postico](https://eggerapps.at/postico/docs/v1.5.6/) (GUI app for working with Postgres) and [Homebrew](https://brew.sh/) to install [Postgres 10](https://formulae.brew.sh/formula/postgresql@10#default) as shown below.

```
brew install postgresql@10
```

> Note: There is a [good setup guide located here](https://medium.freecodecamp.org/how-to-get-started-with-postgresql-9d3bc1dd1b11) but make sure you that you use the above command when installing Postgres as you need version 10.

To get post gres commandline setup after installing with brew execute

```
echo 'export PATH="/usr/local/opt/postgresql@10/bin:$PATH"' >> ~/.bash_profile
```

To test if you successfully installed Postgrew **open a new terminal windows** and execute the following command

```
postgres --version
```

If you see the following then Postgres is installed!

```
postgres (PostgreSQL) 10.6
```

### Starting the Database

If you installed using brew as described above then execute the following command to start Postgres version 10.

```
brew services start postgresql@10
```

### Stopping the Database

If you installed using brew and need to stop the database, use the following command

```
brew services stop postgresql
```

To verify that the DB started [run PSQL](#running-PSQL) as decribed below.

### Running PSQL

PSQL allows running PSQL commands. To run PSQL use the command below.

```
psql postgres
```

and you should see something like this

```
psql (10.6)
Type "help" for help.

postgres=#
```

This is the `psql` client. To exit type

```
/q
```

and hit enter.

## Setting Up the Database

After you have Postgress installed and available on the commandline issue the following commands to setup the user and database.

> Note: This guide uses default configurations in the API for DB Name, username, password, etc... see [Changing default configurations](#changing-default-configurations) below for details on how to change the defaults that are USED by the API.

### Creating User

```
createuser -P -e vice_boilerplate_user
```

Enter `vicesofware` as the password when prompted to use boilerplate defaults.

To verify the user was created execute `\du` in psql and you should see something like shown below.

```
vice_boilerplate=# \du
                                         List of roles
       Role name       |                         Attributes                         | Member of
-----------------------+------------------------------------------------------------+-----------
 ryan.vice             | Superuser, Create role, Create DB, Replication, Bypass RLS | {}
 vice_boilerplate_user |                                                            | {}
```

### Creating the DB

```
createdb -O vice_boilerplate_user vice_boilerplate
```

To verify the user was created execute `\l` in psql and you should see something like shown below.

```
vice_boilerplate=# \l
                                               List of databases
       Name       |         Owner         | Encoding |   Collate   |    Ctype    |      Access privileges
------------------+-----------------------+----------+-------------+-------------+-----------------------------
 postgres         | ryan.vice             | UTF8     | en_US.UTF-8 | en_US.UTF-8 |
 template0        | ryan.vice             | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/"ryan.vice"             +
                  |                       |          |             |             | "ryan.vice"=CTc/"ryan.vice"
 template1        | ryan.vice             | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/"ryan.vice"             +
                  |                       |          |             |             | "ryan.vice"=CTc/"ryan.vice"
 vice_boilerplate | vice_boilerplate_user | UTF8     | en_US.UTF-8 | en_US.UTF-8 |
```

### Initializing the DB Schema

To get the DB Schema needed for the boilerplate setup execute the SQL command in server/cmd/webserver/models/createTables.sql


# Installing Depedencies, Building and Running the App

Installing dependencies and Buidling can both be done by executing the command below

`go build`

from the `./cmd/webserver` directory.

To run the app simply execute

`./webserver`

from the `./cmd/webserver` directory.

Once you run the server you can see that it's working by opening the swagger page in the browser at http://127.0.0.1:8423/swagger/index.html

# Changing default configurations

If you want to configure a different databasename, username, password, etc... then make sure you update the database settings as described below:

### On a Mac:

Edit the .bash_profile in the user directory (cd ~)

	sudo nano ~/.bash_profile

### On Linux:

Edit the .bashrc file in the home directory (cd ~)

	sudo nano ~/.bashrc

	Add the follwing lines (or the values you used for your local database):
	
	export PD_HTTP_LISTENER="127.0.0.1:8423"
	export PD_DB_HOST="127.0.0.1"
	export PD_DB_PORT="5432"
	export PD_DB_USER="vice_boilerplate_user"
	export PD_DB_PASS="vicesoftware"
	export PD_DB_NAME="vice_boilerplate"
	export PD_DB_SSL_MODE="disable"
	export PD_APP_URL=http://localhost:3000/

### On Windows:

Edit the environment variables to set the same values 

If there are issues connecting to the database:

https://dba.stackexchange.com/questions/30818/problems-connecting-to-windows-based-postgresql

# Our Values and Priorities

Software is all about tradeoffs. The boilerplate for for projects and teams who:

1. Want to get started quickly with a GO API
2. Want to use [Gorrilla Mux](https://www.gorillatoolkit.org/pkg/mux) for their outing
3. Want to use [POSTGres](https://www.postgresql.org/about/news/1786/)
4. Like an ORM that provides some higher abstractions like [GORM](http://gorm.io/)
5. Prefer [Go Modules](https://github.com/golang/go/wiki/Modules)
6. Prefer a thin architecture stack without a lot of layers. Note that you can add layers as your app grows complex but this boilerplate provides a thin starting point
