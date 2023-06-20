package ctrl

import (
	"fine-grained-openai-proxy/svc"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

/*
AllFineGrainedKeys Get all fine-grained keys

Path:

	/admin/fgkey/all

Args:

	GET auth: Admin token
*/
func AllFineGrainedKeys(c *fiber.Ctx) error {
	fgSvc := svc.FineGrainedKeySvc{}
	fgs, err := fgSvc.All()
	if err != nil {
		return c.Status(fiber.StatusTeapot).JSON(Resp{Code: 1, Error: "Get All Fine Grained Key Error: " + err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0,
		Data: fgs,
	})
}

/*
FineGrainedKeysByParentID Get fine-grained keys by parent ID

Path:

	/admin/fgkey/parentid

Args:

	GET auth: Admin token
	POST parent_id: Parent OpenAI API Key ID
*/
func FineGrainedKeysByParentID(c *fiber.Ctx) error {
	pids := c.FormValue("parent_id")
	pid, err := strconv.ParseInt(pids, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Parent ID Error: " + err.Error()})
	}

	fgSvc := svc.FineGrainedKeySvc{}
	fgs, err := fgSvc.ByParentID(pid)
	if err != nil {
		return c.Status(fiber.StatusTeapot).JSON(Resp{Code: 1, Error: "Get Fine Grained Key By Parent ID Error: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0,
		Data: fgs,
	})
}

/*
InsertFineGrainedKey Insert new fine-grained key

Path:

	/admin/fgkey/insert

Args:

	GET auth: Admin token
	POST parent_id: Parent OpenAI API Key ID
	POST type: whitelist or blacklist
	POST list: OpenAI Model ID JSON int arr (e.g. [1, 2, 3])
	POST expire: Expire time (e.g. 2023-12-31)
	POST remain_calls: Remain calls (int)
*/
func InsertFineGrainedKey(c *fiber.Ctx) error {
	parentIDs := c.FormValue("parent_id")
	types := c.FormValue("type")
	listJson := c.FormValue("list")
	expires := c.FormValue("expire")
	remainCallss := c.FormValue("remain_calls")

	parentID, err := strconv.ParseInt(parentIDs, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Parent ID Error: " + err.Error()})
	}

	if types != "whitelist" && types != "blacklist" {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Type Error: " + err.Error()})
	}

	// xxx?expire=2023-01-01 --> time.Time, if want unix timestamp, use time.Unix()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	expire, err := time.ParseInLocation("2006-01-02", expires, loc)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Expire is invalid, Error at: " + err.Error()})
	}

	remainCalls, err := strconv.ParseInt(remainCallss, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Remain Calls Error: " + err.Error()})
	}

	fgSvc := svc.FineGrainedKeySvc{}
	key, err := fgSvc.Insert(&svc.FineGrainedKey{
		ParentID:    parentID,
		Type:        types,
		List:        listJson,
		Expire:      expire.Unix(),
		RemainCalls: remainCalls,
	})
	if err != nil {
		return c.Status(fiber.StatusTeapot).JSON(Resp{Code: 1, Error: "Insert Fine Grained Key Error: " + err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0,
		Msg: "Insert Fine Grained Key Success. \n" +
			"Notice: For security reasons, you won't be able to view it again. \n" +
			"If you lose this secret key, you'll need to generate a new one.",
		Data: key,
	})
}

/*
UpdateFineGrainedKey Update fine-grained key

Path:

	/admin/fgkey/update

Args:

	GET auth: Admin token
	POST id: Fine-grained key ID
	POST parent_id: Parent OpenAI API Key ID
	POST type: whitelist or blacklist
	POST list: OpenAI Model ID JSON int arr (e.g. [1, 2, 3])
	POST expire: Expire time (e.g. 2023-12-31)
	POST remain_calls: Remain calls (int)
*/
func UpdateFineGrainedKey(c *fiber.Ctx) error {
	ids := c.FormValue("id")
	parentIDs := c.FormValue("parent_id")
	types := c.FormValue("type")
	listJson := c.FormValue("list")
	expires := c.FormValue("expire")
	remainCallss := c.FormValue("remain_calls")

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "ID Error: " + err.Error()})
	}

	parentID, err := strconv.ParseInt(parentIDs, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Parent ID Error: " + err.Error()})
	}

	if types != "whitelist" && types != "blacklist" {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Type Error: " + err.Error()})
	}

	// xxx?expire=2023-01-01 --> time.Time, if want unix timestamp, use time.Unix()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	expire, err := time.ParseInLocation("2006-01-02", expires, loc)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Expire is invalid, Error at: " + err.Error()})
	}

	remainCalls, err := strconv.ParseInt(remainCallss, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Remain Calls Error: " + err.Error()})
	}

	fgSvc := svc.FineGrainedKeySvc{}
	err = fgSvc.Update(&svc.FineGrainedKey{
		ID:          id,
		ParentID:    parentID,
		Type:        types,
		List:        listJson,
		Expire:      expire.Unix(),
		RemainCalls: remainCalls,
	})
	if err != nil {
		return c.Status(fiber.StatusTeapot).JSON(Resp{Code: 1, Error: "Update Fine Grained Key Error: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0,
		Msg:  "Update Fine Grained Key Success",
	})
}

/*
DeleteFineGrainedKey Delete fine-grained key

Path:

	/admin/fgkey/delete

Args:

	GET auth: Admin token
	POST id: Fine-grained key ID
*/
func DeleteFineGrainedKey(c *fiber.Ctx) error {
	ids := c.FormValue("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "ID Error: " + err.Error()})
	}

	fgSvc := svc.FineGrainedKeySvc{}
	err = fgSvc.Delete(&svc.FineGrainedKey{ID: id})
	if err != nil {
		return c.Status(fiber.StatusTeapot).JSON(Resp{Code: 1, Error: "Delete Fine Grained Key Error: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0,
		Msg:  "Delete Fine Grained Key Success",
	})
}
