# Email Domain Checker

## Overview

The Domain Details Checker is a Go program that enables you to inspect various details of domains, including MX (Mail Exchange), SPF (Sender Policy Framework), and DMARC (Domain-based Message Authentication, Reporting, and Conformance) records. This tool reads domain names from standard input, performs DNS lookups, and presents the results in an easy-to-read format.

## Pre-Requisites

Before using this program, ensure that you have Go installed on your system. You can download and install Go from the official website: https://golang.org/dl/

## Usage

1. Clone or download the program's source code.
2. Open a terminal and navigate to the directory containing the program.
3. Compile the program using the following command:
  ```go build do`main_details_checker.go```
4. Run the compiled program: ```./domain_details_checker [flags]```
```
Available flags:
    -input: Path to the input file containing domain names.
    -output: Path to the output file to write results.
    -workers: Number of concurrent workers for domain processing.
    -log: Path to the log file. If not provided, logs will be printed to the console.
    -verbose: Enable verbose logging.
```
5. Enter domain names one per line and press Enter. The program will check and display details for each domain.
6. To stop the program, press **Ctrl+C.**

## Output
The program will provide information in the following format for each domain:
```domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord```

    domain: The domain name being checked.
    hasMX: Indicates whether the domain has MX records (true/false).
    hasSPF: Indicates whether the domain has SPF records (true/false).
    spfRecord: The SPF record content, if present.
    hasDMARC: Indicates whether the domain has DMARC records (true/false).
    dmarcRecord: The DMARC record content, if present.

## Credits

This project was inspired by, and builds upon, code originally developed by educator [AkhilSharma90](https://github.com/AkhilSharma90). While the original project served as a foundation, this version introduces several enhancements:

- Concurrency: Implemented goroutines to process multiple domains concurrently.
- Command-Line Flags: Added CL-Flags for better usability.
- Improved Error-Handling: Enhanced error handling mechanism for better clarity and debugging.
- Code Refactoring: Code has been restructured and optimized for better readability and efficiency.
