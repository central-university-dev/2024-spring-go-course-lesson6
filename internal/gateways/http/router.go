package http

import (
	"homework/internal/domain"
	"homework/internal/dto"
	"homework/internal/validate"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/gin-gonic/gin"
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
	r.POST("/users/:userId/sensors", SetupCreateSensorByUserId(uc))
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

		userDto := dto.User{}
		if err := ctx.BindJSON(&userDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		userDto.InitData()
		if err := validate.Validate(userDto); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		user := domain.User{ID: userDto.ID, Name: userDto.Name}

		if _, err := uc.User.RegisterUser(ctx, &user); err != nil {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.JSON(http.StatusOK, userDto)
	}
}

func getSensors(ctx *gin.Context, uc *UseCases) ([]dto.Sensor, int) {
	sensors, err := uc.Sensor.GetSensors(ctx)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	sensorsDto := make([]dto.Sensor, len(sensors))
	for i, s := range sensors {
		sensor := s
		sensorsDto[i] = dto.Sensor{
			ID:           &sensor.ID,
			SerialNumber: &sensor.SerialNumber,
			Type:         sensor.Type,
			CurrentState: sensor.CurrentState,
			Description:  sensor.Description,
			IsActive:     sensor.IsActive,
			RegisteredAt: sensor.RegisteredAt,
			LastActivity: sensor.LastActivity,
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

		sensorDto := dto.Sensor{}
		if err := ctx.BindJSON(&sensorDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		sensorDto.InitData()
		if err := validate.Validate(sensorDto); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		sensor := domain.Sensor{ID: *sensorDto.ID, SerialNumber: *sensorDto.SerialNumber, Type: sensorDto.Type, CurrentState: sensorDto.CurrentState, Description: sensorDto.Description, IsActive: sensorDto.IsActive, RegisteredAt: sensorDto.RegisteredAt, LastActivity: sensorDto.LastActivity}
		if _, err := uc.Sensor.RegisterSensor(ctx, &sensor); err != nil {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.JSON(http.StatusOK, sensor)
	}
}

func getSensorById(ctx *gin.Context, uc *UseCases) (*dto.Sensor, int) {
	sensorId, err := strconv.ParseInt(ctx.Param("sensorId"), 10, 64)
	if err != nil {
		return nil, http.StatusUnprocessableEntity
	}

	sensor, err := uc.Sensor.GetSensorByID(ctx, sensorId)
	if err != nil {
		return nil, http.StatusNotFound
	}

	sensorDto := dto.Sensor{ID: &sensor.ID, SerialNumber: &sensor.SerialNumber, Type: sensor.Type, CurrentState: sensor.CurrentState, Description: sensor.Description, IsActive: sensor.IsActive, RegisteredAt: sensor.RegisteredAt, LastActivity: sensor.LastActivity}
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

func getSensorsByUserId(ctx *gin.Context, uc *UseCases) ([]dto.Sensor, int) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		return nil, http.StatusUnprocessableEntity
	}

	sensors, err := uc.User.GetUserSensors(ctx, userId)
	if err != nil {
		return nil, http.StatusNotFound
	}

	sensorsDto := make([]dto.Sensor, len(sensors))
	for i, s := range sensors {
		sensor := s
		sensorsDto[i] = dto.Sensor{
			ID:           &sensor.ID,
			SerialNumber: &sensor.SerialNumber,
			Type:         sensor.Type,
			CurrentState: sensor.CurrentState,
			Description:  sensor.Description,
			IsActive:     sensor.IsActive,
			RegisteredAt: sensor.RegisteredAt,
			LastActivity: sensor.LastActivity,
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

func SetupCreateSensorByUserId(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		sensorDto := dto.Sensor{}
		if err := ctx.BindJSON(&sensorDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		sensorDto.InitData()
		if err := validate.Validate(sensorDto); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		sensor := domain.Sensor{ID: *sensorDto.ID, SerialNumber: *sensorDto.SerialNumber, Type: sensorDto.Type, CurrentState: sensorDto.CurrentState, Description: sensorDto.Description, IsActive: sensorDto.IsActive, RegisteredAt: sensorDto.RegisteredAt, LastActivity: sensorDto.LastActivity}

		if err := uc.User.CreateSensorToUser(ctx, userId, &sensor); err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusCreated)
	}
}

func SetupCreateEvent(uc UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkRequestHeader(ctx)

		eventDto := dto.Event{}
		if err := ctx.BindJSON(&eventDto); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		eventDto.InitData()
		if err := validate.Validate(eventDto); err != nil {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		event := domain.Event{Timestamp: eventDto.Timestamp, SensorSerialNumber: eventDto.SensorSerialNumber, SensorID: eventDto.SensorID, Payload: eventDto.Payload}

		if err := uc.Event.EventRepo.SaveEvent(ctx, &event); err != nil {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.Status(http.StatusCreated)
	}
}
