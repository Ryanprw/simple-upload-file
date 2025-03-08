package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ryanprw/simple-upload-file/core/utils"
	"github.com/gofiber/fiber/v2"
)

func Route(router fiber.Router) {
	router.Post("/", uploadFiles)

	router.Get("/", listUploadedFiles)
}

func listUploadedFiles(c *fiber.Ctx) error {
	uploadDir := "./public/uploads"
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		fmt.Println("Error reading upload directory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.BaseResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to read upload directory",
		})
	}

	var fileUrls []string
	for _, file := range files {
		fileUrls = append(fileUrls, "uploads/"+file.Name())
	}

	return c.Status(fiber.StatusOK).JSON(utils.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Uploaded files",
		Data:    fileUrls,
	})
}

func uploadFiles(c *fiber.Ctx) error {
	fmt.Println("Receiving file upload request...")

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Failed to parse uploaded files",
		})
	}

	files := form.File["file"]
	if len(files) == 0 {
		fmt.Println("No files received")
		return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
			Code:    fiber.StatusBadRequest,
			Message: "No files uploaded",
		})
	}

	uploadDir := "./public/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating upload directory:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(utils.BaseResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to create upload directory",
			})
		}
	}

	var fileUrls []string
	for _, file := range files {
		filePath := filepath.Join(uploadDir, file.Filename)
		fileUrl := "uploads/" + file.Filename

		err := c.SaveFile(file, filePath)
		if err != nil {
			fmt.Println("Error saving file:", file.Filename, err)
			return c.Status(fiber.StatusInternalServerError).JSON(utils.BaseResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to save file",
			})
		}

		fmt.Println("File saved:", file.Filename, "➡", fileUrl)
		fileUrls = append(fileUrls, fileUrl)
	}

	fmt.Println("All files uploaded successfully!")

	return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
		Code:    fiber.StatusCreated,
		Message: "Files uploaded successfully",
		Data:    fileUrls,
	})
}
