/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package inputdevices

import (
	"gir/gio-2.0"
	"sync"
)

var gsLocker sync.Mutex

func (kbd *Keyboard) handleGSettings() {
	kbd.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case kbdKeyRepeatEnable, kbdKeyRepeatDelay,
			kbdKeyRepeatInterval:
			kbd.setRepeat()
		case kbdKeyLayout:
			kbd.setLayout()
		case kbdKeyCursorBlink:
			kbd.setCursorBlink()
		case kbdKeyLayoutOptions:
			kbd.setOptions()
		case kbdKeyUserLayoutList:
			kbd.setGreeterLayoutList()
		}
	})
}

func (m *Mouse) handleGSettings() {
	m.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case mouseKeyLeftHanded:
			m.enableLeftHanded()
		case mouseKeyDisableTouchpad:
			m.disableTouchpad()
		case mouseKeyNaturalScroll:
			m.enableNaturalScroll()
		case mouseKeyMiddleButton:
			m.enableMidBtnEmu()
		case mouseKeyAcceleration:
			m.motionAcceleration()
		case mouseKeyThreshold:
			m.motionThreshold()
		case mouseKeyScaling:
			m.motionScaling()
		case mouseKeyDoubleClick:
			m.doubleClick()
		case mouseKeyDragThreshold:
			m.dragThreshold()
		}
	})
}

func (tp *TrackPoint) handleGSettings() {
	tp.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case trackPointKeyMidButton:
			tp.enableMiddleButton()
		case trackPointKeyMidButtonTimeout:
			tp.middleButtonTimeout()
		case trackPointKeyWheel:
			tp.enableWheelEmulation()
		case trackPointKeyWheelButton:
			tp.wheelEmulationButton()
		case trackPointKeyWheelTimeout:
			tp.wheelEmulationTimeout()
		case trackPointKeyWheelHorizScroll:
			tp.enableWheelHorizScroll()
		case trackPointKeyAcceleration:
			tp.motionAcceleration()
		case trackPointKeyThreshold:
			tp.motionThreshold()
		case trackPointKeyScaling:
			tp.motionScaling()
		}
	})
}

func (tpad *Touchpad) handleGSettings() {
	tpad.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case tpadKeyEnabled:
			tpad.enable(tpad.TPadEnable.Get())
		case tpadKeyLeftHanded:
			tpad.enableLeftHanded()
			tpad.enableTapToClick()
		case tpadKeyTapClick:
			tpad.enableTapToClick()
		case tpadKeyNaturalScroll:
			tpad.enableNaturalScroll()
		case tpadKeyScrollDelta:
			tpad.setScrollDistance()
		case tpadKeyEdgeScroll:
			tpad.enableEdgeScroll()
		case tpadKeyVertScroll, tpadKeyHorizScroll:
			tpad.enableTwoFingerScroll()
		case tpadKeyWhileTyping:
			tpad.disableWhileTyping()
		case tpadKeyAcceleration:
			tpad.motionAcceleration()
		case tpadKeyThreshold:
			tpad.motionThreshold()
		case tpadKeyScaling:
			tpad.motionScaling()
		}
	})
}

func (w *Wacom) handleGSettings() {
	w.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case wacomKeyLeftHanded:
			w.enableLeftHanded()
		case wacomKeyCursorMode:
			w.enableCursorMode()
		case wacomKeyUpAction:
			w.setKeyAction(btnNumUpKey, w.KeyUpAction.Get())
		case wacomKeyDownAction:
			w.setKeyAction(btnNumDownKey, w.KeyDownAction.Get())
		case wacomKeyDoubleDelta:
			w.setClickDelta()
		case wacomKeyPressureSensitive:
			w.setPressure()
		}
	})
}
