package main

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

const (
	KEY_NAME        = `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	VAL_NAME_SYSTEM = `SystemUsesLightTheme`
	VAL_NAME_APP    = `AppsUseLightTheme`
)

type DarkModeSwitcher struct {
}

func (dms DarkModeSwitcher) LightModeOn() (bool, error) {
	sys, err := dms.getRegValue(VAL_NAME_SYSTEM)
	if err != nil {
		return false, err
	}
	app, err := dms.getRegValue(VAL_NAME_APP)
	if err != nil {
		return false, err
	}
	if app != sys {
		return false, DarkModeSwitcherError{"System and App settings varies", nil}
	}

	return app, nil
}

func (dms DarkModeSwitcher) SetLightModeOn() error {
	err := dms.setRegValue(VAL_NAME_SYSTEM, 1)
	if err != nil {
		return err
	}
	dms.setRegValue(VAL_NAME_APP, 1)
	return err
}
func (dms DarkModeSwitcher) SetDarkModeOn() error {
	err := dms.setRegValue(VAL_NAME_SYSTEM, 0)
	if err != nil {
		return err
	}
	dms.setRegValue(VAL_NAME_APP, 0)
	return err

}
func (dms DarkModeSwitcher) ToggleLightMode() error {

	on, err := dms.LightModeOn()
	if err != nil {
		return err
	}
	if on {
		err = dms.SetDarkModeOn()
	} else {
		err = dms.SetLightModeOn()
	}
	return err
}

func (DarkModeSwitcher) getRegValue(name string) (bool, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, KEY_NAME, registry.QUERY_VALUE)

	if err != nil {
		return false, DarkModeSwitcherError{function: "OpenKey: ", reg_error: err}
	}

	val, _, err := k.GetIntegerValue(name)
	if err != nil {
		return false, DarkModeSwitcherError{function: "GetIntegerValue: ", reg_error: err}
	}
	//log.Print(VAL_NAME_SYSTEM, "=", val, " ", valtype)
	err = k.Close()
	if err != nil {
		return false, DarkModeSwitcherError{function: "Close: ", reg_error: err}
	}

	return val == 1, nil
}

func (DarkModeSwitcher) setRegValue(name string, val uint32) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, KEY_NAME, registry.QUERY_VALUE|registry.SET_VALUE)

	if err != nil {
		return DarkModeSwitcherError{function: "OpenKey(with SET_VALUE): ", reg_error: err}
	}

	err = k.SetDWordValue(name, val)
	if err != nil {
		return DarkModeSwitcherError{function: "SetDWordValue: ", reg_error: err}
	}

	err = k.Close()
	if err != nil {
		return DarkModeSwitcherError{function: "Close: ", reg_error: err}
	}

	return nil
}

type DarkModeSwitcherError struct {
	function  string
	reg_error error
}

func (e DarkModeSwitcherError) Error() string {
	return e.function + ": " + e.reg_error.Error()
}

func main() {

	dms := DarkModeSwitcher{}

	err := dms.ToggleLightMode()
	if err != nil {
		log.Fatal(err)
	}
}
