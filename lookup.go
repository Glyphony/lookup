package lookup

import (
	"fmt"
	"io/ioutil"

	//I am NOT using the net package in this file, but the HTTP request package I needed
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
I'm using capitalized functions for most of this package as I believe these would be
functions that other tools would need to use to build their applications as well,
which makes sense to keep parity between the team's code

I'm also using CamelCase here as that is the preferred casing for GoLang
*/

//Would typically use the net package for this in a production environment
func VerifyIPv4Address(ip string) error {

	//Checking the user has used 3 dots in their string input
	if len(strings.Split(ip, dot)) != 4 {
		return fmt.Errorf("invalid IP: correct ipv4 address should consist of 3 periods")
	}

	//Splitting string into separate strings based on a '.' and making sure everything is a valid number
	for _, number := range strings.Split(ip, dot) {
		ipInt, err := strconv.Atoi(number)
		if err != nil {
			return fmt.Errorf("invalid IP: only numbers should be used in an ipv4 address")
		}
		if ipInt < 0 || ipInt > 255 {
			return fmt.Errorf("invalid IP: only valid numbers between 0-255 should be used for ipv4 addresses")
		}
		continue
	}
	return nil
}

func VerifyIPv6Address(ip string) error {
	return nil
}

//Determines whether the CONFIG_FILE_PATH env var is set and returns data as a string if it is, otherwise uses web URL
func ObtainIPAndASNData() (string, error) {
	envVar, exists := os.LookupEnv(lookupKey)
	if !exists {
		fmt.Println("Did not find local file, using web URL")
	}
	file, err := ioutil.ReadFile(envVar)
	if err != nil {
		fmt.Printf("ioutil.ReadFile(): %v", err)
		fmt.Println("\nDid not find local file, using web URL")
		return GetStaticIPDataFromWeb(ipLookupURL)
	}
	return string(file), nil
}

//This will grab the IP/ASN data from the web URL specified in consts.go
func GetStaticIPDataFromWeb(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return noValue, fmt.Errorf("http.Get(): %v", err)
	}
	strResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return noValue, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	return string(strResp), nil
}

//Takes the entire file of data and returns the individual lines as a string : int pair
func ParseTXTFileByReturnDelimiter(data string) (map[string]int, error) {
	ipLookupData := make(map[string]int)
	for _, line := range strings.Split(data, returnDelimiter) {
		if len(strings.Split(line, spaceDelimiter)) != 2 {
			continue
		}

		//We are making sure the 2nd value of each line is actually an ASN number
		asn, err := strconv.Atoi(strings.Split(line, spaceDelimiter)[1])
		if err != nil {
			return ipLookupData, fmt.Errorf("the lookup data file is corrupt - 2nd value should be an ASN number")
		}
		if asn < 0 || asn > 4294967295 {
			return ipLookupData, fmt.Errorf("the lookup data file is corrupt - asn number out of range")
		}
		ipLookupData[strings.Split(line, spaceDelimiter)[0]] = asn
	}
	return ipLookupData, nil
}

//May look somewhat confusing, but essentially taking the amount of IPs we know a subnet has available to it and testing if any of them match the IP
func getSubnetResult(prefixOctet string, bitMatch int, ipOctet string) bool {
	for offset := bitMatch; offset >= 0; offset-- {
		intVal, err := strconv.Atoi(prefixOctet)
		if err != nil {

			//We dont care about the error here as this function is not importable and we've already done our IP checks
			continue
		}
		if strconv.Itoa(intVal+offset) == ipOctet {
			return true
		}
	}
	return false
}

//This function is not importable as its lower case, so I'm not going to worry about error checking IP data again
//There might be a better way to figure out the matching subnets, but this was easy to read and the best method I could think of
func determineSubnetFit(prefix string, prefixLength int, ip string) bool {

	//This code takes the octet that fits within the mask length and figures out how many IPs are available in that subnet
	//and then runs getSubnetResult()
	if prefixLength == 32 && prefix == ip {
		return true
	}
	if prefixLength > 23 && prefixLength < 32 {
		if strings.Split(prefix, dot)[2] != strings.Split(ip, dot)[2] {
			return false
		}
		if strings.Split(prefix, dot)[1] != strings.Split(ip, dot)[1] {
			return false
		}
		if strings.Split(prefix, dot)[0] != strings.Split(ip, dot)[0] {
			return false
		}
		return getSubnetResult(strings.Split(strings.Split(prefix, dot)[3], forwardslash)[0], PowInt(2, (32-prefixLength))-1, strings.Split(ip, dot)[3])
	}
	if prefixLength > 15 && prefixLength < 24 {
		if strings.Split(prefix, dot)[1] != strings.Split(ip, dot)[1] {
			return false
		}
		if strings.Split(prefix, dot)[0] != strings.Split(ip, dot)[0] {
			return false
		}
		return getSubnetResult(strings.Split(prefix, dot)[2], PowInt(2, (24-prefixLength))-1, strings.Split(ip, dot)[2])
	}
	if prefixLength > 7 && prefixLength < 16 {
		if strings.Split(prefix, dot)[0] != strings.Split(ip, dot)[0] {
			return false
		}
		return getSubnetResult(strings.Split(prefix, dot)[1], PowInt(2, (16-prefixLength))-1, strings.Split(ip, dot)[1])
	}
	return false
}

//The error check may seem redundant, but if someone else uses just this function from this package they may want this error checking
func IsIPv4PartOfPrefix(prefix string, ip string) (bool, error) {
	err := VerifyIPv4Address(ip)
	if err != nil {
		return false, fmt.Errorf("VerifyIPv4Address(): %v", err)
	}
	err = VerifyIPv4Address(strings.Split(prefix, forwardslash)[0])
	if err != nil {
		return false, fmt.Errorf("VerifyIPv4Address(): %v", err)
	}
	prefixLength, err := strconv.Atoi(strings.Split(prefix, forwardslash)[1])
	if err != nil {
		return false, fmt.Errorf("prefix provided does not have a valid prefix length")
	}
	if prefixLength < 8 || prefixLength > 32 {
		return false, fmt.Errorf("prefix provided has an out-of-bounds prefix length")
	}
	return determineSubnetFit(prefix, prefixLength, ip), nil
}

//This is the magic sauce that returns all of the matching prefixes in a map[string]int
func findIPInResultsMap(data map[string]int, ip string) (map[string]int, error) {
	prefixMatchMap := map[string]int{}
	for prefix, asn := range data {

		//In this case I'm discarding the error since the data has already been validated
		//and this function cannot be imported to someone elses tool
		if ok, _ := IsIPv4PartOfPrefix(prefix, ip); !ok {
			continue
		}
		prefixMatchMap[prefix] = asn
	}
	if len(prefixMatchMap) == 0 {
		return prefixMatchMap, fmt.Errorf("no match was found")
	}
	return prefixMatchMap, nil
}

func printUserRequestedPrefixes(data OrderedData) {
	for index := range data.prefixes {
		fmt.Println(data.prefixes[index], data.asns[index])
	}
}

func IPLookup(ip string) error {
	err := VerifyIPv4Address(ip)
	if err != nil {

		//For troubleshooting the function path the error originated, I chain the function names into the error msg
		return fmt.Errorf("VerifyIPv4Address(): %v", err)
	}
	data, err := ObtainIPAndASNData()
	if err != nil {
		return fmt.Errorf("ObtainIPAndASNData(): %v", err)
	}
	resultsMap, err := ParseTXTFileByReturnDelimiter(data)
	if err != nil {
		return fmt.Errorf("ParseTXTFileByReturnDelimiter(): %v", err)
	}
	completedMap, err := findIPInResultsMap(resultsMap, ip)
	if err != nil {
		return fmt.Errorf("FindIPInResultsMap(): %v", err)
	}
	printUserRequestedPrefixes(SortBySpecificPrefix(completedMap, 32))

	return nil
}
