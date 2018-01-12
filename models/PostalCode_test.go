package models

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestVerifyAllPostalCodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestVerifyAllPostalCodes in short mode.")
	}

	fmt.Println("TestVerifyAllPostalCodes")
	startTime := time.Now()

	for thisCode, thisPc := range PostalCodeDict {
		if len(thisCode) != 5 {
			fmt.Printf("Postal code '%v' is invalid\r\n", thisCode)
			t.Fail()
		}

		if thisPc.Value != thisCode {
			fmt.Printf("Postal code value '%v' does not match index '%v'\r\n", thisPc.Value, thisCode)
			t.Fail()
		}

		if thisPc.City == nil || len(*thisPc.City) == 0 {
			fmt.Printf("Postal code city name '%v' invalid for index '%v'\r\n", thisPc.City, thisCode)
			t.Fail()
		}

		if len(thisPc.Counties) < 1 {
			fmt.Printf("No counties for index '%v'\r\n", thisCode)
			t.Fail()
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("TestVerifyAllPostalCodes completed in %v\n", elapsed)
}

func TestNonNullReturn(t *testing.T) {
	fmt.Println("TestNonNullReturn")
	pc := PostalCodeDict["01085"]
	if pc == nil {
		t.Fail()
	}
}

func TestNullReturn(t *testing.T) {
	fmt.Println("TestNullReturn")
	pc := PostalCodeDict["00000"]
	if pc != nil {
		t.Fail()
	}
}

func TestOutOfRatableStates(t *testing.T) {
	fmt.Println("TestOutOfRatableStates")
	pc := PostalCodeDict["41472"]
	if *pc.State != "KY" || pc.HasRatedLocations != false {
		t.Fail()
	}
}

func TestRatable(t *testing.T) {
	fmt.Println("TestRatable")
	pc := PostalCodeDict["22963"]
	if pc == nil || !pc.HasRatedLocations {
		t.Fail()
	}
}

func TestUniqueAndPOBoxCodes(t *testing.T) {
	fmt.Println("TestUniqueAndPOBoxCodes")
	// po box
	pc := PostalCodeDict["24142"]
	if strings.ToLower(*pc.PostalCodeType) != "po box" {
		fmt.Printf("Postal code type is '%v', should be 'PO BOX' for code '%v'", pc.PostalCodeType, pc.Value)
		t.Fail()
	}
	// unique
	pc2 := PostalCodeDict["23273"]
	if strings.ToLower(*pc2.PostalCodeType) != "unique" {
		fmt.Printf("Postal code type is '%v', should be 'UNIQUE' for code '%v'", pc2.PostalCodeType, pc2.Value)
		t.Fail()
	}
}
