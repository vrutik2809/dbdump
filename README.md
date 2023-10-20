<div align="center"><img src="./other/logo-name.png"  alt="logo" width=450px></div>
<div align="center">
    <img src="https://img.shields.io/github/issues/vrutik2809/dbdump"  alt="issue_badge">
    <img src="https://img.shields.io/github/issues-pr/vrutik2809/dbdump?logo=git"  alt="pr_badge">
    <img src="https://img.shields.io/github/actions/workflow/status/vrutik2809/dbdump/run_tests.yml?label=GitHub%20Workflows&logo=github"  alt="workflow_badge">
</div>
<br>

# DBDump

CLI tool for dumping various databases like `MongoDB`, `PostgreSQL`, `MySQL` into multiple formats like `json`, `tsv`, `csv`

## Installation Guide

<details open>
<summary>Linux</summary>

```bash
wget https://github.com/vrutik2809/dbdump/releases/download/latest/linux_amd64_dbdump

# make the binary executable
chmod +x linux_amd64_dbdump

# move the binary to the desired location (optional)
sudo mv linux_amd64_dbdump /usr/local/bin/dbdump
```

</details>
<details>
<summary>MacOS</summary>

```bash
curl https://github.com/vrutik2809/dbdump/releases/download/latest/darwin_amd64_dbdump

# make the binary executable
chmod +x darwin_amd64_dbdump

# move the binary to the desired location (optional)
sudo mv darwin_amd64_dbdump /usr/local/bin/dbdump
```

</details>

<details>
<summary>Windows</summary>

**Download the binary:** by clicking [here](https://github.com/vrutik2809/dbdump/releases/download/latest/windows_amd64_dbdump.exe)

</details>

## Usage Examples

```bash
# dump all mongodb collections to json
dbdump mongodb -u admin --password admin123 --host localhost -p 27017 -d test

# dump specific mongodb collections to bson
dbdump mongodb -u admin --password admin123 --host localhost -p 27017 -d test -c "users,posts" -o bson

# dump all postgresql tables to csv
dbdump pg -u postgres --password 123456 --host localhost -p 5432 -d test -o csv

# dump all postgresql tables excluding specific tables to json
dbdump pg -u postgres --password 123456 --host localhost -p 5432 -d test -e users,photos -o json

# dump all mysql tables to tsv
dbdump mysql -u root --password root --host localhost -p 3306 -d test -o tsv

```

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
