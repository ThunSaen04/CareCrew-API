package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/project/carecrew/config"
	_ "github.com/project/carecrew/docs"
	"github.com/project/carecrew/handlers"
)

// @title Backend API CareCrew
// @description API สำหรับแอปพลิเคชันจัดสรรงานบุคลากรแม่บ้านและภารโรง
// @contact.name Thun Saen
// @contact.email biz.thunsaen@gmail.com
// @host api.lcadv.online
// @BasePath /
func main() {

	limitfilesize := 50 //หน่อยเป็น M

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		AppName:   fmt.Sprintf("%s v%s", handlers.Name, handlers.Versions),
		BodyLimit: limitfilesize * 1024 * 1024,
	})

	if config.Corss { //เปิดใช้เฉพาะตอนทดสอบ API ด้วย web tools เท่านั้น
		app.Use(cors.New())
		log.Print("[System] เปิดการใช้งาน CORS")
	} else {
		log.Print("[System] ปิดการใช้งาน CORS")
	}
	if config.STF {
		app.Static("/imgs", fmt.Sprintf(`%s`, config.StaticFile))
		log.Print("[System] เปิดการใช้งาน STATIC FILE")
	} else {
		log.Print("[System] ปิดการใช้งาน STATIC FILE")
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	//Get Methods
	app.Get("/api/personnels", handlers.GetPersonnelsInfo)                               // เรียกข้อมูลผู้ใช้ทั้งหมด
	app.Get("/api/personnels/:personnelID", handlers.GetPersonnelsInfoWithID)            // เรียกข้อมูลผู้ใช้จาก ID
	app.Get("/api/tasks", handlers.GetTasks)                                             // เรียกข้อมูลงานทั้งหมด
	app.Get("/api/tasks/:taskID", handlers.GetTasksWithID)                               // เรียกข้อมูลงานจาก ID
	app.Get("/api/greport", handlers.GetReport)                                          // เรียกข้อมูลการรีพอร์ต
	app.Get("/api/lrubTasks", handlers.GetlrubTasksCount)                                // เรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง
	app.Get("/api/perlrubTasks", handlers.GerperlrubTask)                                // ผู้ใช้งานรับงานไหนไปแล้วบ้าง
	app.Get("/api/persubmittasksbor/:personnelID/:taskID", handlers.GetSumbitTaskWithID) // ผู้ใช้คนนี้ส่งงานนี้ยัง?
	app.Get("/api/tasktypelist", handlers.GetTaskTypeList)                               // เรียกข้อมูลประเภทงานที่มี
	app.Get("/api/taskprioritylist", handlers.GetTaskPriorityList)                       // เรียกข้อมูลลำดับความสำคัญ
	app.Get("/api/notis", handlers.GetNotis)                                             // เรียกข้อมูลแจ้งเตือน

	//Post Methods
	//app.Post("/api/login", handlers.Auth)          // ล็อคอิน
	app.Post("/api/loginv2", handlers.Authv2)      // ล็อคอินv2 ส่ง FCM Device Token มาด้วยด้วย..
	app.Post("/api/register", handlers.Register)   // สมัคร ปล.ใช้ได้เฉพาะผู้ดูแลระบบ ที่มีสิทธิในการสร้างผู้ใช้งาน..
	app.Post("/api/removeacc", handlers.Removeacc) // ลบบัญชีผู้ใช้งาน..

	app.Post("/api/lrubTasksv2", handlers.GetlrubTasksCountv2) // เรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง แบบ Post ระบบ task_id

	//app.Post("/api/report", handlers.Report)     // รีพอร์ต
	app.Post("/api/reportv2", handlers.Reportv2) // รีพอร์ต2 รับหลายรูป..

	//ผู้มอบหมายงาน
	app.Post("/api/addtask", handlers.Addtask)                        // เพิ่มงาน..
	app.Post("/api/removetask", handlers.Removetask)                  // ลบงาน..
	app.Post("/api/removereport", handlers.Removereport)              // ลบการแจ้งรายงาน..
	app.Post("/api/edittask", handlers.Edittask)                      // แก้ไขงาน..
	app.Get("/api/gettaskevidence/:taskID", handlers.GetTaskEvidence) // เรียกข้อมูลที่ผู้ใช้ส่งงานมา
	//app.Post("/api/tasksuccescheck", nil)                             // โหลดข้อมูลที่ผู้ใช้ส่งเข้ามา สำหรับตรวจสอบงาน
	app.Post("/api/tasksuccess", handlers.TaskSuccess) // การอนุมัติงาน(เสร็จสิ้นงาน)..
	app.Post("/api/nosuccess", handlers.NoSuccess)     // การอนุมัติงาน(ปฏิเสธ)
	app.Post("/api/readnotis", handlers.ReadNotis)     // อัพเดทว่าดูแจ้งเตือนรึยัง
	//app.Post("/api/loaddetailfromreport", nil)                        // โหลดข้อมูลงานจาก รีพอร์ต

	//ผู้รับงาน
	app.Post("/api/lrubtask", handlers.Lrubtask)               // รับงาน..
	app.Post("/api/yoklerktask", handlers.Yoklerk)             // ยกเลิกงาน..
	app.Post("/api/songtask", handlers.Songtask)               // ส่งงาน..
	app.Post("/api/yoklerksongtask", handlers.YokLerkSongTask) // ยกเลิกส่งงาน..

	app.Put("/api/personnels/:id/profile", handlers.UpdateProfile) // เปลี่ยนโปรไฟล์

	//FCM
	//app.Post("/api/fcmsend", handlers.SendFCM)
	//app.Post("/api/fcmsendtoperintask", handlers.SendFCM2PInT) //ใช้ทดสอบส่งแจ้งเตือนไปหาผู้ใช้ที่รับงานนั้น

	app.Get("/", handlers.ApiInfo)

	app.Listen(":3000")
}
