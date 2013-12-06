package main

import (
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
)

type ExtDevManager struct {
	DevInfoList []ExtDeviceInfo
}

type ExtDeviceInfo struct {
	DevicePath string
	DeviceType string
}

type MouseEntry struct {
	UseHabit       dbus.Property `access:"readwrite"`
	MoveSpeed      dbus.Property `access:"readwrite"`
	MoveAccuracy   dbus.Property `access:"readwrite"`
	ClickFrequency dbus.Property `access:"readwrite"`
	DeviceID       string
}

type TPadEntry struct {
	TPadEnable     dbus.Property `access:"readwrite"`
	UseHabit       dbus.Property `access:"readwrite"`
	MoveSpeed      dbus.Property `access:"readwrite"`
	MoveAccuracy   dbus.Property `access:"readwrite"`
	ClickFrequency dbus.Property `access:"readwrite"`
	DragDelay      dbus.Property `access:"readwrite"`
	DeviceID       string
}

type KeyboardEntry struct {
	RepeatDelay    dbus.Property `access:"readwrite"`
	RepeatSpeed    dbus.Property `access:"readwrite"`
	CursorBlink    dbus.Property `access:"readwrite"`
	DisableTPad    dbus.Property `access:"readwrite"`
	KeyboardLayout dbus.Property `access:"readwrite"`
	DeviceID       string
}

const (
	_EXT_DEV_NAME = "com.deepin.daemon.ExtDevManager"
	_EXT_DEV_PATH = "/com/deepin/daemon/ExtDevManager"
	_EXT_DEV_IFC  = "com.deepin.daemon.ExtDevManager"

	_EXT_ENTRY_PATH = "/com/deepin/daemon/ExtDevManager/"
	_EXT_ENTRY_IFC  = "com.deepin.daemon.ExtDevManager."

	_KEYBOARD_REPEAT_SCHEMA = "org.gnome.settings-daemon.peripherals.keyboard"
	_LAYOUT_SCHEMA          = "org.gnome.libgnomekbd.keyboard"
	_DESKTOP_INFACE_SCHEMA  = "org.gnome.desktop.interface"
	_MOUSE_SCHEMA           = "org.gnome.settings-daemon.peripherals.mouse"
	_TPAD_SCHEMA            = "org.gnome.settings-daemon.peripherals.touchpad"
)

var (
	mouseGSettings     *gio.Settings
	tpadGSettings      *gio.Settings
	infaceGSettings    *gio.Settings
	layoutGSettings    *gio.Settings
	keyRepeatGSettings *gio.Settings
)

func InitGSettings() bool {
	mouseGSettings = gio.NewSettings(_MOUSE_SCHEMA)
	tpadGSettings = gio.NewSettings(_TPAD_SCHEMA)
	infaceGSettings = gio.NewSettings(_DESKTOP_INFACE_SCHEMA)
	layoutGSettings = gio.NewSettings(_LAYOUT_SCHEMA)
	keyRepeatGSettings = gio.NewSettings(_KEYBOARD_REPEAT_SCHEMA)
	return true
}

func NewKeyboardEntry() *KeyboardEntry {
	keyboard := &KeyboardEntry{}

	keyboard.DeviceID = "Keyboard"
	keyboard.RepeatDelay = property.NewGSettingsProperty(keyboard,
		"RepeatDelay", keyRepeatGSettings, "delay")
	keyboard.RepeatSpeed = property.NewGSettingsProperty(keyboard,
		"RepeatSpeed", keyRepeatGSettings, "repeat-interval")
	keyboard.DisableTPad = property.NewGSettingsProperty(keyboard,
		"DisableTPad", tpadGSettings, "disable-while-typing")
	keyboard.CursorBlink = property.NewGSettingsProperty(keyboard,
		"CursorBlink", infaceGSettings, "cursor-blink-time")
	keyboard.KeyboardLayout = property.NewGSettingsProperty(keyboard,
		"KeyboardLayout", layoutGSettings, "layouts")
	return keyboard
}

func (keyboard *KeyboardEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + keyboard.DeviceID,
		_EXT_ENTRY_IFC + keyboard.DeviceID,
	}
}

func NewMouseEntry() *MouseEntry {
	mouse := &MouseEntry{}

	mouse.DeviceID = "Mouse"
	mouse.UseHabit = property.NewGSettingsProperty(mouse,
		"UseHabit", mouseGSettings, "left-handed")
	mouse.MoveSpeed = property.NewGSettingsProperty(mouse,
		"MoveSpeed", mouseGSettings, "motion-acceleration")
	mouse.MoveAccuracy = property.NewGSettingsProperty(mouse,
		"MoveAccuracy", mouseGSettings, "motion-threshold")
	mouse.ClickFrequency = property.NewGSettingsProperty(mouse,
		"ClickFrequency", mouseGSettings, "double-click")

	return mouse
}

func (mouse *MouseEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + mouse.DeviceID,
		_EXT_ENTRY_IFC + mouse.DeviceID,
	}
}

func NewTPadEntry() *TPadEntry {
	tpad := &TPadEntry{}

	tpad.DeviceID = "TouchPad"
	tpad.TPadEnable = property.NewGSettingsProperty(tpad,
		"TPadEnable", tpadGSettings, "touchpad-enabled")
	tpad.UseHabit = property.NewGSettingsProperty(tpad,
		"UseHabit", tpadGSettings, "left-handed")
	tpad.MoveSpeed = property.NewGSettingsProperty(tpad,
		"MoveSpeed", tpadGSettings, "motion-acceleration")
	tpad.MoveAccuracy = property.NewGSettingsProperty(tpad,
		"MoveAccuracy", tpadGSettings, "motion-threshold")
	tpad.DragDelay = property.NewGSettingsProperty(tpad,
		"DragDelay", mouseGSettings, "drag-threshold")
	tpad.ClickFrequency = property.NewGSettingsProperty(tpad,
		"ClickFrequency", mouseGSettings, "double-click")

	return tpad
}

func (tpad *TPadEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + tpad.DeviceID,
		_EXT_ENTRY_IFC + tpad.DeviceID,
	}
}

func NewExtDevManager() *ExtDevManager {
	return &ExtDevManager{}
}

func (dev *ExtDevManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_EXT_DEV_NAME, _EXT_DEV_PATH, _EXT_DEV_IFC}
}
