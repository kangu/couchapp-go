# couchapp-go

couchapp-go is a command-line tool for parsing a folder structure and uploading its contents as a JSON document into a CouchDB database. This tool is designed to simplify the process of migrating a local folder structure into a CouchDB database, making it easier to manage and query the data using CouchDB's document-oriented database features.

## Why another CouchApp tool?

The reason behind creating couchapp-go is to address the limitations of existing options in the realm of CouchApp development.

Many existing tools rely on a complex web of dependencies, making installation and cross-platform compatibility a cumbersome task.

In contrast, couchapp-go has been designed with simplicity and convenience in mind, offering a streamlined and efficient solution. As a single, self-contained executable, it eliminates the need for managing various dependencies, simplifying the installation process and ensuring that developers can seamlessly work with CouchDB and deploy CouchApps across different platforms without the hassle of dealing with complex setups.

## Features

- Parses a local folder structure and its contents.
- Uploads the folder structure and its contents to a CouchDB database as a JSON document.
- Supports authentication with CouchDB to ensure data security.
- Supports live watching of the folder for changes and pushes updates to the database.

## Prerequisites

Before using couchapp-go, you should have the following prerequisites installed on your system:

- A CouchDB database: [Install CouchDB](https://couchdb.apache.org/)
- Go 1.21 (if you want to build from source): [Install Go](https://go.dev/)

## Installation

Pre-build packages are available for download in the Releases section.
The application itself consists of a single executable file that doesn't need any installation, just copy it somewhere on your system and add it to the PATH environment to access it from anywhere.

If you want to build it from source, use the following commands to output an executable for your own system architecture under `/bin`:

```bash
git clone https://github.com/kangu/couchapp-go
cd couchapp-go
go build -o bin/couchapp-go
```
You might want to use the filename `couchapp-go.exe` if you build on Windows.

## Usage

To use couchapp-go, follow these steps:

1. Open a terminal or command prompt.
2. Navigate to the folder you want to upload to CouchDB.
3. Run the following command, replacing the placeholders with your actual CouchDB credentials and database URL:

```bash
couchapp-go --db=test_db --user=username --pass=password
```

## Options

- `--db`: (Required) Target database where the design doc should be uploaded
- `--source`: (Optional) Folder path where the couchapp is located. Defaults to current folder
- `--user`: (Optional) Username for authentication
- `--password`: (Optional) Password for authentication
- `--watch`: (Optional) Watch the folder for changes and push to database on file updates

## Contributing

If you find any issues or have ideas for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the Apache License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- Thank you to the CouchDB community for their fantastic database system.
- Inspired by the need to simplify folder-to-database migrations, following on the footsteps of [couchapp](https://github.com/benoitc/couchapp), [Erica](https://github.com/benoitc/erica), [Kanso](https://github.com/kanso/kanso), [couchdb-push](https://github.com/jo/couchdb-push) and probably others.
