package main

import (
	"Config"
	"OSM"
	"OSM/OSMHouses"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/bmizerany/pat"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var osm *dbsearch.Searcher
var dbh *dbsearch.Searcher
var config *Config.Config

const config_file string = "config/mapspapyrus.json"

type ReqParams struct {
	Id            int
	Limit         int
	Offset        int
	FiasId        string
	Params        map[string]interface{}
	Filters       map[string]string
	SearchFilters map[string]string
	ParamsPoi     map[string]interface{}
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func GetPoiList(w http.ResponseWriter, req *http.Request) {

	Rq := GetMetaParams(req)
	List, err_c := AddressPOI.Get(dbh, Rq.ParamsPoi, Rq.Limit, Rq.Offset)

	if err_c != nil {
		// SEND ERROR JSON
		ErrorFound(err_c, w, req)
		return
	}

	if List.IsEmpty() {
		NotFound(w, req)
		return
	}

	JsonSuccessPOI(List, w, req)
}

func GetOnePoi(w http.ResponseWriter, req *http.Request) {

	Rq := GetMetaParams(req)
	One, err_c := AddressPOI.GetOne(dbh, Rq.Id)

	if err_c != nil {
		// SEND ERROR JSON
		log.Printf("GetOnePoi: %s\n", err_c)
		NotFound(w, req)
		return
	}

	if One.IsEmpty() {
		NotFound(w, req)
		return
	}

	JsonSuccessPOI(One, w, req)
}

func GetAddresses(w http.ResponseWriter, req *http.Request) {

	Rq := GetMetaParams(req)
	List, err_c := Address.GetPlaces(dbh, config.SortTopIds(), Rq.Id, Rq.Limit, Rq.Offset, Rq.Filters, Rq.SearchFilters)
	if err_c != nil {
		ErrorFound(err_c, w, req)
		return
	}

	if List.IsEmpty() {
		NotFound(w, req)
		return
	}

	List.Places = Matching.MergeFromAddress(dbh, List.Places)
	List.Places = HandMade.MarkEdited(dbh, List.Places)

	JsonSuccess(List, w, req)
}

func get_parent_osm_id(One map[string]interface{}) int {
	var err error
	for {

		if p, find := One["polygon_osm_id"]; find {
			polygon_osm_id := iutils.AnyToInt(p)
			if polygon_osm_id != 0 {
				return polygon_osm_id
			}
		}

		if pid, find := One["parent_id"]; find {
			if parent_id := iutils.AnyToInt(pid); parent_id > 0 {
				One, err = Address.GetOneAddress(dbh, parent_id)
				if err != nil {
					log.Printf("get_parent_osm_id: %s\n", err)
				}
				continue
			}
		}

		return 0
	}
}

func houses_fias(Rq ReqParams) ([]map[string]interface{}, map[string]interface{}, map[string]interface{}, []map[string]interface{}, error) {

	List, err := Houses.GetList(dbh, Rq.Id)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	One, err_c := Address.GetOneAddress(dbh, Rq.Id)
	if err_c != nil {
		return nil, nil, nil, nil, err_c
	}

	FiasItem := map[string]interface{}{}
	FiasHousesList := []map[string]interface{}{}

	if fias_id, find := One["fias_id"]; find {
		FiasItemList, err_c := AddrObj.GetPlaces(dbh, map[string]interface{}{"aoguid": fias_id})
		if err_c != nil {
			return nil, nil, nil, nil, err_c
		}

		if len(FiasItemList) > 0 {
			FiasItem = FiasItemList[0]
			FiasHousesList, err = FiasHouses.GetPlaces(dbh, FiasItem["regioncode"].(string), map[string]interface{}{"aoguid": fias_id})
			if err != nil {
				return nil, nil, nil, nil, err
			}
		}

	}

	return List, One, FiasItem, FiasHousesList, nil
}

func GetHouses(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	List, One, FiasItem, FiasHousesList, err_h := houses_fias(Rq)

	if err_h != nil {
		ErrorFound(err_h, w, req)
		return
	}

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"request": req.URL.Path, // url
	}

	out := map[string]interface{}{
		"children":      List,
		"item":          One,
		"item_fias":     FiasItem,
		"children_fias": FiasHousesList,
		"meta":          meta,
	}

	JsonSend(out, w, req)
}

func GetHousesMap(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	List, One, FiasItem, FiasHousesList, err_f := houses_fias(Rq)

	if err_f != nil {
		ErrorFound(err_f, w, req)
		return
	}
	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"request": req.URL.Path, // url
	}

	polygon_osm_id := get_parent_osm_id(One)

	osm_houses, err_h := OSMHouses.GetHouses(osm, polygon_osm_id, One["name"].(string))
	if err_h != nil {
		ErrorFound(err_h, w, req)
		return
	}

	osm_street, err_st := OSMHouses.GetStreet(osm, polygon_osm_id, One["name"].(string))
	if err_st != nil {
		ErrorFound(err_st, w, req)
		return
	}

	center := OSM.New(nil)
	street := map[string]interface{}{}
	street["type"] = "FeatureCollection"
	features_st := []interface{}{}
	for _, v := range osm_street {
		one := map[string]interface{}{
			"type":     "Feature",
			"geometry": v["geojson"],
		}

		t := OSM.New(v["geojson"])
		center.Append(t)
		features_st = append(features_st, one)
	}
	street["features"] = features_st

	features := []interface{}{}
	features_fias := []interface{}{}
	features_osm := []interface{}{}

	for _, h_osm := range osm_houses {
		one := map[string]interface{}{
			"type":     "Feature",
			"geometry": h_osm["geojson"],
		}
		no_add := true
		osm_id := iutils.AnyToInt(h_osm["osm_id"])
		for _, h := range List {
			if iutils.AnyToInt(h["osm_id"]) != osm_id {
				continue
			}
			if iutils.AnyToString(h["fias_id"]) != "" {
				features_fias = append(features_fias, one)
			} else {
				features = append(features, one)
			}
			no_add = false
			break
		}

		if no_add {
			features_osm = append(features_osm, one)
		}
	}
	jsonhouse_osm := map[string]interface{}{}
	jsonhouse_osm["type"] = "FeatureCollection"
	jsonhouse_osm["features"] = features_osm

	jsonhouse := map[string]interface{}{}
	jsonhouse["type"] = "FeatureCollection"
	jsonhouse["features"] = features

	jsonhousefias := map[string]interface{}{}
	jsonhousefias["type"] = "FeatureCollection"
	jsonhousefias["features"] = features_fias

	maps_center, _ := center.Center()

	out := map[string]interface{}{
		"children":         List,
		"item":             One,
		"item_fias":        FiasItem,
		"children_fias":    FiasHousesList,
		"osm_houses":       osm_houses,
		"json_street":      street,
		"json_houses":      jsonhouse,
		"json_houses_fias": jsonhousefias,
		"json_houses_osm":  jsonhouse_osm,
		"maps_center":      maps_center,
		"meta":             meta,
	}

	JsonSend(out, w, req)
}

func SaveAddress(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	if Rq.Id == 0 {
		JsonSend(NotFoundBody("SaveAddress: нет переменной Id"), w, req)
		return
	}

	OneOld, err_c := Address.GetOneAddress(dbh, Rq.Id)
	if err_c != nil {
		ErrorFound(err_c, w, req)
		return
	}

	sdata := map[string]interface{}{}

	sdata["alt_names"] = iutils.AnyToString(Rq.Params["alt_names"])
	sdata["fias_id"] = iutils.AnyToString(Rq.Params["fias_id"])
	sdata["info"] = iutils.AnyToString(Rq.Params["info"])
	sdata["is_capital"] = iutils.AnyToBoolInt(Rq.Params["is_capital"])
	sdata["lat"] = iutils.AnyToString(Rq.Params["lat"])
	sdata["lon"] = iutils.AnyToString(Rq.Params["lon"])
	sdata["msg"] = iutils.AnyToString(Rq.Params["msg"])
	sdata["name"] = iutils.AnyToString(Rq.Params["name"])
	sdata["parent_id"] = iutils.AnyToInt(Rq.Params["parent_id"])
	sdata["path"] = iutils.AnyToString(Rq.Params["path"])
	sdata["point_osm_id"] = iutils.AnyToInt(Rq.Params["point_osm_id"])
	sdata["polygon_osm_id"] = iutils.AnyToInt(Rq.Params["polygon_osm_id"])

	if iutils.AnyToString(Rq.Params["type"]) != "" {
		sdata["type"] = iutils.AnyToString(Rq.Params["type"])
	}

	if iutils.AnyToString(Rq.Params["path"]) != "" {
		sdata["path"] = iutils.AnyToString(Rq.Params["path"])
	}

	sdata["group_type"] = iutils.AnyToString(Rq.Params["group_type"])
	sdata["place_type"] = iutils.AnyToString(Rq.Params["place_type"])
	sdata["fias_type_ru"] = iutils.AnyToString(Rq.Params["fias_type_ru"])
	sdata["official_status_ru"] = iutils.AnyToString(Rq.Params["official_status_ru"])

	Address.UpdateOne(dbh, Rq.Id, sdata, 0)

	OneNew, err_c2 := Address.GetOneAddress(dbh, Rq.Id)
	if err_c2 != nil {
		ErrorFound(err_c2, w, req)
		return
	}

	al := HandMade.New()
	al.AddParam("action", "update")
	al.AddParam("id", Rq.Id)
	al.AddParam("object_after", OneNew)
	al.AddParam("object_before", OneOld)
	al.AddParam("table_name", "address")
	al.AddParam("author", "papyrus-web")
	al.Save(dbh)

	rnd := fmt.Sprintf("/?t=%s", time.Now().Format("200601020304"))
	SendRedirect("/#/one_item_edit/"+iutils.AnyToString(Rq.Id)+rnd, w, req)
}

func SaveNewFIAS(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	if Rq.FiasId == "" || len(Rq.FiasId) != 36 {
		JsonSend(NotFoundBody("SaveNewFIAS: нет переменной FiasId"), w, req)
		return
	}

	if Rq.Id == 0 {
		JsonSend(NotFoundBody("SaveNewFIAS: нет переменной Id"), w, req)
		return

	}

	OneOld, err_c := Address.GetOneAddress(dbh, Rq.Id)
	if err_c != nil {
		log.Fatal("SaveNewFIAS: ", err_c)
	}

	if id, find := OneOld["id"]; !find || id != Rq.Id {
		JsonSend(NotFoundBody("SaveNewFIAS: ошибка нахождения"), w, req)
		return
	}

	Address.SetFiasId(dbh, Rq.Id, Rq.FiasId, 0)

	OneNew, err_c2 := Address.GetOneAddress(dbh, Rq.Id)
	if err_c2 != nil {
		ErrorFound(err_c2, w, req)
		return
	}

	al := HandMade.New()
	al.AddParam("action", "update")
	al.AddParam("id", Rq.Id)
	al.AddParam("object_after", OneNew)
	al.AddParam("object_before", OneOld)
	al.AddParam("table_name", "address")
	al.AddParam("author", "papyrus-web")
	al.Save(dbh)

	SendRedirect("/#/merge_fias/"+iutils.AnyToString(Rq.Id)+"?path="+Rq.SearchFilters["path"], w, req)
}

func GetAddress(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"fias_id": Rq.FiasId,
		"request": req.URL.Path, // url
	}

	One := map[string]interface{}{}
	var err_o error
	if Rq.Id > 0 {
		One, err_o = Address.GetOneAddress(dbh, Rq.Id)
		if err_o != nil {
			log.Printf("ERROR: %s\n", err_o)
		}
	} else if Rq.FiasId != "" {
		One, err_o = Address.GetOneFias(dbh, Rq.FiasId)
		if err_o != nil {
			log.Printf("ERROR: %s\n", err_o)
		}
	}

	out := map[string]interface{}{
		"item": One,
		"meta": meta,
	}

	JsonSend(out, w, req)
}

func GetAddressTree(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"fias_id": Rq.FiasId,
		"request": req.URL.Path, // url
	}

	One, FullLayers, err_c := Address.GetParentsLyers(dbh, Rq.Id, Rq.FiasId)
	if err_c != nil {
		ErrorFound(err_c, w, req)
		return
	}

	if One != nil {
		One["get_id"] = One["id"]
	}

	out := map[string]interface{}{
		"layers": FullLayers,
		"item":   One,
		"meta":   meta,
	}

	JsonSend(out, w, req)
}

func CountFiasMatching(w http.ResponseWriter, req *http.Request) {
	CountHouses, CountFiasHouses, err := Houses.CountFias(dbh)
	if err != nil {
		ErrorFound(err, w, req)
		return
	}

	CountAddress, CountFiasAddress, err2 := Address.CountFias(dbh)
	if err != nil {
		ErrorFound(err2, w, req)
		return
	}

	out := map[string]interface{}{
		"what": "count_fias_house",
		"house": map[string]interface{}{
			"total": CountHouses,
			"fias":  CountFiasHouses,
		},
		"address": map[string]interface{}{
			"total": CountAddress,
			"fias":  CountFiasAddress,
		},
	}

	JsonSend(out, w, req)
}

func GetFIAS(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"fias_id": Rq.FiasId,
		"request": req.URL.Path, // url
	}

	/*
		ChildrenList, err_c := AddrObj.GetPlaces(dbh, map[string]interface{}{"parentguid": Rq.FiasId})
		if err_c != nil {
			log.Fatal("GetFIAS: ", err_c)
		}
	*/
	One, FullLayers, err_c := AddrObj.GetParentsLyers(dbh, Rq.FiasId)
	if err_c != nil {
		ErrorFound(err_c, w, req)
		return
	}

	if One != nil {
		One["get_id"] = One["aoguid"]
	}

	out := map[string]interface{}{
		"layers": FullLayers,
		"item":   One,
		"meta":   meta,
	}

	JsonSend(out, w, req)
}

func SearchAddress(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"request": req.URL.Path, // url
	}

	One, err_o := Address.GetOneAddress(dbh, Rq.Id)
	if err_o != nil {
		log.Printf("ERROR Address.GetOneAddress: %s\n", err_o)
	}

	List, err_c := Address.GetSearch(dbh, string(Rq.Filters["search_line"]), string(Rq.Filters["parent_list"]))
	if err_c != nil {
		ErrorFound(err_c, w, req)
		return
	}

	out := map[string]interface{}{
		"path":     Rq.SearchFilters["path"],
		"children": List,
		"item":     One,
		"meta":     meta,
		"parents":  Address.GetParents(dbh, iutils.AnyToInt(Rq.Id)),
	}

	JsonSend(out, w, req)
}

func _GetRegionCode(Id int, Deeps ...int) (string, string, error) {

	Deep := 0
	if len(Deeps) > 0 {
		Deep = Deeps[0]
	}

	Deep++

	if Deep > 10 {
		return "", "", errors.New("regioncode not found")
	}

	One, err_c := Address.GetOneAddress(dbh, Id)
	if err_c != nil {
		return "", "", err_c
	}

	fias, find := One["fias_id"]
	fias_id := iutils.AnyToString(fias)
	if !find || "" == fias_id {
		return _GetRegionCode(iutils.AnyToInt(One["parent_id"]), Deep)
	}

	return AddrObj.GetRegionCode(dbh, fias_id)
}

func EditHistoryLoad(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"request": req.URL.Path, // url
	}

	SendError := ""
	One := map[string]interface{}{}
	History := []map[string]interface{}{}
	var err error

	if iutils.AnyToInt(Rq.Id) > 0 {
		One, err = Address.GetOneAddress(dbh, Rq.Id)
		if err != nil {
			SendError += fmt.Sprintf("GetFIASPageView. GetOneAddress - Error: %s\n", err)
		}

		History, err = HandMade.LoadHistory(dbh, "address", iutils.AnyToInt(Rq.Id))
		if err != nil {
			SendError += fmt.Sprintf("GetFIASPageView. LoadHistory - Error: %s\n", err)
		}
	} else {
		History, err = HandMade.LoadLastHistory(dbh, "address")
		if err != nil {
			SendError += fmt.Sprintf("GetFIASPageView. LoadHistory - Error: %s\n", err)
		}
	}

	out := map[string]interface{}{
		"path":    Rq.SearchFilters["path"],
		"item":    One,
		"meta":    meta,
		"history": History,
		"error":   SendError,
		"parents": Address.GetParents(dbh, iutils.AnyToInt(Rq.Id)),
	}

	JsonSend(out, w, req)
}

func GetFIASPageView(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"request": req.URL.Path, // url
	}

	One, err_c := Address.GetOneAddress(dbh, Rq.Id)
	if err_c != nil {
		ErrorFound(err_c, w, req)
		return
	}

	out := map[string]interface{}{
		"path":    Rq.SearchFilters["path"],
		"item":    One,
		"meta":    meta,
		"parents": Address.GetParents(dbh, iutils.AnyToInt(Rq.Id)),
	}

	JsonSend(out, w, req)
}

func GetFIASLike(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	meta := map[string]interface{}{
		"error":   "",
		"result":  1,
		"id":      Rq.Id,
		"request": req.URL.Path, // url
	}

	One, err_c := Address.GetOneAddress(dbh, Rq.Id)
	if err_c != nil {
		ErrorFound(err_c, w, req)
		return
	}

	parent_fias, RegionCode, err_code := _GetRegionCode(iutils.AnyToInt(Rq.Params["parent_id"]))
	if err_code != nil || RegionCode == "" {
		out := map[string]interface{}{
			"children_path": []map[string]interface{}{},
			"item":          One,
			"meta":          meta,
			"path":          Rq.SearchFilters["path"],
		}

		JsonSend(out, w, req)
		return
	}

	FiasFindList, err_ffl := AddrObj.GetLikeIt(dbh, Rq.SearchFilters["path"], RegionCode)
	if err_ffl != nil {
		out := map[string]interface{}{
			"children_path": []map[string]interface{}{},
			"item":          One,
			"meta":          meta,
			"path":          Rq.SearchFilters["path"],
		}

		JsonSend(out, w, req)
		return
	}

	if len(FiasFindList) == 0 {

		out := map[string]interface{}{
			"children_path": []map[string]interface{}{},
			"item":          One,
			"meta":          meta,
			"path":          Rq.SearchFilters["path"],
		}

		JsonSend(out, w, req)
		return
	}

	PathList, err_path := AddrObj.GetMassivePaths(dbh, FiasFindList)
	if err_path != nil {
		ErrorFound(err_path, w, req)
		return
	}

	/* Get list by FIAS ID */
	all_fias_ids := []interface{}{}
	all_fias_str_ids := []string{}
	for i := range PathList {
		fias_id := iutils.AnyToString(PathList[i]["aoguid"])
		if fias_id != "" {
			all_fias_ids = append(all_fias_ids, fias_id)
			all_fias_str_ids = append(all_fias_str_ids, iutils.AnyToString(fias_id))
		}
	}

	house_count := FiasHouses.GetCountByFias(dbh, RegionCode, all_fias_str_ids...)

	if find_by_fias_id, err := Address.GetListByField(dbh, "fias_id", all_fias_ids); err == nil {

		for i := range PathList {
			PathList[i]["id"] = ""
			aoguid := iutils.AnyToString(PathList[i]["aoguid"])
			for _, v := range find_by_fias_id {
				if aoguid == v["fias_id"] {
					PathList[i]["id"] = v["id"]
					break
				}
			}

			if c, find := house_count[aoguid]; find {
				PathList[i]["count_house"] = c
			}
		}
	}

	PathListOut := []map[string]interface{}{}
	if len(parent_fias) > 0 {
		for _, v := range PathList {
			if parent_fias == iutils.AnyToString(iutils.GetKey("aoguid", v)) {
				PathListOut = append(PathListOut, v)
				break
			}

			list := iutils.AnyToStringArray(iutils.GetKey("parent_fias_ids", v))
			for _, f := range list {
				if parent_fias == f {
					PathListOut = append(PathListOut, v)
					break
				}
			}
		}
	} else {
		PathListOut = PathList
	}

	out := map[string]interface{}{
		"path":          Rq.SearchFilters["path"],
		"children_path": PathListOut,
		"item":          One,
		"meta":          meta,
	}

	JsonSend(out, w, req)
}

func GetDiffFIAS(w http.ResponseWriter, req *http.Request) {
	Rq := GetMetaParams(req)

	res, _ := FIAS.GetDiff(db, Rq.Id)

	if res["meta"].(map[string]interface{})["total"] == 0 {
		NotFound(w, req)
		return
	}

	fiasList := []string{}
	for i, v := range res["children"].([]map[string]interface{}) {
		if fias_id := iutils.AnyToString(v["fias_id"]); fias_id != "" {
			res["children"].([]map[string]interface{})[i]["fias_diff"] = 1
			fiasList = append(fiasList, fias_id)
		}
	}
	in := map[string]interface{}{"IN": fiasList}
	params := map[string]interface{}{"fias_id": in}
	addList, err := Address.StGetPlaces(dbh, params)
	if err == nil {
		newList := Merge2array("fias_id", res["children"].([]map[string]interface{}), addList)

		res["children"] = Matching.MergeFromAddress(dbh, newList)
	} else {
		res["children"] = Matching.MergeFromAddress(dbh, res["children"].([]map[string]interface{}))
	}
	JsonSend(res, w, req)
}

func NotFoundBody(serr string) map[string]interface{} {
	return map[string]interface{}{
		"error":  serr,
		"status": 0,
	}
}

func NotFound(w http.ResponseWriter, req *http.Request) {
	JsonSend(NotFoundBody("not found"), w, req)
}

func ErrorFound(err error, w http.ResponseWriter, req *http.Request) {
	log.Println(err)
	JsonSend(NotFoundBody(fmt.Sprintf("%s", err)), w, req)
}

func Status(w http.ResponseWriter, req *http.Request) {
	status := make(map[string]interface{}) //:= Address.Status(dbh)

	status["version"] = Version()

	status["config"] = map[string]interface{}{
		"db": map[string]interface{}{
			"PoolSize": config.DB.PoolSize,
			"DSN":      config.DB.DSN,
		},
		"ids": config.Top_ids,
	}

	status["info"] = map[string]interface{}{
		"address": Address.Status(dbh),
		"fias":    FIAS.Status(db),
	}

	status["status"] = "ok"

	for _, value := range status["info"].(map[string]interface{}) {
		if value.(map[string]interface{})["status"] == "error" {
			status["status"] = "error"
		}
	}

	JsonSend(status, w, req)
}

func GetMetaParams(req *http.Request) ReqParams {

	rq := ReqParams{}

	var err error

	rq.Id, err = strconv.Atoi(string(req.URL.Query().Get(":id")))
	if err != nil {
		rq.Id = 0
	}

	rq.Limit, err = strconv.Atoi(string(req.FormValue("limit")))
	if err != nil {
		rq.Limit = 200
	}
	if rq.Limit == 0 {
		rq.Limit = 200
	}

	rq.Offset, err = strconv.Atoi(string(req.FormValue("offset")))
	if err != nil {
		rq.Offset = 0
	}

	rq.FiasId = string(req.FormValue("fias_id"))
	if len(rq.FiasId) != 36 {
		rq.FiasId = ""
	}

	/* get text-search fielda from request */
	rq.SearchFilters = map[string]string{}
	check_SearchFilters := false
	for name, _ := range Address.SearchFields() {
		if search_line := string(req.FormValue(name)); search_line != "" {
			rq.SearchFilters[string(name)] = search_line
			check_SearchFilters = true
		}
	}

	if check_SearchFilters {
		rq.SearchFilters["searchwhere"] = string(req.FormValue("searchwhere"))
	}

	rq.Filters = map[string]string{}
	for _, v := range Filter.GetFilterNames() {
		val := string(req.FormValue(v))
		if val != "" {
			rq.Filters[string(v)] = val
		}
	}

	rq.Params = map[string]interface{}{}
	ch_params := []string{
		"alt_names", "fias_id", "fias_type_ru", "group_type", "info",
		"is_capital", "lat", "lon", "msg", "name", "official_status_ru", "parent_id",
		"path", "place_type", "point_osm_id", "polygon_osm_id", "type",
	}
	for _, v := range ch_params {
		val := req.FormValue(v)

		if iutils.AnyToString(val) != "" {
			rq.Params[v] = val
		}
	}

	rq.ParamsPoi = map[string]interface{}{}
	pi_params := []string{
		"id", "address_id", "house_id", "point_osm_id", "polygon_osm_id",
		"src", "type", "name", "alt_names", "lat", "lon",
		"path", "info", "place_type", "official_status_ru", "msg", "action_log_id",
	}

	for _, v := range pi_params {
		val := req.FormValue(v)

		if iutils.AnyToString(val) != "" {
			rq.ParamsPoi[v] = val
		}
	}

	return rq
}

//func JsonSuccess( MyAddress *Address.Data, w http.ResponseWriter, req *http.Request ) {
func JsonSuccessPOI(Data *AddressPOI.Data, w http.ResponseWriter, req *http.Request) {
	meta := map[string]interface{}{
		"limit":   Data.Limit,
		"offset":  Data.Offset,
		"total":   Data.Total,
		"what":    Data.What, // countries|regions|places|districts|cites et.c.
		"error":   "",
		"result":  1,
		"request": req.URL.Path, // url
	}

	filters, _ := Filter.GetFilterPoiList(dbh)

	out := map[string]interface{}{
		// Get object methods
		"children": Data.Dumper(),

		// Get package function
		"filters":       filters,
		"search_fields": AddressPOI.SearchFields(),

		"poi": Data.Poi,
	}

	out["meta"] = meta

	JsonSend(out, w, req)
}

//func JsonSuccess( MyAddress *Address.Data, w http.ResponseWriter, req *http.Request ) {
func JsonSuccess(MyAddress *Address.Data, w http.ResponseWriter, req *http.Request) {
	meta := map[string]interface{}{
		"limit":   MyAddress.Limit,
		"offset":  MyAddress.Offset,
		"total":   MyAddress.Total,
		"what":    MyAddress.What, // countries|regions|places|districts|cites et.c.
		"error":   "",
		"result":  1,
		"request": req.URL.Path, // url
	}

	filters, _ := Filter.GetFilterList(dbh)

	out := map[string]interface{}{
		// Get object methods
		"parents":  MyAddress.Parents,
		"children": MyAddress.Dumper(),

		// Get package function
		"filters": filters,
		// Get package function
		"search_fields": Address.SearchFields(),
	}

	out["meta"] = meta

	JsonSend(out, w, req)
}

//func JsonSend( MyAddress *Address.Data, w http.ResponseWriter, req *http.Request ) {
func JsonSend(out map[string]interface{}, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	b, err := json.Marshal(out)
	if err != nil {
		log.Println("error:", err)
	}

	io.WriteString(w, string(b))
}

func ListenStatus(port int) {
	s := pat.New()

	log.Printf("Init *:%v%s\n", port, "/status")
	s.Get("/status", http.HandlerFunc(Status))
	opts := fmt.Sprintf(":%v", port)

	err := http.ListenAndServe(opts, s)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func SendRedirect(urlStr string, w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, urlStr, 302)
}

func init_handlers(m *pat.PatternServeMux, list map[string]interface{}) {
	for path, fun := range list {
		log.Printf("Init %s\n", path)
		m.Get(path, http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
		m.Get(path+"/", http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
	}
}

func main() {

	runtime.GOMAXPROCS(MaxParallelism())

	makeDebug := flag.Int("d", 0, "-d")
	flag.Parse()

	if *makeDebug > 0 {
		IsDebug = true
	}

	config = JSConfig.LoadJsonConfig(config_file)
	port := config.HTTP.PortServer
	go ListenStatus(config.HTTP.PortStatus)

	log.Printf("config: %#v\n", config)

	var err_db error

	log.Printf("Connect DB. PoolSize: %d, DNS: %s\n", config.DB.PoolSize, config.DB.DSN)
	dbh, err_db = dbsearch.DBI(config.DB.PoolSize, config.DB.DSN)
	db = FIAS.InitDb(config.DB.DSN)
	if err_db != nil {
		log.Println("DB error: ", err_db)
	}

	log.Printf("Connect DB. PoolSize: %d, DNS: %s\n", config.DBOSM.PoolSize, config.DBOSM.DSN)
	osm, err_db = dbsearch.DBI(config.DBOSM.PoolSize, config.DBOSM.DSN)
	log.Println("osm: ", osm)

	if err_db != nil {
		log.Println("DB-OSM error: ", err_db)
	}

	// init filters
	if dbh != nil {
		_, filterErr := Filter.GetFilterList(dbh)

		if filterErr != nil {
			log.Printf("somthing wrong with filter: %s", filterErr)
		}

		_, filterErr = Filter.GetFilterPoiList(dbh)

		if filterErr != nil {
			log.Printf("somthing wrong with filter: %s", filterErr)
		}

		dbh.StartReConnect(10)
		osm.StartReConnect(10)
	}

	m := pat.New()

	m.Post("/api/save/item/:id", http.HandlerFunc(SaveAddress))
	m.Post("/api/save/item/:id/", http.HandlerFunc(SaveAddress))

	get_list := map[string]interface{}{
		/*
			urls which contain :id must be more up then the same url without ":id"
		*/
		"/404":                         NotFound,
		"/api/address_view/:id":        GetAddressTree,
		"/api/address_view":            GetAddressTree,
		"/api/countries":               GetAddresses,
		"/api/countries/:id":           GetAddresses,
		"/api/diff/fias/:id":           GetDiffFIAS,
		"/api/fias_merge_load/:id":     GetFIASLike,
		"/api/fias_merge_view/:id":     GetFIASPageView,
		"/api/fias_view":               GetFIAS,
		"/api/get/:id":                 GetAddresses,
		"/api/get":                     GetAddresses,
		"/api/get/fias_merge_view/:id": GetFIASPageView,
		"/api/get/fias_merge_load/:id": GetFIASLike,
		"/api/get_address/:id":         GetAddress,
		"/api/get_address":             GetAddress,
		"/api/house/get/:id":           GetHouses,
		"/api/house/map/:id":           GetHousesMap,
		"/api/one_country_search/:id":  SearchAddress,
		"/api/one_country_search":      SearchAddress,
		"/api/save/new_fias/:id":       SaveNewFIAS,
		"/api/stat/fiasmatching":       CountFiasMatching,
		"/api/history_load/:id":        EditHistoryLoad,
		"/api/poi/:id":                 GetOnePoi,
		"/api/poi_list":                GetPoiList,
	}

	init_handlers(m, get_list)

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	http.Handle("/", m)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Printf("ListenAndServe: ", err)
	} else {
		log.Printf("Listen on %d", port)
	}
}
