package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord\n")

	for scanner.Scan() {
		inspectDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func inspectDomain(targetDomain string) {
	isMXPresent := validateMX(targetDomain)
	isSPFPresent, detectedSPF := detectSPF(targetDomain)
	isDMARCPresent, detectedDMARC := detectDMARC(targetDomain)

	fmt.Printf("%v, %v, %v, %q, %v, %q\n", targetDomain, isMXPresent, isSPFPresent, detectedSPF, isDMARCPresent, detectedDMARC)
}

func validateMX(domainName string) bool {
	mxEntries, mxErr := net.LookupMX(domainName)
	if mxErr != nil {
		log.Printf("MX Lookup Error: %v\n", mxErr)
		return false
	}
	return len(mxEntries) > 0
}

func detectSPF(domainName string) (bool, string) {
	txtEntries, txtErr := net.LookupTXT(domainName)
	if txtErr != nil {
		log.Printf("TXT Lookup Error: %v\n", txtErr)
		return false, ""
	}
	for _, entry := range txtEntries {
		if strings.HasPrefix(entry, "v=spf1") {
			return true, entry
		}
	}
	return false, ""
}

func detectDMARC(domainName string) (bool, string) {
	dmarcEntries, dmarcErr := net.LookupTXT("_dmarc." + domainName)
	if dmarcErr != nil {
		log.Printf("DMARC Lookup Error: %v\n", dmarcErr)
		return false, ""
	}
	for _, entry := range dmarcEntries {
		if strings.HasPrefix(entry, "v=DMARC1") {
			return true, entry
		}
	}
	return false, ""
}
