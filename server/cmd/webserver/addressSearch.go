package main

import (
//    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
	"regexp"
	"strings"
//	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/vicesoftware/vice-go-boilerplate/cmd/webserver/models"
	"github.com/vicesoftware/vice-go-boilerplate/cmd/webserver/citysandiego"
	"github.com/vicesoftware/vice-go-boilerplate/pkg/logWriter"
	"go.uber.org/zap"
    // "github.com/vicesoftware/vice-go-boilerplate/cmd/webserver/middleware"
)

type ParcelSource int 
const (
	SR ParcelSource = 1 + iota
)

/* 
Manage the number of searches - needs to be done per account so will need to get 
added to the database as a column in the account
These are placeholder functions and will need to access the DB
*/
// var localSearchCount

// func iterateSearchCount(searchCount int64) int64{
// 	searchCount +=1
// 	return searchCount
// }

// func decrementSearchCount(searchCount int64) int64{
// 	if searchCount > 0 {
// 		searchCount -=1
// 	}
// 	return searchCount
// }

// func resetSearchCount(searchCount int64) int64{
// 	return searchCount = 0
// }

// func getSearchCount(UserAccount string) int64{
// 	return localSearchCount
// }

func isAPN(inputString string) bool {
	// Run through a regex to determine if it is an APN number
	apnRegex := "^[0-9]{10}$|^([0-9]{3}-){2}[0-9]{2}-[0-9]{2}$"
	matched, err := regexp.MatchString(apnRegex, inputString)
	if err != nil {
		logWriter.Error(err.Error())
	}
	return matched
}

/* get search property
* This function expects one of the following:
* APN
* Property Address
* It will return the Parcel information for the propery or an error message
* if an error occurred 
*/
func (ws *webserver) handleGetSearchProperty(w http.ResponseWriter, r *http.Request) error {
	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// get url params
	logWriter.Info("Entering the search handler...")
	vars := mux.Vars(r)
	searchText:=strings.TrimSpace(vars["searchText"])
	logWriter.Info("Searching for: ", zap.String("addr", searchText))

	var err error = nil
	var parcelInfo = make([]models.SR_Parcel_Struct,0)
	// var parcelInfoId models.SR_Parcel_Struct
	// var parcelId int = 0

	// if yes, call Scoutred on APN address API
	parcelInfo, err = getParcelInfoByAddress(SR,searchText)
	if err != nil {
		logWriter.Error("Failed to get parcel info by Address")	
		parcelInfo, err = getParcelInfoByGeo(SR,searchText)
		if err != nil {
			logWriter.Error("Failed to get address by Geolocation")
			return err
		}
	}
	// check the received structure and add the required array elements to store FAR
	if parcelInfo != nil && len(parcelInfo) > 0{
		logWriter.Debug("-----------Trying to add zoning struct enter 1----------------------")
		if parcelInfo[0].Zoning == nil || len(parcelInfo[0].Zoning) == 0 {
			logWriter.Debug("-----------Trying to add zoning struct enter 2----------------------")
			// zone := make([]models.Zoning, 0)
			var idZone int32 = 01
			zone1 := models.Zoning{ID: &idZone}
			parcelInfo[0].Zoning = append(parcelInfo[0].Zoning, zone1)
			logWriter.Debug("-----------Trying to add zoning struct enter 3----------------------")
			*parcelInfo[0].Zoning[0].ID = 01
			logWriter.Debug("-----------Trying to add zoning struct enter 4----------------------")
		}
		if parcelInfo[0].Zoning[0].Regulations == nil || len(parcelInfo[0].Zoning[0].Regulations) == 0 {
			logWriter.Debug("-----------Trying to add regs struct enter 1----------------------")
			regs := make([]models.ZoningRegulation, 0)
			logWriter.Debug("-----------Trying to add regs struct enter 2----------------------")
			var idTemp int32 = 01
			var baseTemp float32 = 0.0
			var noteTemp string = ""
			reg1 := models.ZoningRegulation{ID: &idTemp, FAR: models.FAR{Base: &baseTemp, Note: &noteTemp}}
			regs = append(regs, reg1)
			logWriter.Debug("-----------Trying to add regs struct enter 3----------------------")
			parcelInfo[0].Zoning[0].Regulations = append(parcelInfo[0].Zoning[0].Regulations, regs[0])
			logWriter.Debug("-----------Trying to add regs struct enter 4----------------------")
			fmt.Println("Parcel Regulations: ", parcelInfo[0].Zoning[0].Regulations)
			logWriter.Debug("-----------Trying to add regs struct leave ----------------------")
		}
	}
	var zoneDes string
	var area float64
	if parcelInfo != nil && len(parcelInfo) > 0 {
		logWriter.Debug("-----------Trying to add FAR info ----------------------")
		if parcelInfo[0].Zoning != nil && len(parcelInfo[0].Zoning) > 0 {
			logWriter.Debug("-----------Trying to add FAR info 2 ----------------------")
			if parcelInfo[0].Zoning[0].Designation != nil {
				zoneDes = *parcelInfo[0].Zoning[0].Designation
			} 
			logWriter.Debug("-----------Trying to add FAR info 3 ----------------------")
			if parcelInfo[0].GeomArea != nil {
				area = *parcelInfo[0].GeomArea
			}
			logWriter.Debug("-----------Trying to add FAR info 4 ----------------------")
			if parcelInfo[0].Zoning[0].Regulations != nil {
				logWriter.Debug("-----------Trying to add FAR info 5 ----------------------")
				oldNote := ""
				if parcelInfo[0].Zoning[0].Regulations[0].FAR.Note != nil {
					oldNote = *parcelInfo[0].Zoning[0].Regulations[0].FAR.Note
				}
				// logWriter.Debug("-----------Trying to add FAR info 6 ----------------------")
				// newNote := ""
				// var far *float32 = nil
				logWriter.Debug("********** for testing zone response remove for production *******************/")
				logWriter.Debug("Zoning information: ", zap.String("designation: ", zoneDes))
				logWriter.Debug("Zoning information: ", zap.Float64("area: ", area))
				logWriter.Debug("Zoning information: ", zap.String("area: ", oldNote))
				logWriter.Debug("/******************************************************************************/")
				far, newNote , isRes := citysandiego.GetFAR(zoneDes, area, -1)
				logWriter.Info("-----------Returned FAR info ----------------------")
				logWriter.Info("FAR calculation result.", zap.Any("FAR: ",far))
				logWriter.Info(newNote)
				logWriter.Info("------------------------------------------")
				parcelInfo[0].Zoning[0].Regulations[0].FAR.Base = far
				newNote = oldNote + "\n" + newNote
				parcelInfo[0].Zoning[0].Regulations[0].FAR.Note = &newNote
				var use string = "Residential"
				parcelInfo[0].Zoning[0].Use = &use
				if !isRes {	
					use = "Unknown"	
					parcelInfo[0].Zoning[0].Use = nil
				}			
				logWriter.Info(use)
			} else {
				logWriter.Debug("Regulations array is nil")			
			}

		} else {
			logWriter.Debug("Zoning array is nil")
		}
	} else {
		logWriter.Debug("Parcel Info array is nil")
	}
	

	// response, err := http.Get("https://scoutred.com/api/search/addresses?q=" + strconv.Itoa(id))
	// fmt.Println("https://scoutred.com/api/search/addresses?q=" + strconv.Itoa(id))
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }

	//fmt.Printf("%+v\n", parcelInfoId)

	// w.Header().Set("Content-Type", "application/json")	w.Write(parcelInfoJson)

	return Ok(w, parcelInfo)
}

func getParcelInfoByGeo(source ParcelSource, searchText string) ([]models.SR_Parcel_Struct, error){
	var googlRes models.GooglePlacesAPIResponse
	// var addrInfo = make([]models.SR_Address_API_Struct,0)
	var parcelId int = 0
	var err error = nil
	// var parcelInfoId models.SR_Parcel_Struct
	var parcelInfo = make([]models.SR_Parcel_Struct,0)
	switch source {
	case SR:
	// then the search Text is an address 
		// Put the search text through the Google Places API to get the 
		// lat and long 
		logWriter.Info("User entered ", zap.String("addr", searchText))
		googlRes, err = RequestAddressInfoGooglePlaces(searchText)
		if err != nil {
			logWriter.Error("Failed Google Placed API call")
			return nil, err
		}
		logWriter.Info("got results")
		if len(googlRes.Results) > 0 {
			logWriter.Info("Google Places Result", zap.Any("res: ", googlRes.Results[0]))
			// send long and lat to SR Parcel API
			//logWriter.Info("Lat = ")
			lat := googlRes.Results[0].Geometry.Location.Lat
			long := googlRes.Results[0].Geometry.Location.Lng
			//addr := googlRes.Results[0].FormattedAddress

			latstr:= strconv.FormatFloat(lat, 'f', 6, 64)
			longstr:= strconv.FormatFloat(long, 'f', 6, 64)


			logWriter.Info("Lat = " + latstr +", "+ "Lon = "+longstr)
			// logWriter.Info("Long = " + longstr)
			parcelInfo, err = RequestParcelInfoSRByGeo(latstr, longstr)
			logWriter.Debug("---------------------got here 1 -----------------------")
			if err != nil {
				logWriter.Error("Failed SR ParcelInfo API call")
				return nil, err
			} else if parcelInfo == nil || !(len(parcelInfo) > 0) {
				logWriter.Info("ScoutRed returned empty results for geo coordinates for ", zap.String("address: ", searchText))
				return nil, err
			} else {
				fmt.Printf("%v\n", parcelInfo)
				if parcelInfo[0].ID == nil{
					logWriter.Error("No ID for parcel info from Scoutred")
					return nil, err
				}
				parcelId = int(*parcelInfo[0].ID)
				logWriter.Info("Scoutred Parcel by Geo: ", zap.Int("result length: ", len(parcelInfo)))
				if parcelId != 0 {
					logWriter.Info("Scoutred Parcel by Geo: ", zap.Int("parcelId : ", parcelId))
				}
				// Not sure why I would do this here so commenting this out
				// parcelInfoId, err = RequestParcelInfoSRByID(strconv.Itoa(parcelId))
				// if err != nil {
				// 	logWriter.Error("Scoutred search by parcel id returned empty address array for %d ")
				// 	logWriter.Error(strconv.Itoa(addrs[0].ParcelID))
				// 	logWriter.Error(err.Error())
				// 	return err
				// }
				// fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
				// fmt.Println("%v\n", parcelInfoId)
				// Also commenting this out...
				// parcelInfo = append(parcelInfo, parcelInfoId)
				// fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
			}
		}
	}
	return parcelInfo, err
}

func getParcelInfoByAddress(source ParcelSource, searchText string) ([]models.SR_Parcel_Struct, error) {
	var addrInfo models.SR_Address_API_Struct
	// var parcelId int = 0
	var err error = nil
	var parcelInfoId models.SR_Parcel_Struct
	var parcelInfo = make([]models.SR_Parcel_Struct,0)
	switch source {
	case SR:
		// Need to get the parcel id and formatted address from Scoutred
		addrInfo, err = RequestAddressInfoSR(searchText)
		if err != nil {
			logWriter.Error("Error from RequestAddressInfoSR", zap.String("error:", err.Error()))
			return nil, err
		}
		parcelIdStr := fmt.Sprintf("%d", *addrInfo.ParcelID)
		addrIdStr := fmt.Sprintf("%d", *addrInfo.ID)
		logWriter.Info("Next search for parcel ", zap.String("parcelId:", parcelIdStr), zap.String("addressId:", addrIdStr))
		parcelInfoId, err = RequestParcelInfoSRByID(parcelIdStr, addrIdStr)

		if err != nil {
			logWriter.Error("Error from RequestParcelInfoSRByID", zap.String("error:", err.Error()))
			return nil, err
		}
		// fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		// fmt.Printf("%v\n", parcelInfoId)
		parcelInfo = append(parcelInfo, parcelInfoId)
		// fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		if err != nil {
			logWriter.Error("Failed SR ParcelInfo by API call", zap.String("error:", err.Error()))
			return nil, err
		}
	}
	
	return parcelInfo, err
}

/* get search property
* This function expects one of the following:
* APN
* Property Address
* It will return the Parcel information for the propery or an error message
* if an error occurred 
*/
func (ws *webserver) handleGetSearchPropertyByCoords(w http.ResponseWriter, r *http.Request) error {
	// get url params
	logWriter.Info("Entering the search by coords handler...")
	vars := mux.Vars(r)
	lat:= vars["lat"]
	lng:= vars["lon"]
	// logWriter.Info("URL sent lat = " + lat)
	// logWriter.Info("URL sent lngt = " + lat)
//	logWriter.Info("Searching for: " + searchText)

	var err error = nil
	var parcelInfo = make([]models.SR_Parcel_Struct, 0)
	//var parcelInfoJson []byte
	// var latstr, longstr string

	// then the search Text is an address 
	// Put the search text through the Google Places API to get the 
	// lat and long 

	if len(lat) == 0 {
		logWriter.Error(err.Error())
		fmt.Println("failed to convert lat %v to string", lat)
		return err
	}
	if len(lng) == 0 {
		logWriter.Error(err.Error())
		fmt.Println("failed to convert lng %v to string", lng)
		return err
	}

	logWriter.Info("Lat = " + lat)
	logWriter.Info("Long = " + lng)
	parcelInfo, err = RequestParcelInfoSRByGeo(lat, lng)
	if err != nil {
		logWriter.Error("Failed SR ParcelInfo API call", zap.String("error:", err.Error()))
		return err

	}
	// parcelInfoJson, err = json.Marshal(parcelInfo)
	if err != nil {
		logWriter.Error("Failed to Marshall parcelInfo: ", zap.String("error:", err.Error()))
		return err
	}
	// Once there is an address, set the state of the
	// address prop

	// fmt.Printf("%+v\n", string(parcelInfoJson))

	return Ok(w, parcelInfo)
}