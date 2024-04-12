package http

import (
	"homework/internal/domain"
	"homework/internal/dto/models"
	"net/http"
	"strconv"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

func setupRouter(r *gin.Engine, uc UseCases) {
	r.HandleMethodNotAllowed = true

	r.POST("/users", SetupCreateUser(uc))
	r.OPTIONS("/users", func(ctx *gin.Context) {
		ctx.Header("Allow", "OPTIONS,POST")
		ctx.AbortWithStatus(http.StatusNoContent)
	})

	r.GET("/sensors", SetupGetSensors(uc))
	r.HEAD("/sensors", SetupHeadGetSensors(uc))
	r.POST("/sensors", SetupCreateSensor(uc))
	r.OPTIONS("/sensors", func(ctx *gin.Context) {
		ctx.Header("Allow", "OPTIONS,POST,GET,HEAD")
		ctx.AbortWithStatus(http.StatusNoContent)
	})

	r.GET("/sensors/:sensorId", SetupGetSensorById(uc))
	r.HEAD("/sensors/:sensorId", SetupHeadGetSensorById(uc))
	r.OPTIONS("/sensors/:sensorId", func(ctx *gin.Context) {
		ctx.Header("Allow", "OPTIONS,GET,HEAD")
		ctx.AbortWithStatus(http.StatusNoContent)
	})

	r.GET("/users/:userId/sensors", SetupGetSensorsByUserId(uc))
	r.HEAD("/users/:userId/sensors", SetupHeadGetSensorsByUserId(uc))
	r.POST("/users/:userId/sensors", SetupAttachSensorByUserId(uc))
	r.OPTIONS("/users/:userId/sensors", func(ctx *gin.Context) {
		ctx.Header("Allow", "OPTIONS,POST,GET,HEAD")
		ctx.AbortWithStatus(http.StatusNoContent)
	})

	r.POST("/events", SetupCreateEvent(uc))
	r.OPTIONS("/events", func(ctx *gin.Context) {
		ctx.Header("Allow", "OPTIONS,POST")
		ctx.AbortWithStatus(http.StatusNoContent)
	})
}

func checkRequestHeader(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodPost {
		if ctx.Request.Header.Get("Content-Type") != "application/json" {
			ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		}
	} else if ctx.Request.Method == http.MethodGet || ctx.Request.Method == http.MethodHead {
		if ctx.Request.Header.Get("Accept") != "application/json" {
			ctx.AbortWithStatus(http.StatusNotAcceptable)
		}
	}
}

func SetupCreateUser(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		userDto := models.UserToCreate{}
		if err := ctx.BindJSON(&userDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := userDto.Validate(strfmt.Default); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		user := domain.User{ID: 1, Name: *userDto.Name}

		if _, err := uc.User.RegisterUser(ctx, &user); err != nil {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.JSON(http.StatusOK, userDto)
	}
}

func getSensors(ctx *gin.Context, uc *UseCases) ([]models.Sensor, int) {
	sensors, err := uc.Sensor.GetSensors(ctx)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	sensorsDto := make([]models.Sensor, len(sensors))
	for i, s := range sensors {
		sensor := s
		strType := string(sensor.Type)
		strReg := strfmt.DateTime(sensor.RegisteredAt)
		strLAct := strfmt.DateTime(sensor.LastActivity)

		sensorsDto[i] = models.Sensor{
			ID:           &sensor.ID,
			SerialNumber: &sensor.SerialNumber,
			Type:         &strType,
			CurrentState: &sensor.CurrentState,
			Description:  &sensor.Description,
			IsActive:     &sensor.IsActive,
			RegisteredAt: &strReg,
			LastActivity: &strLAct,
		}
	}
	return sensorsDto, http.StatusOK
}

func SetupGetSensors(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)
		sensorsDto, hc := getSensors(ctx, &uc)
		ctx.JSON(hc, sensorsDto)
	}
}

func SetupHeadGetSensors(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)
		sensorsDto, hc := getSensors(ctx, &uc)

		sensorsByteSize := int64(0)
		for _, s := range sensorsDto {
			sensorsByteSize += int64(unsafe.Sizeof(s))
		}
		ctx.Header("Content-Length", strconv.FormatInt(sensorsByteSize, 10))
		ctx.Status(hc)
	}
}

func SetupCreateSensor(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		sensorDto := models.SensorToCreate{}
		if err := ctx.BindJSON(&sensorDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := sensorDto.Validate(strfmt.Default); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		tp := domain.SensorType(*sensorDto.Type)
		sensor := domain.Sensor{ID: 1, SerialNumber: *sensorDto.SerialNumber, Type: tp, CurrentState: 0, Description: *sensorDto.Description, IsActive: *sensorDto.IsActive, RegisteredAt: time.Now(), LastActivity: time.Now()}

		if _, err := uc.Sensor.RegisterSensor(ctx, &sensor); err != nil {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.JSON(http.StatusOK, sensor)
	}
}

func getSensorById(ctx *gin.Context, uc *UseCases) (*models.Sensor, int) {
	sensorId, err := strconv.ParseInt(ctx.Param("sensorId"), 10, 64)
	if err != nil {
		return nil, http.StatusUnprocessableEntity
	}
	sensorBindUser := models.SensorToUserBinding{SensorID: &sensorId}
	if err := sensorBindUser.Validate(strfmt.Default); err != nil {
		return nil, http.StatusUnprocessableEntity
	}

	sensor, err := uc.Sensor.GetSensorByID(ctx, *sensorBindUser.SensorID)
	if err != nil {
		return nil, http.StatusNotFound
	}

	strType := string(sensor.Type)
	strReg := strfmt.DateTime(sensor.RegisteredAt)
	strLAct := strfmt.DateTime(sensor.LastActivity)
	sensorDto := models.Sensor{
		ID:           &sensor.ID,
		SerialNumber: &sensor.SerialNumber,
		Type:         &strType,
		CurrentState: &sensor.CurrentState,
		Description:  &sensor.Description,
		IsActive:     &sensor.IsActive,
		RegisteredAt: &strReg,
		LastActivity: &strLAct,
	}
	return &sensorDto, http.StatusOK
}

func SetupGetSensorById(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		sensorDto, hc := getSensorById(ctx, &uc)
		ctx.JSON(hc, sensorDto)
	}
}

func SetupHeadGetSensorById(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		sensorDto, hc := getSensorById(ctx, &uc)
		if sensorDto == nil {
			ctx.Header("Content-Length", "0")
		} else {
			ctx.Header("Content-Length", strconv.FormatInt(int64(unsafe.Sizeof(*sensorDto)), 10))
		}

		ctx.Status(hc)
	}
}

func getSensorsByUserId(ctx *gin.Context, uc *UseCases) ([]models.Sensor, int) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		return nil, http.StatusUnprocessableEntity
	}
	sensorBindUser := models.SensorToUserBinding{SensorID: &userId}
	if err := sensorBindUser.Validate(strfmt.Default); err != nil {
		return nil, http.StatusUnprocessableEntity
	}

	sensors, err := uc.User.GetUserSensors(ctx, *sensorBindUser.SensorID)
	if err != nil {
		return nil, http.StatusNotFound
	}

	sensorsDto := make([]models.Sensor, len(sensors))
	for i, s := range sensors {
		sensor := s
		strType := string(sensor.Type)
		strReg := strfmt.DateTime(sensor.RegisteredAt)
		strLAct := strfmt.DateTime(sensor.LastActivity)

		sensorsDto[i] = models.Sensor{
			ID:           &sensor.ID,
			SerialNumber: &sensor.SerialNumber,
			Type:         &strType,
			CurrentState: &sensor.CurrentState,
			Description:  &sensor.Description,
			IsActive:     &sensor.IsActive,
			RegisteredAt: &strReg,
			LastActivity: &strLAct,
		}
	}

	return sensorsDto, http.StatusOK
}

func SetupGetSensorsByUserId(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		sDto, hc := getSensorsByUserId(ctx, &uc)
		ctx.JSON(hc, sDto)
	}
}

func SetupHeadGetSensorsByUserId(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		sensorsDto, hc := getSensorsByUserId(ctx, &uc)

		sensorsByteSize := int64(0)
		for _, s := range sensorsDto {
			sensorsByteSize += int64(unsafe.Sizeof(s))
		}
		ctx.Header("Content-Length", strconv.FormatInt(sensorsByteSize, 10))
		ctx.Status(hc)
	}
}

func SetupAttachSensorByUserId(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		sensorUserBindDto := models.SensorToUserBinding{}
		if err := ctx.BindJSON(&sensorUserBindDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := sensorUserBindDto.Validate(strfmt.Default); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := uc.User.AttachSensorToUser(ctx, userId, *sensorUserBindDto.SensorID); err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusCreated)
	}
}

func SetupCreateEvent(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		eventDto := models.SensorEvent{}
		if err := ctx.BindJSON(&eventDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := eventDto.Validate(strfmt.Default); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		event := domain.Event{Timestamp: time.Now(), SensorSerialNumber: *eventDto.SensorSerialNumber, SensorID: 1, Payload: *eventDto.Payload}
		if err := uc.Event.ReceiveEvent(ctx, &event); err != nil {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.Status(http.StatusCreated)
	}
}
