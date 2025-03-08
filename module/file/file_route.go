package file

import (
	"strconv"
	"time"

	"github.com/Ryanprw/simple-upload-file/core/utils"
	"github.com/gofiber/fiber/v2"
)

func Route(router fiber.Router) {
	router.Post("/", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Failed to parse uploaded files",
			})
		}

		files := form.File["file"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Code:    fiber.StatusBadRequest,
				Message: "There's no file provided",
			})
		}

		file := files[0]
		milis := strconv.FormatInt(time.Now().UnixMilli(), 10)
		filePath := "./public/uploads/" + milis + "_" + file.Filename
		fileUrl := "uploads/" + milis + "_" + file.Filename
		err = c.SaveFile(file, filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.BaseResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Failed to save uploaded file",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Code:    fiber.StatusCreated,
			Message: "Saved",
			Data:    fileUrl,
		})
	})

	router.Post("/bulk", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Failed to parse uploaded files",
			})
		}

		var compressedFileURLs []string
		files := form.File["file[]"]
		for _, file := range files {
			milis := strconv.FormatInt(time.Now().UnixMilli(), 10)
			filePath := "./public/uploads/" + milis + "_" + file.Filename
			fileUrl := "uploads/" + milis + "_" + file.Filename
			err = c.SaveFile(file, filePath)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(utils.BaseResponse{
					Code:    fiber.StatusBadRequest,
					Message: "Failed to save uploaded file",
				})
			}
			compressedFileURLs = append(compressedFileURLs, fileUrl)
		}
		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Code:    fiber.StatusCreated,
			Message: "Saved",
			Data:    compressedFileURLs,
		})

	})
}
