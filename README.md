# DBDump

<p align="center"><img src="./logo.png"  alt="logo" width="200" hight="200"></p>

CLI tool for dumping various databases like `MongoDB`, `PostgreSQL`, `MySQL` into multiple formats like `json`, `tsv`, `csv`

## Commands

<!-- create a table  -->
| Command | Description |
| --- | --- |
|`mongodb`| command for dumping MongoDB database |
|`pg`| command for dumping PostgreSQL database |
|`mysql`| command for dumping MySQL database |

## Command Options

- `mongodb`

    |flag|shorthand|description|default|
    |---|---|---|---|
    |`--username`|`-u`|username of the database|`''`|
    |`--password`||password of the database|`''`|
    |`--host`||host of the database|`localhost`|
    |`--port`|`-p`|port of the database|`0`|
    |`--db-name`|`-d`|database name||
    |`--dir`||directory to store the dump|`dump`|
    |`--srv`||use SRV connection format|`false`|
    |`--collections`|`-c`|name of the collections to dump|`[]`|
    |`--exclude-collections`|`-e`|name of the collections to exclude|`[]`|
    |`--output`|`-o`|output format of the dump `(json,bson,gzip)`|`json`|
    |`--help`|`-h`|help for the command|
- common for `pg` and `mysql`

    |flag|shorthand|description|default|
    |---|---|---|---|
    |`--username`|`-u`|username of the database|`postgres`|
    |`--password`||password of the database|`123456`|
    |`--host`||host of the database|`localhost`|
    |`--port`|`-p`|port of the database|`5432`|
    |`--db-name`|`-d`|database name||
    |`--dir`||directory to store the dump|`dump`|
    |`--tables`|`-t`|name of the tables to dump|`[]`|
    |`--exclude-tables`|`-e`|name of the tables to exclude|`[]`|
    |`--output`|`-o`|output format of the dump `(json,csv,tsv)`|`json`|
    |`--help`|`-h`|help for the command|

> You can always pass the `--help` flag to get help for the command
