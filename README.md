# WaybackDocs

WaybackDocs is a command-line tool written in Go that downloads historical document snapshots from the Wayback Machine for a given domain. It filters for specific file types (such as `.doc`, `.docx`, and `.pdf`) and downloads them concurrently using a worker pool to speed up the process.

## Features

- **Historical Document Retrieval:** Download snapshots of documents from the Internet Archive's Wayback Machine.
- **File Filtering:** Only downloads files with the extensions `.doc`, `.docx`, and `.pdf`.
- **Concurrent Downloads:** Uses a worker pool to process multiple downloads in parallel.
- **Customizable Delay:** Implements a 10-second delay before starting each download to ensure the website loads properly.
- **Simple Installation:** Easily installable using the Go toolchain.

## Installation

Make sure you have Go installed, then install WaybackDocs with:

```bash
go install github.com/pocahon/waybackdocs@latest
