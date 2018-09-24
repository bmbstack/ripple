package scripts

import (
	"errors"
	"fmt"
	"github.com/bmbstack/ripple/cmd/ripple/logger"
	"github.com/bmbstack/ripple/cmd/ripple/utils"
	"github.com/labstack/gommon/color"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// NewApplication create a new application with the appName
func NewApplication(appName string) {
	logger.Logger.Info("New application " + appName)
	appPath, err := getAppPath(appName)
	if err != nil {
		logger.Logger.Error(err.Error())
		return
	}

	if fileExists(appPath) {
		logger.Logger.Error(fmt.Sprintf("A folder already exists at path %s", appPath))
		return
	}

	// Copy the pristine new site over
	goPathArray := strings.Split(os.Getenv("GOPATH"), ":")
	goPath := goPathArray[0]
	templateAppPath := path.Join(goPath, "src", PACKAGE_TEMPLATES)
	err = copyApplication(templateAppPath, appPath)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error copying project %s", err))
		return
	}
}

//getAppPath According to the appName return abs path to this application
func getAppPath(appName string) (string, error) {
	var path string
	goPathString := os.Getenv("GOPATH")
	if goPathString == "" {
		return "", errors.New("GOPATH must be set to use the ripple tool")
	}
	goPathArray := strings.Split(goPathString, ":")
	goPath := goPathArray[0]
	inGoSrcPath := filepath.Join(goPath, "src", "*")
	currentPath, _ := os.Getwd()
	if matched, _ := filepath.Match(inGoSrcPath, currentPath); matched {
		path = filepath.Join(currentPath, appName)
	} else {
		path = filepath.Join(goPath, "src", appName)
	}
	return path, nil
}

func copyApplication(templateAppPath, appPath string) error {

	// Check that the folders up to the path exist, if not create them
	// Make directory
	err := os.MkdirAll(path.Dir(appPath), permissions)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("The project path could not be created: %s", err))
		return err
	}

	// Now recursively copy over the files from the original repo to project path
	logger.Logger.Info(fmt.Sprintf("Creating files at: %s", appPath))
	_, err = copyPath(templateAppPath, appPath)
	if err != nil {
		return err
	}

	// Delete the .git folder at that path
	gitPath := path.Join(appPath, ".git")
	logger.Logger.Info(fmt.Sprintf("Removing all at:%s", gitPath))
	err = os.RemoveAll(gitPath)
	if err != nil {
		return err
	}

	// Run git init to get a new git repo here
	logger.Logger.Info(fmt.Sprintf("Initializing new git repo at:%s", appPath))
	_, err = RunCommand("git", "init", appPath)
	if err != nil {
		return err
	}

	logCreateAppFiles(appPath)

	// Now reifyApplication
	return reifyApplication(templateAppPath, appPath)
}

// fileExists returns true if this file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// copyPath Copy a path to another one - at present this is unix only
// Unfortunately there is no simple facility for this in golang stdlib,
// so we use unix command (sorry windows!)
func copyPath(src, dst string) ([]byte, error) {
	// Replace this with an os independent version using filepath.Walk
	return RunCommand("cp", "-r", src, dst)
}

// projectPathRelative return the relative path to the appPath
func projectPathRelative(projectPath string) string {
	goSrc := os.Getenv("GOPATH") + "/src/"
	return strings.Replace(projectPath, goSrc, "", 1)
}

func logCreateAppFiles(appPath string) {
	logger.Logger.Info(fmt.Sprintf("Updating import paths to: %s", projectPathRelative(appPath)))
	err := filepath.Walk(appPath, func(file string, info os.FileInfo, err error) error {
		// If we have an err pass it up
		if err != nil {
			return err
		}
		// Deal with files only
		if !info.IsDir() {
			file = strings.Replace(file, path.Join(os.Getenv("GOPATH"), "src"), "", 1)
			logger.Logger.Debug(fmt.Sprintf("Create file: $GOPATH%s", file))
		}
		return nil
	})
	if err != nil {
		logger.Logger.Error("Create application files error!")
	}
}

// reifyApplication changes import refs within go files to the correct format
func reifyApplication(templateAppPath, appPath string) error {
	err := replaceExpressionInTemplates(templateAppPath, appPath, []string{".go", ".example"})
	if err != nil {
		return err
	}

	logger.Logger.Notice(fmt.Sprintf("Run command in bash: %s", color.Bold(color.Green("cd "+appPath))))
	logger.Logger.Notice(fmt.Sprintf("Run command in bash: %s", color.Bold(color.Green("go run main.go s"))))
	logger.Logger.Notice(fmt.Sprintf("Open this url: http://127.0.0.1:%s", HOST_PORT))
	return nil
}

// replaceExpressionInTemplates replace expression in the template
func replaceExpressionInTemplates(templateAppPath, appPath string, extentions []string) error {
	files, err := utils.CollectFiles(appPath, extentions)
	if err != nil {
		return err
	}

	// For each go file within project, make sure the refs are to the new site,
	// not to the template site
	relativeTemplateAppPath := projectPathRelative(templateAppPath)
	relativeAppPath := projectPathRelative(appPath)
	for _, f := range files {
		// Load the file, if it contains refs to goprojectpath, replace them with relative project path imports
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		// Substitutions - consider reifying instead if it is any more complex
		fileString := string(data)
		if strings.Contains(fileString, relativeTemplateAppPath) {
			fileString = strings.Replace(fileString, relativeTemplateAppPath, relativeAppPath, -1)
		}

		if strings.Contains(fileString, EXPRESSION_APP_NAME) {
			appName := utils.Substring(appPath, strings.LastIndex(appPath, "/")+1, len(appPath))
			fileString = strings.Replace(fileString, EXPRESSION_APP_NAME, appName, -1)
		}

		err = ioutil.WriteFile(f, []byte(fileString), permissions)
		if err != nil {
			return err
		}
	}

	return nil
}
