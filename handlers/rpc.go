package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/byuoitav/crestron-control-microservice/crestroncontrol"
	"github.com/byuoitav/crestron-control-microservice/helpers"
	"github.com/byuoitav/crestron-control-microservice/sigfile"
	"github.com/labstack/echo"
)

func HandleCommand(context echo.Context, commandName string) error {

	allSignals, err := sigfile.GetSignalsForAddress(context.Param("address"))

	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return context.JSON(http.StatusBadRequest, helpers.ReturnError(err))
	}

	log.Printf("Getting the signal sequence.")
	//Get the signal value
	values, err := crestroncontrol.GetSignalConfigSequence(context, commandName)
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, helpers.ReturnError(err))
	}
	log.Printf("%v steps.", len(values))

	//Run through our signal sequence
	for _, state := range values {
		log.Printf("Preparing to set %v to %v", state.SignalName, state.Value)
		var signal sigfile.Signal
		var ok bool

		if signal, ok = allSignals[state.SignalName]; !ok {
			err = fmt.Errorf("no signal for %v defined in the signal file", state.SignalName)

			log.Printf("ERROR: %v", err.Error())
			return context.JSON(http.StatusInternalServerError, helpers.ReturnError(err))
		}

		//set the state with the memory address from the config name.
		err = helpers.SetState(signal.MemAddr, state.Value, context.Param("address"))

		if err != nil {
			log.Printf("ERROR: %v", err.Error())
			return context.JSON(http.StatusBadRequest, helpers.ReturnError(err))
		}
	}

	log.Printf("Done")
	return nil
}

func Test(context echo.Context) error {
	log.Printf("Testing...")

	return HandleCommand(context, "Test")
}

//PowerOn handles the power on command.
func PowerOn(context echo.Context) error {
	log.Printf("Powering on %s...", context.Param("address"))

	return HandleCommand(context, "PowerOn")
}

//Standby handles the standby command
func Standby(context echo.Context) error {
	log.Printf("Powering off %s...", context.Param("address"))

	return HandleCommand(context, "Standby")
}

//Update checks for a current Sig file for the device and updates it if necessary
func Update(context echo.Context) error {
	log.Printf("Updating sig file for %s...", context.Param("address"))
	_, err := sigfile.Read(context.Param("address"))

	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, helpers.ReturnError(err))
	}
	log.Printf("Done.")
	return nil
}

//SwitchInput handles the SwitchInput command
func SwitchInput(context echo.Context) error {
	log.Printf("Switching input for %s to %s ...", context.Param("address"), context.Param("port"))
	//address := context.Param("address")
	//port := context.Param("port")

	return HandleCommand(context, "ChangeInput")
}

//SetVolume handles the SetVolume command
func SetVolume(context echo.Context) error {
	// address := context.Param("address")
	// value := context.Param("value")

	// log.Printf("Setting volume for %s to %v...", address, value)

	return HandleCommand(context, "SetVolume")
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
