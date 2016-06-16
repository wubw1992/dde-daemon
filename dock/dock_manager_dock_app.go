package dock

import (
	"errors"
	"gir/gio-2.0"
	"io/ioutil"
	"os"
	"path/filepath"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
	"text/template"
)

const dockedItemTemplate string = `[Desktop Entry]
Name={{ .Name }}
Exec={{ .Exec }}
Icon={{ .Icon }}
Type=Application
Terminal=false
StartupNotify=false
`

type dockedItemInfo struct {
	Name, Icon, Exec string
}

func createScratchDesktopFile(id, title, icon, cmd string) error {
	logger.Debugf("create scratch file for %q", id)
	err := os.MkdirAll(scratchDir, 0775)
	if err != nil {
		logger.Warning("create scratch directory failed:", err)
		return err
	}
	f, err := os.OpenFile(filepath.Join(scratchDir, id+".desktop"),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0744)
	if err != nil {
		logger.Warning("Open file for write failed:", err)
		return err
	}

	defer f.Close()
	temp := template.Must(template.New("docked_item_temp").Parse(dockedItemTemplate))
	dockedItem := dockedItemInfo{title, icon, cmd}
	logger.Debugf("dockedItem: %#v", dockedItem)
	err = temp.Execute(f, dockedItem)
	if err != nil {
		return err
	}
	return nil
}

func removeScratchFiles(id string) {
	extList := []string{"desktop", "sh", "png"}
	for _, ext := range extList {
		file := filepath.Join(scratchDir, id+"."+ext)
		if dutils.IsFileExist(file) {
			logger.Debugf("remove scratch file %q", file)
			err := os.Remove(file)
			if err != nil {
				logger.Warning("remove scratch file %q failed:", file, err)
			}
		}
	}
}

func createScratchDesktopFileWithAppEntry(entry *AppEntry) string {
	appId := "docked:" + entry.innerId

	if entry.appInfo != nil {
		desktopFile := entry.appInfo.GetFilePath()
		newPath := filepath.Join(scratchDir, appId+".desktop")
		// try link
		err := os.Link(desktopFile, newPath)
		if err != nil {
			logger.Warning("link failed try copy file contents")
			err = copyFileContents(desktopFile, newPath)
		}
		if err == nil {
			return appId
		} else {
			logger.Warning(err)
		}
	}

	title := entry.current.getDisplayName()
	// icon
	icon := entry.current.getIcon()
	if strings.HasPrefix(icon, "data:image") {
		path, err := dataUriToFile(icon, filepath.Join(scratchDir, appId+".png"))
		if err != nil {
			logger.Warning(err)
			icon = ""
		} else {
			icon = path
		}
	}
	if icon == "" {
		icon = "application-default-icon"
	}

	// cmd
	scriptContent := "#!/bin/sh\n" + entry.exec
	scriptFile := filepath.Join(scratchDir, appId+".sh")
	ioutil.WriteFile(scriptFile, []byte(scriptContent), 0744)
	cmd := scriptFile + " %U"

	err := createScratchDesktopFile(appId, title, icon, cmd)
	if err != nil {
		logger.Warning("createScratchDesktopFile failed:", err)
		return ""
	}
	return appId
}

func (m *DockManager) getDockedEntries() []*AppEntry {
	var dockedEntries []*AppEntry
	for _, entry := range m.Entries {
		if entry.appInfo != nil && entry.IsDocked == true {
			dockedEntries = append(dockedEntries, entry)
		}
	}
	return dockedEntries
}

func (m *DockManager) getDockedAppEntryByDesktopFilePath(desktopFilePath string) (*AppEntry, error) {
	dockedEntries := m.getDockedEntries()
	// same file
	for _, entry := range dockedEntries {
		file := entry.appInfo.GetFilePath()
		if file == desktopFilePath {
			return entry, nil
		}
	}

	// hash equal
	appInfo := NewAppInfoFromFile(desktopFilePath)
	if appInfo == nil {
		return nil, errors.New("Invalid desktopFilePath")
	}
	hash := appInfo.innerId
	appInfo.Destroy()
	for _, entry := range dockedEntries {
		if entry.appInfo.innerId == hash {
			return entry, nil
		}
	}
	return nil, nil
}

func (m *DockManager) saveDockedApps() {
	var list []string
	for _, entry := range m.getDockedEntries() {
		list = append(list, entry.appInfo.GetId())
	}
	m.DockedApps.Set(list)
}

func (m *DockManager) dockEntry(entry *AppEntry) bool {
	entry.dockMutex.Lock()
	defer entry.dockMutex.Unlock()

	if entry.IsDocked {
		logger.Warningf("dockEntry failed: entry %v is docked", entry.Id)
		return false
	}
	needScratchDesktop := false
	if entry.appInfo == nil {
		logger.Debug("dockEntry: entry.appInfo is nil")
		needScratchDesktop = true
	} else {
		// try create appInfo by desktopId
		desktopId := entry.appInfo.GetDesktopId()
		appInfo := gio.NewDesktopAppInfo(desktopId)
		if appInfo != nil {
			appInfo.Unref()
		} else {
			logger.Debugf("dockEntry: gio.NewDesktopAppInfo failed: desktop id %q", desktopId)
			needScratchDesktop = true
		}
	}

	logger.Debug("dockEntry: need scratch desktop?", needScratchDesktop)
	if needScratchDesktop {
		appId := createScratchDesktopFileWithAppEntry(entry)
		if appId != "" {
			entry.appInfo = NewAppInfo(appId)
			entryOldInnerId := entry.innerId
			entry.innerId = entry.appInfo.innerId
			logger.Debug("dockEntry: createScratchDesktopFile successed, entry use new innerId", entry.innerId)

			if strings.HasPrefix(entryOldInnerId, windowHashPrefix) {
				// entryOldInnerId is window hash
				m.desktopWindowsMapCacheManager.AddKeyValue(entry.innerId, entryOldInnerId)
				m.desktopWindowsMapCacheManager.AutoSave()
			}

			m.desktopHashFileMapCacheManager.SetKeyValue(entry.innerId, entry.appInfo.GetFilePath())
			m.desktopHashFileMapCacheManager.AutoSave()
		} else {
			logger.Warning("createScratchDesktopFileWithAppEntry failed")
			return false
		}
	}

	entry.setIsDocked(true)
	entry.updateMenu()
	m.saveDockedApps()
	return true
}

func (m *DockManager) undockEntry(entry *AppEntry) {
	entry.dockMutex.Lock()
	defer entry.dockMutex.Unlock()

	if !entry.IsDocked {
		logger.Warningf("undockEntry failed: entry %v is not docked", entry.Id)
		return
	}

	if entry.appInfo == nil {
		logger.Warning("undockEntry failed: entry.appInfo is nil")
		return
	}
	appId := entry.appInfo.GetId()

	if !entry.hasWindow() {
		m.removeAppEntry(entry)
	} else {
		dir := filepath.Dir(entry.appInfo.GetFilePath())
		if dir == scratchDir {
			removeScratchFiles(appId)
			// Re-identify Window
			if entry.current != nil {
				var newAppInfo *AppInfo
				entry.innerId, newAppInfo = m.identifyWindow(entry.current)
				entry.setAppInfo(newAppInfo)
			}
		}

		entry.setIsDocked(false)
		entry.updateMenu()
	}
	m.saveDockedApps()
}
