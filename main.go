package main

import (
	"fmt"
	"path/filepath"
	"net/http"
	"os"
	"io"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"golang.org/x/net/context"
	"github.com/akkuman/parseConfig"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

var  downLoadDirPath , accessKey , secretKey , upToken , key , videoDomain string
var  cfg storage.Config
var  config parseConfig.Config

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func main () {
	config = parseConfig.New("./config.json")
	//七牛云配置
	downLoadDirPath = "./"
	accessKey =  config.Get("qiniu_config > accessKey").(string)
	secretKey =  config.Get("qiniu_config > secretKey").(string)
	bucket :=  config.Get("qiniu_config > bucket").(string)
	videoDomain =  config.Get("qiniu_config > url").(string)

	mac := qbox.NewMac(accessKey , secretKey)
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken = putPolicy.UploadToken(mac)
	fmt.Println(upToken)
	cfg := storage.Config{}
	//华北机房
	cfg.Zone = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	//输入参数
	inputPars()
}

func inputPars () {
	fmt.Println("请输入视频URL:");
	var fileName , fileUrl string
	fmt.Scanln(&fileUrl)
	if len(fileUrl) < 1 {
		fmt.Println("参数不能为空")
		return
	}
	fmt.Printf("请输入保存文件名:");
	fmt.Scanln(&fileName)
	if len(fileName) < 1 {
		fmt.Println("参数不能为空")
		return
	}
	key = fileName

	go downLoad(fileName , fileUrl)
	fmt.Printf("视频文件：《%s》已放入后台下载\n" , fileName)
	//再次调用
	inputPars()
}

/**
 *@fileName 文件名 string
 *@fileUrl 要下载的文件url string
 */
func downLoad (fileName , fileUrl string) {
	fileName = filepath.Join(downLoadDirPath , fileName)
	fileRes , err := http.Get(fileUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	file , err := os.Create(fileName + ".mp4")
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(file , fileRes.Body)
	formUploader := storage.NewFormUploader(&cfg)
	ret := MyPutRet{}

	err = formUploader.PutFile(context.Background() , &ret , upToken , key , file.Name() , nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	//mysql配置
	host :=  config.Get("mysql_config > host").(string)
	user :=  config.Get("mysql_config > user").(string)
	passwd :=  config.Get("mysql_config > passwd").(string)
	port :=  config.Get("mysql_config > port").(string)
	database :=  config.Get("mysql_config > database").(string)

	db , _ := sql.Open("mysql" , user + ":" + passwd + "@tcp(" + host + ":" + port +")/" + database + "?charset=utf8")
	defer db.Close()
	videoUrl := videoDomain + key
	//destStr , err := iconv.Open("utf-8","gb2312")
	//if err != nil{
	//	fmt.Println("iconv.Openfailed!")
	//}
	//
	//defer destStr.Close()
	//插入数据
	db.Exec("insert into video(name,url) values('"+ key +"', '"+ videoUrl +"')");
	//ins_id, _ := res.LastInsertId();
	//fmt.Println(ins_id);
	//fmt.Println(ret)
}

