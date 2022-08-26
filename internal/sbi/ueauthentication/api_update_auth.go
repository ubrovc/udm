/*
 * NudmUEAU
 *
 * UDM UE Authentication Service
 *
 */

package ueauthentication

import (
	"net/http"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udm/internal/logger"
	"github.com/free5gc/udm/internal/sbi/producer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
)

// HTTPDeleteAuth ... - Delete an authentication result
func HTTPDeleteAuth(c *gin.Context) {
	var authEvent models.AuthEvent
	// step 1: retrieve http request body
	requestBody, err := c.GetRawData()
	if err != nil {
		// pdErr := models.NewProblemDetailsEx(models.ProblemDetailsCause_SYSTEM_FAILURE, "", err.Error())
		pdErr := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UeauLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, pdErr)
		return
	}

	// step 2: convert requestBody to openapi models
	err = openapi.Deserialize(&authEvent, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		// pdErr := models.NewProblemDetailsEx(models.ProblemDetailsCause_UNSPECIFIED_MSG_FAILURE, "Malformed request syntax", problemDetail)
		pdErr := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
			Cause:  "Malformed request syntax",
		}
		logger.UeauLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, pdErr)
		return
	}

	// TODO: checking if input data are valid. Add valid function check in model and call it:
	// if pdErr := authEvent.Valid(); pdErr != nil {
	// 	logger.UeauLog.Errorln(pdErr)
	// 	c.JSON(int(pdErr.Status), pdErr)
	// 	return
	// }

	req := httpwrapper.NewRequest(c.Request, authEvent)

	// params validation
	// supi
	// TODO: Check supi format. Make supi model and add valid function check.
	// use model: supi := models.Supi(c.Params.ByName("supi"))
	// valid check: pdErr := supi.Valid()
	supi := c.Params.ByName("supi")
	if supi == "" {
		pdErr := models.ProblemDetails{
			Title:  "Mandatory IE missing",
			Status: http.StatusBadRequest,
			Detail: "Missing mandatory param supi",
			Cause:  "MANDATORY_IE_MISSING",
		}
		logger.UeauLog.Errorln(pdErr.Detail)
		c.JSON(int(pdErr.Status), pdErr)
		return
	}
	// TODO: Check supi format. Make supi model and add valid function check. Use valid function at this point:
	// else if pdErr := supi.Valid(); pdErr != nil {
	// 	logger.UeauLog.Errorln(pdErr.Detail)
	// 	c.JSON(int(pdErr.Status), pdErr)
	// 	return
	// }
	req.Params["supi"] = string(supi)

	// authEventId
	req.Params["authEventId"] = c.Params.ByName("authEventId")
	if req.Params["authEventId"] == "" {
		pdErr := models.ProblemDetails{
			Title:  "Mandatory IE missing",
			Status: http.StatusBadRequest,
			Detail: "Missing mandatory param authEventId",
			Cause:  "MANDATORY_IE_MISSING",
		}
		logger.UeauLog.Errorln(pdErr.Detail)
		c.JSON(int(pdErr.Status), pdErr)
		return
	}

	rsp := producer.HandleDeleteAuthDataRequest(req)

	for key, val := range rsp.Header {
		c.Header(key, val[0])
	}
	if err != nil {
		logger.UeauLog.Errorln(err)
		// pdErr := models.NewProblemDetailsEx(models.ProblemDetailsCause_SYSTEM_FAILURE, "", err.Error())
		pdErr := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(int(pdErr.Status), pdErr)
	} else {
		c.Status(rsp.Status)
	}
}
