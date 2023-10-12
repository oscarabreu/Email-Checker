# Domain Details Checker

## Overview

This Go program allows you to check the details of domains, including MX (Mail Exchange), SPF (Sender Policy Framework), 
and DMARC (Domain-based Message Authentication, Reporting, and Conformance) records. It reads domain names from standard 
input, performs DNS lookups, and prints the results.

## Pre-Requisites

You need Go installed on your system to run this program.

## Usage
1. Clone or download the program's source code.
2. Open a terminal and navigate to the directory containing the program.
3. Compile the program using the following command:
  ```go build do`main_details_checker.go```
4. Run the compiled program: ```./domain_details_checker```
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
