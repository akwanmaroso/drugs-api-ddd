package interfaces

import (
	"fmt"
	"github.com/akwanmaroso/ddd-drugs/application"
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Drug struct {
	drugApp          application.DrugAppInterface
	userApp          application.UserAppInterface
	tokenInterface   auth.TokenInterface
	refreshInterface auth.AuthInterface
}

func NewDrug(drugApp application.DrugAppInterface, userApp application.UserAppInterface, tokenInterface auth.TokenInterface, refreshInterface auth.AuthInterface) *Drug {
	return &Drug{drugApp: drugApp, userApp: userApp, tokenInterface: tokenInterface, refreshInterface: refreshInterface}
}

func (drug *Drug) SaveDrug(c *gin.Context) {
	// check if this user is authenticated
	metadata, err := drug.tokenInterface.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// lookup the metadata in redis
	userId, err := drug.refreshInterface.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	var saveDrugError = make(map[string]string)

	name := c.PostForm("name")
	description := c.PostForm("description")
	image := c.PostForm("image")
	if fmt.Sprintf("%T", name) != "string" || fmt.Sprintf("%T", description) != "string" || fmt.Sprintf("%T", image) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}

	emptyDrug := entity.Drug{}
	emptyDrug.Name = name
	emptyDrug.Description = description
	emptyDrug.DrugImage = image
	saveDrugError = emptyDrug.Validate("")

	if len(saveDrugError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveDrugError)
		return
	}

	// check if the user exists
	_, err = drug.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	var newDrug = entity.Drug{}
	newDrug.UserID = userId
	newDrug.Name = name
	newDrug.Description = description
	newDrug.DrugImage = image
	savedDrug, saveErr := drug.drugApp.SaveDrug(&newDrug)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusOK, savedDrug)
}

func (drug *Drug) UpdateDrug(c *gin.Context) {
	// check if this user is authenticated
	metadata, err := drug.tokenInterface.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// lookup the metadata in redis
	userId, err := drug.refreshInterface.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	var updateDrugError = make(map[string]string)

	drugId, err := strconv.ParseUint(c.Param("drug_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	image := c.PostForm("image")
	if fmt.Sprintf("%T", name) != "string" || fmt.Sprintf("%T", description) != "string" || fmt.Sprintf("%T", image) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
		return
	}

	emptyDrug := entity.Drug{}
	emptyDrug.Name = name
	emptyDrug.Description = description
	emptyDrug.DrugImage = image
	updateDrugError = emptyDrug.Validate("update")
	if len(updateDrugError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateDrugError)
		return
	}
	user, err := drug.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	// check if the drug exist:
	drugToUpdate, err := drug.drugApp.GetDrug(drugId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	if user.ID != drugToUpdate.UserID {
		c.JSON(http.StatusUnauthorized, "your not the owner of this drug")
		return
	}

	drugToUpdate.Name = name
	drugToUpdate.Description = description
	drugToUpdate.UpdatedAt = time.Now()
	updatedDrug, dbUpdateErr := drug.drugApp.UpdateDrug(drugToUpdate)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}

	c.JSON(http.StatusOK, updatedDrug)
}

func (drug *Drug) GetAllDrug(c *gin.Context) {
	drugs, err := drug.drugApp.GetAllDrug()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, drugs)
}

func (drug *Drug) GetDrugAndCreator(c *gin.Context) {
	drugId, err := strconv.ParseUint(c.Param("drug_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	drugs, err := drug.drugApp.GetDrug(drugId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	user, err := drug.userApp.GetUser(drugs.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	drugAndUser := map[string]interface{}{
		"drug":    drugs,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, drugAndUser)
}

func (drug *Drug) DeleteDrug(c *gin.Context) {
	metadata, err := drug.tokenInterface.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	drugId, err := strconv.ParseUint(c.Param("drug_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	_, err = drug.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = drug.drugApp.DeleteDrug(drugId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Drug Delete")
}
