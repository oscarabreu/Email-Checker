package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	SPFRecordPrefix   = "v=spf1"
	DMARCRecordPrefix = "v=DMARC1"
	DMARCDomainPrefix = "_dmarc."
)

func main() {
	// We use a Scanner to buffer the input stream and read line by line
	// This prevents loading the entire input into memory at once!
	// However, if domain line is >64k, it will return false but this is rare (I hope!)
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	// Read eachline from the input stream and call our controller function

		
	for scanner.Scan() {
		domain := scanner.Text()
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			inspectDomain(d)
		}(domain)
	}
	// Error handling
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}

	wg.Wait()
}

// inspectDomain takes a domain name and prints a CSV line with the results
func inspectDomain(domain string) {
	var printMutex sync.Mutex


	isMXPresent := detectMX(domain) // Checks if MX records exist
	isSPFPresent, detectedSPF := detectSPF(domain) // Checks if SPF records exist and returns the record
	isDMARCPresent, detectedDMARC := detectDMARC(domain) // Checks if DMARC records exist and returns the record

	printMutex.Lock()
	fmt.Printf("%v, %v, %v, %q, %v, %q\n", domain, isMXPresent, isSPFPresent, detectedSPF, isDMARCPresent, detectedDMARC)
	printMutex.Unlock()}

// Function to validate MX records using LookupMX from the net package.
// Returns true if MX records exist (more than 0), false if not.
func detectMX(domainName string) bool {
	mxEntries, mxErr := net.LookupMX(domainName)
	// Error handling.
	if mxErr != nil {
		log.Printf("MX Lookup Error: %v\n", mxErr)
		return false
	}
	return len(mxEntries) > 0
}

// Function to validate SPF records using LookupTXT from the net package.
// Returns true if SPF records exist (more than 0), false if not.
func detectSPF(domainName string) (bool, string) {
	txtEntries, txtErr := net.LookupTXT(domainName)
	// Error handling.
	if txtErr != nil {
		log.Printf("TXT Lookup Error: %v\n", txtErr)
		return false, ""
	}
	// Loop through the TXT entries and check if any of them start with "v=spf1"
	for _, entry := range txtEntries {
		if strings.HasPrefix(entry, SPFRecordPrefix) {
			return true, entry
		}
	}
	return false, ""
}

// Function to validate DMARC records using LookupTXT from the net package.
// Returns true if SPF records exist (more than 0), false if not.
func detectDMARC(domainName string) (bool, string) {
	// To check the DMARC record for example.com, 
	// you'd look up the TXT records for _dmarc.example.com.
	dmarcEntries, dmarcErr := net.LookupTXT(DMARCDomainPrefix + domainName)
	// Error handling.
	if dmarcErr != nil {
		log.Printf("DMARC Lookup Error: %v\n", dmarcErr)
		return false, ""
	}
	// Loop through the TXT entries and check if any of them start with "v=DMARC1"
	for _, entry := range dmarcEntries {
		if strings.HasPrefix(entry, DMARCRecordPrefix) {
			return true, entry
		}
	}
	return false, ""
}
