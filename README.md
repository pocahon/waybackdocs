#waybackdocs

Features

    Historical Document Retrieval: Download snapshots of documents from the Internet Archive's Wayback Machine.
    File Filtering: Only downloads files with the extensions .doc, .docx, and .pdf (ignores .txt and .eml files).
    Concurrent Downloads: Uses a worker pool to process multiple downloads in parallel.
    Customizable Delay: Implements a 10-second delay before starting each download to ensure the website loads properly.
    Simple Installation: Easily installable using the Go toolchain.

Installation

Make sure you have Go installed, then install WaybackDocs with:

go install github.com/<your-username>/waybackdocs@latest

Replace <your-username> with your actual GitHub username.
Usage

After installation, you can run WaybackDocs from the command line:

waybackdocs -d example.com -n 10

    -d specifies the target domain (e.g., example.com).
    -n specifies the maximum number of downloads (use 0 for unlimited downloads).

Downloaded files are saved in the output directory.
Example

To download 7 documents from flevoland.nl, run:

waybackdocs -d flevoland.nl -n 7

Contributing

Contributions are welcome! Please feel free to open issues, suggest improvements, or submit pull requests.
License

This project is licensed under the MIT License. See the LICENSE file for details.
