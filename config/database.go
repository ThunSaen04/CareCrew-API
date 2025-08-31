package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB
var Corss = false
var STF = false
var StaticFile = ""
var BasePath = ""
var FullFCMPath = ""
var fcmPath = ""
var APIURL = ""

func ConnectDB() {
	////////////////////////////////////////////////////////////////////
	///////////////////////////// โหลด ENV /////////////////////////////
	err := godotenv.Load()
	if err != nil {
		log.Fatal("โหลดไฟล์ .env ไม่ได้อิอิ")
	}
	log.Print("[System] โหลดไฟล์ .env สำเร็จ")
	////////////////////////////////////////////////////////////////////
	//////////////////////// เปิดใช้งานบลา ๆ ๆ ป่าว ////////////////////////
	Corss, err = strconv.ParseBool(os.Getenv("CORS"))
	if err != nil {
		log.Fatal("แปลงเป็น Bool ไม่สำรเร็จ")
	}

	STF, err = strconv.ParseBool(os.Getenv("STF"))
	if err != nil {
		log.Fatal("แปลงเป็น Bool ไม่สำรเร็จ")
	}
	StaticFile = os.Getenv("STFPATH")
	if len(StaticFile) == 0 {
		log.Fatal("กะรุณาระบบที่อยู่รูปภาพ")
	}
	BasePath = os.Getenv("BASTPATH")
	if len(BasePath) == 0 {
		log.Fatal("ไม่พบตำแหน่งของไฟล์โปรเจค")
	}
	fcmPath = os.Getenv("SERVICE_ACCOUNT")
	if len(fcmPath) == 0 {
		log.Fatal("ไม่พบ SERVICE ACCOUNT FILE")
	}
	FullFCMPath = BasePath + fcmPath
	APIURL = os.Getenv("APIURL")
	if len(APIURL) == 0 {
		log.Fatal("ไม่พบ API URL")
	}

	////////////////////////////////////////////////////////////////////
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("[System] DB OK!")
}
