package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/crestron-control-microservice/crestroncontrol"
	"github.com/byuoitav/crestron-control-microservice/helpers"
	"github.com/byuoitav/crestron-control-microservice/sigfile"
	"github.com/labstack/echo"
)

//PowerOn handles the power on command.
func PowerOn(context echo.Context) error {
	log.Printf("Powering on %s...", context.Param("address"))

	allSignals, err := sigfile.GetSignalsForAddress(context.Param("address"))
	if err != nil {
		return err
	}

	value := crestroncontrol.GetSignalConfigValue(context, "PowerOn")

	config := crestroncontrol.SignalConfigFile.Mapping["PowerOn"]

	err = helpers.SetState(allSignals["PowerOn"].MemAddr, value, context.Param("address"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, helpers.ReturnError(err))
	}

	if config.HighLow {
		err = helpers.SetState(allSignals["PowerOn"].MemAddr, "0", context.Param("address"))

		if err != nil {
			return context.JSON(http.StatusInternalServerError, helpers.ReturnError(err))
		}
	}

	log.Printf("Done")
	return nil
}

//Standby handles the standby command
func Standby(context echo.Context) error {
	log.Printf("Powering off %s...", context.Param("address"))

	log.Printf("Done")
	return nil
}

//SwitchInput handles the SwitchInput command
func SwitchInput(context echo.Context) error {
	log.Printf("Switching input for %s to %s ...", context.Param("address"), context.Param("port"))
	// address := context.Param("address")
	// port := context.Param("port")

	log.Printf("Done")
	return nil
}

//SetVolume handles the SetVolume command
func SetVolume(context echo.Context) error {
	// address := context.Param("address")
	// value := context.Param("value")

	// log.Printf("Setting volume for %s to %v...", address, value)

	log.Printf("Done")
	return nil
}

//VolumeUnmute hanldes the unmute command
func VolumeUnmute(context echo.Context) error {
	// address := context.Param("address")
	// log.Printf("Unmuting %s...", address)

	log.Printf("Done")
	return nil
}

//VolumeMute handles the mute command
func VolumeMute(context echo.Context) error {
	log.Printf("Muting %s...", context.Param("address"))

	log.Printf("Done")
	return nil
}

//BlankDisplay handles the blank command
func BlankDisplay(context echo.Context) error {
	return Standby(context)
}

//UnblankDisplay handles the unblank display
func UnblankDisplay(context echo.Context) error {
	return PowerOn(context)
}

//GetVolume handles the request for volume levels
func GetVolume(context echo.Context) error {
	log.Printf("Getting volume for %s...", context.Param("address"))

	log.Printf("Done")
	return nil
}
