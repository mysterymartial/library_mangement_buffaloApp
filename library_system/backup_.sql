pop v6.1.1

A tasty treat for all your database needs

Usage:
  buffalo-pop pop [flags]
  buffalo-pop pop [command]

Aliases:
  pop, db

Available Commands:
  create      Creates databases for you
  destroy     Allows to destroy generated code.
  drop        Drops databases for you
  fix         Brings pop, soda, and fizz files in line with the latest APIs
  generate    Generates config, model, and migrations files.
  migrate     Runs migrations against your database.
  reset       Drop, then recreate databases
  schema      Tools for working with your database schema

Flags:
  -c, --config string   The configuration file you would like to use.
  -d, --debug           Use debug/verbose mode
  -e, --env string      The environment you want to run migrations against. Will use $GO_ENV if set. (default "development")
  -h, --help            help for pop
  -p, --path string     Path to the migrations folder (default "./migrations")
  -v, --version         Show version information

Use "buffalo-pop pop [command] --help" for more information about a command.
