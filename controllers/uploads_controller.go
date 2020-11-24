package controllers

import (
	"net/http"

	"github.com/jpr98/apis_pf_back/datastore"
	"github.com/labstack/echo/v4"
)

// Uploads represents an uploads controller
type Uploads struct {
	uploadsStore datastore.StorageDatastore
}

// NewUploadsController creates a new uploads controlelr with a store
func NewUploadsController(us datastore.StorageDatastore) Uploads {
	return Uploads{uploadsStore: us}
}

// Upload uploads a file
func (u *Uploads) Upload(c echo.Context) error {
	name := c.FormValue("name")
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	err = u.uploadsStore.Upload(name, src)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Couldn't upload file")
	}

	url := u.uploadsStore.URL + "/" + name
	return c.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}
