package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/crestron-control-microservice/crestroncontrol"
	"github.com/byuoitav/crestron-control-microservice/helpers"
	"github.com/byuoitav/crestron-control-microservice/sigfile"
	"github.com/labstack/echo"
)

func PowerOn(context echo.Context) error {
	log.Printf("Powering on %s...", context.Param("address"))

	allSignals, err := sigfile.GetSignalsForAddress(context.Param("address"))
	if err != nil {
		return err
	}

	value := crestroncontrol.GetSignalConfigValue(context, "PowerOn")

	err = helpers.SetState(allSignals["PowerOn"].MemAddr, value, context.Param("address"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, helpers.ReturnError(err))
	}

	log.Printf("Done")
	return nil
}

func Standby(context echo.Context) error {
	log.Printf("Powering off %s...", context.Param("address"))

	log.Printf("Done")
	return nil
}

func SwitchInput(context echo.Context) error {
	log.Printf("Switching input for %s to %s ...", context.Param("address"), context.Param("port"))
	// address := context.Param("address")
	// port := context.Param("port")

	log.Printf("Done")
	return nil
}

func SetVolume(context echo.Context) error {
	// address := context.Param("address")
	// value := context.Param("value")

	// log.Printf("Setting volume for %s to %v...", address, value)

	log.Printf("Done")
	return nil
}

func VolumeUnmute(context echo.Context) error {
	// address := context.Param("address")
	// log.Printf("Unmuting %s...", address)

	log.Printf("Done")
	return nil
}

func VolumeMute(context echo.Context) error {
	log.Printf("Muting %s...", context.Param("address"))

	log.Printf("Done")
	return nil
}

func BlankDisplay(context echo.Context) error {
	return Standby(context)
}

func UnblankDisplay(context echo.Context) error {
	return PowerOn(context)
}

func GetVolume(context echo.Context) error {
	log.Printf("Getting volume for %s...", context.Param("address"))

	log.Printf("Done")
	return nil
}
