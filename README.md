# DBDump

CLI tool for dumping various databases like `MongoDB`, `PostgreSQL`, `MySQL` into multiple formats like `json`, `tsv`, `csv`

## Commands

<!-- create a table  -->
| Command | Description |
| --- | --- |
|`mongodb`| command for dumping MongoDB database |
|`pg`| command for dumping PostgreSQL database |

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
    |`--output`|`-o`|output type `(json,bson,gzip)`|`json`|
    |`--help`|`-h`|help for the command|

