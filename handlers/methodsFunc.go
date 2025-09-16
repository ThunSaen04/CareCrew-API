package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/project/carecrew/config"
	"github.com/project/carecrew/models"
	"github.com/project/carecrew/models/assignor"
	v2 "github.com/project/carecrew/models/v2"
	"github.com/project/carecrew/models/worker"
	"github.com/project/carecrew/orther"
)

type Res struct { //สำหรับ Swagger
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var Name = "CareCrew Backend API"
var Versions = "0.2.3"
var Last_Update = "09-17-25"

///////////////////////////////////////////////////////////////////////////////////////////////

func FixFilename(name string) string {
	no := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, f := range no {
		name = strings.ReplaceAll(name, f, "_")
	}
	return name
}

func ApiInfo(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf(
		`
			%s Version %s
			LAST UPDATE: %s
		`, Name, Versions, Last_Update))
}

// เรียกข้อมูลผู้ใช้ทั้งหมด
// GetPersonnelsInfo godoc
// @Summary เรียกข้อมูลผู้ใช้งานทั้งหมด
// @Tags GetMethods
// @Produce json
// @Success 200 {array} models.PersonnelsInfo
// @Failure 400 {object} Res
// @Router /api/personnels [get]
func GetPersonnelsInfo(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	data, err := models.GetPersonnelsInfo(config.DB)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลผู้ใช้งาน")
	return c.JSON(data)
}

// เรียกข้อมูลผู้ใช้จาก ID
// GetPersonnelsInfoWithID godoc
// @Summary เรียกข้อมูลจากหมายเลขผู้ใช้งาน
// @Tags GetMethods
// @Produce json
// @Success 200 {object} models.PersonnelsInfo
// @Failure 400 {object} Res
// @Failure 404 "Not Found"
// @Param personnelID path int true "Personnel ID"
// @Router /api/personnels/{personnelID} [get]
func GetPersonnelsInfoWithID(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	id, err := strconv.Atoi(c.Params("personnelID"))
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลขอเรียกข้อมูลผู้ใช้งานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลขอเรียกข้อมูลผู้ใช้งานไม่ถูกต้อง",
		})
	}

	data, err := models.GetPersonnelsInfo_With_ID(config.DB, id)
	if err != nil {
		log.Print("[Error] เกิดข้อผิดพลาดในการเรียกข้อมูลผู้ใช้งานหมายเลข: ", id, " ", err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	log.Print("[System] พบการเรียกข้อมูลผู้ใช้งานหมายเลข: ", data.PersonnelID)
	return c.JSON(data)
}

// เรียกข้อมูลงานทั้งหมด
// GetTasks godoc
// @Summary เรียกข้อมูลงานทั้งหมด
// @Tags GetMethods
// @Produce json
// @Success 200 {array} models.TasksInfo
// @Failure 400 {object} Res
// @Router /api/tasks [get]
func GetTasks(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	data, err := models.GetTasks(config.DB)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลงาน")
	return c.JSON(data)
}

// เรียกข้อมูลงานจาก ID
// GetTasksWithID godoc
// @Summary เรียกข้อมูลจากหมายเลขงาน
// @Tags GetMethods
// @Produce json
// @Success 200 {object} models.TasksInfo
// @Failure 400 {object} Res
// @Failure 404 Not Found
// @Param taskID path int true "Task ID"
// @Router /api/tasks/{taskID} [get]
func GetTasksWithID(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	id, err := strconv.Atoi(c.Params("taskID"))
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลขอเรียกข้อมูลงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลขอเรียกข้อมูลงานไม่ถูกต้อง",
		})
	}

	data, err := models.GetTasks_With_ID(config.DB, id)
	if err != nil {
		log.Print("[Error] เกิดข้อผิดพลาดในการเรียกข้อมูลงานหมายเลข: ", id, " ", err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	log.Print("[System] พบการเรียกข้อมูลงานหมายเลข: ", data.Task_id)
	return c.JSON(data)
}

// เรียกข้อมูลการแจ้งรายงาน
// GetReport godoc
// @Summary เรียกข้อมูลการแจ้งรายงาน
// @Tags GetMethods
// @Produce json
// @Success 200 {array} models.ReportInfo
// @Failure 400 {object} Res
// @Router /api/greport [get]
func GetReport(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	data, err := models.GetReport(config.DB)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลการแจ้งรายงาน")
	return c.JSON(data)
}

// เรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง
// GetlrubTasksCount godoc
// @Summary เรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง
// @Tags GetMethods
// @Produce json
// @Success 200 {array} models.LrubTasksCountInfo
// @Failure 400 {object} Res
// @Router /api/lrubTasks [get]
func GetlrubTasksCount(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	data, err := models.LrubTasksCount(config.DB)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลจำนวนคนที่รับงาน")
	return c.JSON(data)
}

// ผู้ใช้งานรับงานไหนไปแล้วบ้าง
// GerperlrubTask godoc
// @Summary ผู้ใช้งานรับงานไหนไปแล้วบ้าง
// @Tags GetMethods
// @Produce json
// @Success 200 {array} models.PersonnelLrubTaskInfo
// @Failure 400 {object} Res
// @Router /api/perlrubTasks [get]
func GerperlrubTask(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	data, err := models.PersonnelLrubTask(config.DB)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลรับงาน")
	return c.JSON(data)
}

// ผู้ใช้คนนี้ส่งงานนี้ยัง?
// GetSumbitTaskWithID godoc
// @Summary ผู้ใช้คนนี้ส่งงานนี้ยัง
// @Tags GetMethods
// @Produce json
// @Success 200 {object} models.SubmitTaskWithID
// @Failure 400 {object} Res
// @Param personnelID path int true "Personnel ID"
// @Param taskID path int true "Task ID"
// @Router /api/persubmittasksbor/{personnelID}/{taskID} [get]
func GetSumbitTaskWithID(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	personnel_id, err := strconv.Atoi(c.Params("personnelID"))
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลขอเรียกข้อมูลงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลขอเรียกข้อมูลงานไม่ถูกต้อง",
		})
	}
	task_id, err := strconv.Atoi(c.Params("taskID"))
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลขอเรียกข้อมูลงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลขอเรียกข้อมูลงานไม่ถูกต้อง",
		})
	}

	data, err := models.Get_Submit_Task_With_ID(config.DB, personnel_id, task_id)
	if err != nil {
		log.Print("[Error] เกิดข้อผิดพลาดในการเรียกข้อมูลการส่งงานของผูใช้: ", personnel_id, " (", task_id, ")", " ", err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	log.Print("[System] พบการเรียกข้อมูลการส่งงานของผูใช้: ", personnel_id, " (", task_id, ")")
	return c.JSON(data)
}

// เรียกข้อมูลประเภทงานที่มี
// GetTaskTypeList godoc
// @Summary เรียกข้อมูลประเภทงานที่มี
// @Tags GetMethods
// @Produce json
// @Success 200 {array} models.TasktypeInfo
// @Failure 400 {object} Res
// @Router /api/tasktypelist [get]
func GetTaskTypeList(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	data, err := models.Get_Task_Type_Info(config.DB)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลประเภทงาน")
	return c.JSON(data)
}

// เรียกข้อมูลลำดับความสำคัญ
// GetTaskPriorityList godoc
// @Summary เรียกข้อมูลลำดับความสำคัญ
// @Tags GetMethods
// @Produce json
// @Success 200 {array} models.TaskpriorityInfo
// @Failure 400 {object} Res
// @Router /api/taskprioritylist [get]
func GetTaskPriorityList(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	data, err := models.Get_Task_Priority_Info(config.DB)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลลำดับความสำคัญ")
	return c.JSON(data)
}

// เรียกข้อมูลที่ผู้ใช้ส่งงานมา
// GetTaskEvidence godoc
// @Summary เรียกข้อมูลที่ผู้ใช้ส่งงานมา
// @Tags GetMethods
// @Produce json
// @Success 200 {object} models.TaskEvidenceInfo
// @Failure 400 {object} Res
// @Param taskID path int true "Task ID"
// @Router /api/gettaskevidence/{taskID} [get]
func GetTaskEvidence(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("taskID"))
	data, err := models.GetTaskEvidence(config.DB, id)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}
	log.Print("[System] พบการเรียกข้อมูลส่งงานของผู้ใช้")
	return c.JSON(data)
}

// ล็อคอิน
func Auth(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var authInfo struct {
		PersonnelID int    `json:"personnel_id"`
		Password    string `json:"password"`
	}
	err := c.BodyParser(&authInfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการเข้าสู่ระบบไม่ถูกต้อง")
		// return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการเข้าสู่ระบบไม่ถูกต้อง",
		})
	} else {
		data, err := orther.Auth(config.DB, authInfo.PersonnelID, authInfo.Password)
		if err != nil {
			log.Print("[Error] เกิดข้อผิดพลาดในการเข้าสู่ระบบ")
			// return c.Status(fiber.StatusBadRequest).SendString(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "หมายเลขผู้ใช้งานหรือรหัสผ่านผิด",
			})
		} else {
			log.Print("[System] สามาชิกหมายเลข: ", data.PersonnelID, " เข้าสู่ระบบ")
			return c.JSON(fiber.Map{
				"success":      true,
				"message":      "เข้าสู่ระบบสำเร็จ",
				"personnel_id": data.PersonnelID,
				"role":         data.Role,
			})
		}
	}
}

type AuthInfoV2 struct {
	PersonnelID int    `json:"personnel_id"`
	Password    string `json:"password"`
	Token       string `json:"token"`
}
type Authv2Res struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	PersonnelID int    `json:"personnel_id"`
	Role        string `json:"role"`
}

// ล็อคอิน V2 (ส่ง Token สำหรับใช้ FCM ด้วย)
// Authv2 godoc
// @Summary ล็อคอินเวอร์ชัน 2 (ส่ง Token สำหรับ FCM)
// @Tags Account
// @Accept json
// @Produce json
// @Param request body AuthInfoV2 true "ข้อมูลการเข้าสู่ระบบ"
// @Success 200 {object} Authv2Res "เข้าสู่ระบบสำเร็จ"
// @Failure 400 {object} Res "เข้าสู่ระบบไม่สำเร็จ"
// @Router /api/loginv2 [post]
func Authv2(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var authInfo AuthInfoV2
	err := c.BodyParser(&authInfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการเข้าสู่ระบบไม่ถูกต้อง")
		// return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการเข้าสู่ระบบไม่ถูกต้อง",
		})
	} else {
		data, err := orther.AuthV2(config.DB, authInfo.PersonnelID, authInfo.Password, authInfo.Token)
		if err != nil {
			log.Print(err)
			// return c.Status(fiber.StatusBadRequest).SendString(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		} else {
			log.Print("[System] สามาชิกหมายเลข: ", data.PersonnelID, " เข้าสู่ระบบ")
			return c.JSON(Authv2Res{
				Success:     true,
				Message:     "เข้าสู่ระบบสำเร็จ",
				PersonnelID: data.PersonnelID,
				Role:        data.Role,
			})
		}
	}
}

type RegisterRes struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	PersonnelID int    `json:"personnel_id"`
}

// สมัครบัญชีผู้ใช้
// Register godoc
// @Summary สมัครบัญชีผู้ใช้
// @Tags Account
// @Accept json
// @Produce json
// @Param request body orther.RegisterUserInfo true "ข้อมูลการสมัครบัญชีผู้ใช้"
// @Success 200 {object} RegisterRes "สมัครบัญชีผู้ใช้สำเร็จ"
// @Failure 400 {object} Res "สมัครบัญชีผู้ใช้ไม่สำเร็จ"
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var registerInfo orther.RegisterUserInfo
	var personnelID int

	err := c.BodyParser(&registerInfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการสมัครบัญชีผู้ใช้ไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการสมัครบัญชีผู้ใช้ไม่ถูกต้อง",
		})
	} else {
		personnelID, err = orther.RegisterUser(config.DB, &registerInfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดในการสมัครบัญชีผู้ใช้",
			})
		}
	}
	log.Print("[System] สมัครสามาชิกหมายเลข: ", personnelID)
	return c.JSON(RegisterRes{
		Success:     true,
		Message:     "สมัครสมาชิกสำเร็จ",
		PersonnelID: personnelID,
	})
}

// ลบบัญชี
// Removeacc godoc
// @Summary ลบบัญชี
// @Tags Account
// @Accept json
// @Produce json
// @Param request body AuthInfoV2 true "ข้อมูลการลบบัญชี"
// @Success 200 {object} Res "ลบบัญชีสำเร็จ"
// @Failure 400 {object} Res "ลบบัญชีไม่สำเร็จ"
// @Router /api/removeacc [post]
func Removeacc(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var removeaccinfo orther.RemoveAccInfo

	err := c.BodyParser(&removeaccinfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการลบบัญชีผู้ใช้ไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการลบบัญชีผู้ใช้ไม่ถูกต้อง",
		})
	} else {
		err = orther.RemoveAcc(config.DB, &removeaccinfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}
	log.Print("[System] พบการลบบัญชีผู้ใช้")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบบัญชีผู้ใช้สำเร็จ",
	})
}

type LrubTasksCountvTwo struct {
	TaskId int `json:"task_id"`
}

// เรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง
// GetlrubTasksCountv2 godoc
// @Summary เรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง
// @Tags PostMethods
// @Accept json
// @Produce json
// @Param request body LrubTasksCountvTwo true "ข้อมูลการเรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง"
// @Success 200 {object} models.LrubTasksCountInfo
// @Failure 400 {object} Res
// @Router /api/lrubTasksv2 [post]
func GetlrubTasksCountv2(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var lrubtaskscountvtwo LrubTasksCountvTwo

	err := c.BodyParser(&lrubtaskscountvtwo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการเรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้างไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการเรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้างไม่ถูกต้อง",
		})
	} else {
		data, err := v2.LrubTasksCountV2(config.DB, lrubtaskscountvtwo.TaskId)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}

		log.Print("[System] พบการเรียกข้อมูลว่างานนั้นมีคนรับไปกี่คนและมีใครบ้าง")
		return c.JSON(data)
	}
}

// รายงาน
func Report(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var greportInfo models.ReportInfo

	err := c.BodyParser(&greportInfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการแจ้งรายงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการแจ้งรายงานไม่ถูกต้อง",
		})
	} else {
		err = models.Report(config.DB, &greportInfo)
		if err != nil {
			log.Print("[Error] เกิดข้อผิดพลาดในการแจ้งรายงาน")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดในการแจ้งรายงาน",
			})
		}
	}
	log.Print("พบการแจ้งรายงานจาก: ", greportInfo.PersonnelID)
	return c.JSON(fiber.Map{
		"success":     true,
		"message":     "ส่งรายงานสำเร็จแล้ว",
		"PersonnelID": greportInfo.PersonnelID,
		"Detail":      greportInfo.Detail,
		"Location":    greportInfo.Location,
	})
}

type ReportInfoV2 struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	Title       string   `json:"Title"`
	PersonnelID int      `json:"PersonnelID"`
	Detail      string   `json:"Detail"`
	Location    string   `json:"Location"`
	Files       []string `json:"Files"`
}

// รายงาน v2
// Reportv2 godoc
// @Summary รายงานเวอร์ชั่น2
// @Tags Worker
// @Accept json
// @Produce json
// @Param title formData string true "ชื่อรายงาน"
// @Param personnel_id formData int true "รหัสบุคลากร"
// @Param detail formData string true "รายละเอียดรายงาน"
// @Param location formData string true "สถานที่เกิดเหตุ"
// @Param img formData file true "ไฟล์รูปภาพ"
// @Success 200 {object} ReportInfoV2
// @Failure 400 {object} Res
// @Failure 500 {object} Res
// @Router /api/reportv2 [post]
func Reportv2(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	title := c.FormValue("title")
	personnel_id := c.FormValue("personnel_id")
	detail := c.FormValue("detail")
	location := c.FormValue("location")

	// if !strings.HasSuffix(email, "@rmuti.ac.th") {
	// 	log.Print("[Warning] รูปแบบอีเมล์ไม่ถูกต้อง")
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"success": false,
	// 		"message": "ต้องใช้อีเมล @rmuti.ac.th เท่านั้น",
	// 	})
	// }

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}

	if len(location) == 0 || len(detail) == 0 || len(title) == 0 {
		log.Print("[Warning] กรุณาระบุข้อมูลให้ครบถ้วน")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณาระบุข้อมูลให้ครบถ้วน",
		})
	}

	personnelid, err := strconv.Atoi(personnel_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบของ PersonnelID ไม่ถูกต้อง",
		})
	}

	files := form.File["img"]
	if len(files) == 0 {
		log.Print("[Warning] กรุณาอัพโหลดไฟล์อย่างน้อย 1 ไฟล์")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณาอัพโหลดไฟล์อย่างน้อย 1 ไฟล์",
		})
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	uploadDir := "./imgs/reports"
	os.MkdirAll(uploadDir, os.ModePerm)

	var savedPaths []string
	for index, file := range files {
		repTitle := strings.ReplaceAll(title, " ", "_")
		repLocation := strings.ReplaceAll(location, " ", "_")
		timestmp := time.Now().Format("20060102_150405")
		ext := filepath.Ext(file.Filename)
		filename := FixFilename(fmt.Sprintf("%d_%s_%s_%s%s", index, repTitle, repLocation, timestmp, ext))
		//filename := fmt.Sprintf("%s_%s", repLocation, timestmp)

		fileType := file.Header.Get("Content-Type")
		if !allowedTypes[fileType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("พบไฟล์ที่มีชนิดไม่ถูกต้อง (%s)", fileType),
			})
		}

		savePath := filepath.Join(uploadDir, filename)
		err := c.SaveFile(file, savePath)
		if err != nil {
			log.Print("[Error] ไม่สามารถบันทึกไฟล์ได้")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถบันทึกไฟล์ได้",
			})
		}
		savedPaths = append(savedPaths, "/imgs/reports/"+filename)
	}

	greportInfo := &models.ReportInfo{
		Title:       title,
		PersonnelID: personnelid,
		Detail:      detail,
		Location:    location,
		File:        strings.Join(savedPaths, ","),
	}

	err = models.Report(config.DB, greportInfo)
	if err != nil {

		for _, path := range savedPaths {
			fullpath := config.BasePath + path
			_ = os.Remove(fullpath)
		}

		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการบันทึกข้อมูล",
		})
	}

	log.Print("[System] พบการแจ้งปัญหาใหม่")
	return c.JSON(ReportInfoV2{
		Success:     true,
		Message:     "ส่งรายงานสำเร็จแล้ว",
		Title:       title,
		PersonnelID: personnelid,
		Detail:      detail,
		Location:    location,
		Files:       savedPaths,
	})
}

// เพิ่มงาน
// Addtask godoc
// @Summary เพิ่มงาน
// @Tags Assignor
// @Accept json
// @Produce json
// @Param request body assignor.AddTaskInfo true "ข้อมูลการเพิ่มงาน"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/addtask [post]
func Addtask(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var addtaskinfo assignor.AddTaskInfo

	err := c.BodyParser(&addtaskinfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการเพิ่มงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการเพิ่มงานไม่ถูกต้อง",
		})
	} else {
		err = assignor.AddTask(config.DB, &addtaskinfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}
	log.Print("[System] พบการเพื่มงานใหม่")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "เพิ่มงานใหม่แล้ว",
	})
}

// แก้ไขงาน
// Edittask godoc
// @Summary แก้ไขงาน
// @Tags Assignor
// @Accept json
// @Produce json
// @Param request body assignor.EditTaskInfo true "ข้อมูลการแก้ไขงาน"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/edittask [post]
func Edittask(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var edittaskinfo assignor.EditTaskInfo

	err := c.BodyParser(&edittaskinfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการแก้ไขงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการแก้ไขงานไม่ถูกต้อง",
		})
	} else {
		err = assignor.EditTask(config.DB, &edittaskinfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}

	log.Print("[System] แก้ไขงานหมายเลข: ", edittaskinfo.Task_id, " สำเร็จ")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "แก้ไขงานสำเร็จแล้ว",
	})
}

// ลบงาน
// Removetask godoc
// @Summary ลบงาน
// @Tags Assignor
// @Accept json
// @Produce json
// @Param request body assignor.RemoveTaskInfo true "ข้อมูลการลบงาน"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/removetask [post]
func Removetask(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	var removetaskinfo assignor.RemoveTaskInfo

	err := c.BodyParser(&removetaskinfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการลบงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการลบงานไม่ถูกต้อง",
		})
	} else {
		err = assignor.RemoveTask(config.DB, &removetaskinfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}

	log.Print("[System] ลบงานหมายเลข: ", removetaskinfo.Task_id, " สำเร็จ")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบงานสำเร็จแล้ว",
	})
}

// ลบการแจ้งรายงาน
// Removereport godoc
// @Summary ลบการแจ้งรายงาน
// @Tags Assignor
// @Accept json
// @Produce json
// @Param request body assignor.RemoveReportInfo true "ข้อมูลการลบการแจ้งรายงาน"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/removereport [post]
func Removereport(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	var removereportinfo assignor.RemoveReportInfo

	err := c.BodyParser(&removereportinfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการลบการแจ้งรายงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการลบการแจ้งรายงานไม่ถูกต้อง",
		})
	} else {
		err = assignor.RemoveReport(config.DB, &removereportinfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}

	log.Print("[System] ลบการแจ้งรายงานหมายเลข: ", removereportinfo.Report_id, " สำเร็จ")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบการแจ้งรายงานสำเร็จแล้ว",
	})
}

type EditTaskSt struct {
	TaskID int `json:"task_id"`
}

// การอนุมัติงาน(เสร็จสิ้นงาน)
// TaskSuccess godoc
// @Summary การอนุมัติงาน(เสร็จสิ้นงาน)
// @Tags Assignor
// @Accept json
// @Produce json
// @Param request body EditTaskSt true "ข้อมูลการอนุมัติงาน(เสร็จสิ้นงาน)"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/tasksuccess [post]
func TaskSuccess(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	var tasksuccesinfo EditTaskSt

	err := c.BodyParser(&tasksuccesinfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการอนุมัติงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการอนุมัติงานไม่ถูกต้อง",
		})
	} else {
		err = assignor.TaskSuccess(config.DB, tasksuccesinfo.TaskID)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}
	log.Print("[System] เสร็จสิ้นการตรวจสอบ: ", tasksuccesinfo.TaskID, " สำเร็จ")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "เสร็จสิ้นการตรวจสอบงาน",
	})
}

// รับงาน
// Lrubtask godoc
// @Summary การรับงาน
// @Tags Worker
// @Accept json
// @Produce json
// @Param request body worker.PerLrubTaskInfo true "ข้อมูลการรับงาน"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/lrubtask [post]
func Lrubtask(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var perLrubTaskInfo worker.PerLrubTaskInfo
	err := c.BodyParser(&perLrubTaskInfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการเพิ่มงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการเพิ่มงานไม่ถูกต้อง",
		})
	} else {
		err = worker.PerLrubTask(config.DB, &perLrubTaskInfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "รับงานใหม่แล้ว",
	})
}

// ยกเลิกงาน
// Yoklerk godoc
// @Summary การยกเลิกงาน
// @Tags Worker
// @Accept json
// @Produce json
// @Param request body worker.YokLerkTaskInfo true "ข้อมูลการยกเลิกงาน"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/yoklerktask [post]
func Yoklerk(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var yoklerkTaskInfo worker.YokLerkTaskInfo
	err := c.BodyParser(&yoklerkTaskInfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการยกเลิกไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการยกเลิกไม่ถูกต้อง",
		})
	} else {
		err = worker.YokLerkTask(config.DB, &yoklerkTaskInfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ยกเลิกงานแล้ว",
	})
}

type SongtaskRes struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Files   []string `json:"Files"`
}

// ส่งงาน
// Songtask godoc
// @Summary การส่งงาน
// @Tags Worker
// @Accept json
// @Produce json
// @Param task_id formData string true "ชื่อรายงาน"
// @Param personnel_id formData int true "รหัสบุคลากร"
// @Param img formData file true "ไฟล์รูปภาพ"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/songtask [post]
func Songtask(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	taskid := c.FormValue("task_id")
	personnelid := c.FormValue("personnel_id")
	taskID, err := strconv.Atoi(taskid)
	if err != nil {
		return err
	}
	personnelID, err := strconv.Atoi(personnelid)
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}

	files := form.File["img"]
	if len(files) == 0 {
		log.Print("[Warning] กรุณาแนบรูปหลักฐานอย่างน้อย 1 รูป")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณาแนบรูปหลักฐานอย่างน้อย 1 รูป",
		})
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	uploadDir := "./imgs/attachments"
	os.MkdirAll(uploadDir, os.ModePerm)

	var savedPaths []string
	for index, file := range files {
		timestmp := time.Now().Format("20060102_150405")
		ext := filepath.Ext(file.Filename)
		filename := FixFilename(fmt.Sprintf("%d_tid%d_pid%d_%s%s", index, taskID, personnelID, timestmp, ext))

		fileType := file.Header.Get("Content-Type")
		if !allowedTypes[fileType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("พลไฟล์ที่มีชนิดไม่ถูกต้อง (%s)", fileType),
			})
		}

		savePath := filepath.Join(uploadDir, filename)
		err := c.SaveFile(file, savePath)
		if err != nil {
			log.Print("[Error] ไม่สามารถบันทึกไฟล์ได้")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถบันทึกไฟล์ได้",
			})
		}
		savedPaths = append(savedPaths, "/imgs/attachments/"+filename)
	}

	perSongTaskInfo := worker.PerSongTaskInfo{
		Task_id:      taskID,
		Personnel_id: personnelID,
		File:         strings.Join(savedPaths, ","),
	}

	err = worker.Songtask(config.DB, &perSongTaskInfo)
	if err != nil {

		for _, path := range savedPaths {
			fullpath := config.BasePath + path
			_ = os.Remove(fullpath)
		}
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
		})
	}

	log.Print("[System] ผู้ใช้งานหมายเลข: ", personnelID, " ส่งงานหมายเลข: ", taskID, " สำเร็จ")
	return c.JSON(SongtaskRes{
		Success: true,
		Message: "ส่งงานสำเร็จแล้ว",
		Files:   savedPaths,
	})
}

// ยกเลิกส่งงาน
// YokLerkSongTask godoc
// @Summary การยกเลิกส่งงาน
// @Tags Worker
// @Accept json
// @Produce json
// @Param request body worker.YokLerkSongTaskInfo true "ข้อมูลการยกเลิกส่งงาน"
// @Success 200 {object} Res
// @Failure 400 {object} Res
// @Router /api/yoklerksongtask [post]
func YokLerkSongTask(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")

	var yoklerkSongTaskInfo worker.YokLerkSongTaskInfo
	err := c.BodyParser(&yoklerkSongTaskInfo)
	if err != nil {
		log.Print("[Error] รูปแบบข้อมูลการยกเลิกการส่งงานไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลการยกเลิกการส่งงานไม่ถูกต้อง",
		})
	} else {
		err = worker.Yoklerksongtask(config.DB, &yoklerkSongTaskInfo)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "เกิดข้อผิดพลาดกรุณาติดต่อทีมงานที่เกี่ยวข้องเพื่อแก้ไข",
			})
		}

	}

	log.Print("[System] ผู้ใช้งาน: ", yoklerkSongTaskInfo.Personnel_id, " ยกเลิกการส่งงาน: ", yoklerkSongTaskInfo.Task_id, " สำเร็จ")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "ยกเลิกการส่งงานแล้ว",
	})
}

// ส่งแจ้งเตือน
func SendFCM(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	var fcmInfo struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	err := c.BodyParser(&fcmInfo)
	if err != nil {
		log.Print("[Error] รูปแบบส่งแจ้งเตือนไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบส่งแจ้งเตือนไม่ถูกต้อง",
		})
	}
	err = orther.SendNotificationToAll(config.DB, fcmInfo.Title, fcmInfo.Body)
	if err != err {
		log.Print("[Error] เกิดข้อผิดพลาดในการส่งแจ้งเตือน")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการส่งแจ้งเตือน",
		})
	}

	log.Print("[System] ส่งแจ้งเตือนหาผู้ใช้งานสำเร็จ")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "ส่งแจ้งเตือนสำเร็จแล้ว",
	})
}

// ส่งแจ้งเตือนเฉพาะคนที่รับงาน
func SendFCM2PInT(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	var s orther.SendNotiInfo
	err := c.BodyParser(&s)
	if err != nil {
		log.Print("[Error] รูปแบบส่งแจ้งเตือนไม่ถูกต้อง")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบส่งแจ้งเตือนไม่ถูกต้อง",
		})
	}
	err = orther.SendNotiSuccessToPerInTask(config.DB, &s)
	if err != err {
		log.Print("[Error] เกิดข้อผิดพลาดในการส่งแจ้งเตือน")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการส่งแจ้งเตือน",
		})
	}

	log.Print("[System] ส่งแจ้งเตือนหาผู้ใช้งานสำเร็จ")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "ส่งแจ้งเตือนสำเร็จแล้ว",
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////
