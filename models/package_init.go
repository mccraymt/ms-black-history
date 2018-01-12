package models

import (
	"regexp"
	"strings"
)

// TODO: put this into config
const returnRatedLocations = false

// PostalCodeDict holds the whole body of the data
var PostalCodeDict map[string]*PostalCode
var countyDict map[string]*County

var statesRx = regexp.MustCompile("IL|IN|MD|TN|TX|VA")

func init() {
	PostalCodeDict = make(map[string]*PostalCode)
	countyDict = make(map[string]*County)

	for _, r := range dataArray {
		c := createCounty(&r)
		pc, found := PostalCodeDict[r.Zip]

		if found {
			updatePostalCode(&r, c, pc)
		} else {
			PostalCodeDict[r.Zip] = createPostalCode(&r, c)
		}
	}
}

func createCounty(r *GeoDataRecord) *County {
	c, found := countyDict[r.FIPS]
	if found {
		return c
	}

	var ctptr *string
	if len(r.CountyNameType) > 0 {
		ctptr = &(r.CountyNameType)
	}

	cty := County{
		Name:           titleCase(r.CountyName),
		FipsCode:       r.FIPS,
		State:          r.State,
		CountyNameType: ctptr,
	}

	countyDict[r.FIPS] = &cty
	return &cty
}

func createPostalCode(r *GeoDataRecord, c *County) *PostalCode {
	counties := []*County{c}
	var pcc *string
	var pcs *string
	var pcca *string
	if len(r.LL) > 0 {
		city := titleCase(r.City)
		state := strings.ToUpper(r.State)
		if len(r.CityAbb) > 0 {
			cabb := titleCase(r.CityAbb)
			pcca = &cabb
		}
		pcc = &city
		pcs = &state
	}

	rla := []*RatedLocation{}

	rtn := PostalCode{
		Value:             r.Zip,
		Counties:          counties,
		HasRatedLocations: isRatedLocation(r),
		City:              pcc,
		CityAbbreviation:  pcca,
		//RatedLocations:    rla,
		State:          pcs,
		PostalCodeType: pCodeType(r),
	}

	if returnRatedLocations {

		if isRatedLocation(r) {
			rl := RatedLocation{
				PostalCode:       r.Zip,
				State:            strings.ToUpper(r.State),
				County:           titleCase(r.CountyName),
				City:             titleCase(r.City),
				CityAbbreviation: titleCase(r.CityAbb),
			}

			rla = append(rla, &rl)
			rtn.RatedLocations = rla
		}
	}

	return &rtn
}

func updatePostalCode(r *GeoDataRecord, c *County, pc *PostalCode) {
	isrl := isRatedLocation(r)
	pc.HasRatedLocations = pc.HasRatedLocations || isrl

	if isrl && returnRatedLocations {
		// add rated location
		rl := RatedLocation{
			PostalCode:       r.Zip,
			State:            strings.ToUpper(r.State),
			County:           titleCase(r.CountyName),
			City:             titleCase(r.City),
			CityAbbreviation: titleCase(r.CityAbb),
		}

		rlfound := false
		for _, orl := range pc.RatedLocations {
			if (orl.City == rl.City || orl.CityAbbreviation == rl.City) && orl.County == rl.County && orl.State == rl.State {
				rlfound = true
			}

			if !rlfound && orl.City == rl.CityAbbreviation && orl.County == rl.County && orl.State == rl.State {
				orl.City = rl.City
				orl.CityAbbreviation = rl.CityAbbreviation
				rlfound = true
			}
		}
		if !rlfound {
			pc.RatedLocations = append(pc.RatedLocations, &rl)
		}
	}

	foundCty := false
	for _, cty := range pc.Counties {
		if cty.FipsCode == r.FIPS {
			foundCty = true
		}
	}

	if !foundCty {
		pc.Counties = append(pc.Counties, countyDict[r.FIPS])
	}

	if len(r.LL) > 0 {
		city := titleCase(r.City)
		state := strings.ToUpper(r.State)
		var pcca *string
		if len(r.CityAbb) > 0 {
			cabb := titleCase(r.CityAbb)
			pcca = &cabb
		}
		pc.City = &city
		pc.State = &state
		pc.PostalCodeType = pCodeType(r)
		pc.CityAbbreviation = pcca
	}
}

func isRatedLocation(r *GeoDataRecord) bool {
	if !statesRx.MatchString(strings.ToUpper(r.State)) {
		return false
	}

	if len(r.ZipType) > 0 {
		return false
	}

	return true
}

func pCodeType(r *GeoDataRecord) *string {
	pctype := r.ZipType
	var tptr *string
	if len(pctype) > 0 {
		tptr = &pctype
	}
	return tptr
}

func titleCase(t string) string {
	return strings.Title(strings.ToLower(t))
}
