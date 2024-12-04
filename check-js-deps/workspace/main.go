package workspace

import (
	"check-js-deps/sets"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// "sync"

	"github.com/gobwas/glob"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

type Workspace struct {
	Packages []string
}

type Dependency struct {
	Version   string `yaml:"version"`
	Specifier string `yaml:"specifier"`
}

type Importer struct {
	DevDependencies map[string]Dependency `yaml:"devDependencies"`
	Dependencies    map[string]Dependency `yaml:"dependencies"`
}

type PnpmLockFile struct {
	Importers map[string]Importer `yaml:"importers"`
}

type PackageJson struct {
	Version string `json:"version"`
}

func readPackageJson(path string) (PackageJson, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return PackageJson{}, err
	}
	defer file.Close()

	// Read the file
	byteValue, err := io.ReadAll(file)
	var payload PackageJson
	if err != nil {
		return PackageJson{}, err
	}
	err = json.Unmarshal(byteValue, &payload)
	if err != nil {
		return PackageJson{}, err
	}
	return payload, nil
}

func readPnpmLockfile(path string) (PnpmLockFile, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return PnpmLockFile{}, err
	}
	defer file.Close()

	// Read the file
	byteValue, err := io.ReadAll(file)
	var payload PnpmLockFile
	if err != nil {
		return PnpmLockFile{}, err
	}
	err = yaml.Unmarshal(byteValue, &payload)
	if err != nil {
		return PnpmLockFile{}, err
	}
	return payload, nil
}

func readMainWorkspace(path string) (Workspace, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return Workspace{}, err
	}
	defer file.Close()

	// Read the file
	byteValue, err := io.ReadAll(file)
	var payload Workspace
	if err != nil {
		return Workspace{}, err
	}
	err = yaml.Unmarshal(byteValue, &payload)
	if err != nil {
		return Workspace{}, err
	}
	return payload, nil
}

func debugPrint(args ...interface{}) {
	_, isOk := os.LookupEnv("DEBUG_GO")
	if isOk {
		fmt.Println(args...)
	}
}

func excludePackage(notPackages []string, path string) bool {
	for _, notPackage := range notPackages {
		g := glob.MustCompile(notPackage[1:])
		if !g.Match(path) {
			return false
		}
	}
	return true
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func isLink(input string) bool {
	re := regexp.MustCompile(`^(.*?)\:`) // regex to capture everything up to the first '('
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1] == "link"
	}
	return false
}

func getLinkAbsPath(originalPath string, fullLinkStr string) (string, error) {
	splitedStrings := strings.Split(fullLinkStr, ":")
	path := splitedStrings[1]
	newLinkPath, pathError := filepath.Abs(originalPath + "/" + path)
	if pathError != nil {
		return "", pathError
	}
	return newLinkPath, nil
}

func containsGitOrGithub(input string) bool {
	pattern := regexp.MustCompile(`\b(git|github)\b`)
	return pattern.MatchString(input)
}

// GetStringUntilParenthesis returns the substring up until the first '(' using regex.
func extractVersion(input string) string {
	re := regexp.MustCompile(`^(.*?)\(`) // regex to capture everything up to the first '('
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return input // return the original string if '(' is not found
}

func checkPackage(path string, packageName string, pkgDetails Dependency, links *sets.DoubleSet) (bool, error) {
	packagePath := path + "/node_modules/" + packageName
	if !fileExists(packagePath) {
		// fmt.Println("Missing package " + packageName)
		return false, errors.New("Missing package " + packageName + " 'pnpm install' should help")
	}
	if isLink(pkgDetails.Version) {
		newLinkPath, pathError := getLinkAbsPath(path, pkgDetails.Version)
		if pathError != nil {
			return false, pathError
		}
		links.Add(newLinkPath)
		return true, nil
	}
	if containsGitOrGithub(pkgDetails.Version) {
		debugPrint("skipping " + pkgDetails.Version)
		return true, nil
	}
	version := extractVersion(pkgDetails.Version)
	pkgJson, err := readPackageJson(packagePath + "/package.json")
	if err != nil {
		return false, err
	}
	if version == pkgJson.Version {
		debugPrint("exact Version! ", packageName, version)
	} else {
		return false, errors.New("not same version " + packageName + " - " + version + ", pnpm install should help")
	}
	return true, nil
}

// func CheckProject(path string, links *sets.DoubleSet, wg *sync.WaitGroup) (bool, error) {
func CheckProject(ctx context.Context, path string, links *sets.DoubleSet, g *errgroup.Group) (bool, error) {
	// func CheckProject(path string, links *sets.DoubleSet) (bool, error) {
	select {
	case <-ctx.Done():
		return true, nil
	default:
		if links.HasBeenChecked(path) {
			return true, nil
		}
		links.Check(path)
		filePath := path + "/pnpm-lock.yaml"
		debugPrint(filePath)
		if !fileExists(filePath) {
			return false, errors.New("Missing pnpm lock file (" + filePath + "), try pnpm install")
		}
		lockFile, err := readPnpmLockfile(filePath)
		if err != nil {
			return false, err
		}
		debugPrint(lockFile.Importers["."].DevDependencies)
		for packageName, pkgDetails := range lockFile.Importers["."].DevDependencies {
			res, err := checkPackage(path, packageName, pkgDetails, links)
			if !res {
				return res, err
			}
		}
		debugPrint("*")
		for packageName, pkgDetails := range lockFile.Importers["."].Dependencies {
			res, err := checkPackage(path, packageName, pkgDetails, links)
			if !res {
				return res, err
			}
		}
		nonChecked := len(links.GetNoneChecked())
		debugPrint(nonChecked)
		for _, newPath := range links.GetNoneChecked() {
			// wg.Add(1)
			// go func(str string) {
			// 	defer wg.Done()
			// 	CheckProject(str, links, wg)
			// }(newPath)

			func(str string) {
				g.Go(func() error {
					_, err := CheckProject(ctx, newPath, links, g)
					return err
				})
			}(newPath)

			// return CheckProject(newPath, links)
		}
		return true, nil
	}
}

func findPackage(path string, notPackages []string) {
	f, err := os.Open(path)
	if err != nil {
		debugPrint(err)
		return
	}
	files, err := f.Readdir(0)
	if err != nil {
		debugPrint(err)
		return
	}

	for _, v := range files {
		if v.IsDir() {
			folder := path + v.Name()
			if excludePackage(notPackages, folder) {
				debugPrint("excluding " + folder + "&&&")
				continue
			}

			packJsonPath := folder + "/package.json"
			if fileExists(packJsonPath) {
				// checkPackage(folder)
				debugPrint(packJsonPath + "  Found!")
			} else {
				debugPrint(packJsonPath + "  Not Found!")
			}
		}
		// debugPrint(v.Name(), v.IsDir())
	}
}

func omitLast(path string) (string, string) {
	splitedStrings := strings.Split(path, "/")
	last := splitedStrings[len(splitedStrings)-1]
	joined := strings.Join(splitedStrings[:len(splitedStrings)-1], "/") + "/"
	return joined, last
}

func partition(input []string) ([]string, []string) {
	var with []string
	var without []string
	for _, str := range input {
		if len(str) > 0 && str[0] == '!' {
			with = append(with, str)
		} else {
			without = append(without, str)
		}
	}
	return without, with
}

func Read(path string) (Workspace, error) {
	workspace, err := readMainWorkspace(path)
	if err != nil {
		return Workspace{}, err
	}

	workspacePaths, notPackages := partition(sets.Unique(workspace.Packages))
	workspacePath, _ := omitLast(path)

	for _, projectsPathGlob := range workspacePaths {
		projectPath, last := omitLast(projectsPathGlob)
		switch last {
		case "*":
			{
				findPackage(workspacePath+projectPath, notPackages)

				fmt.Println("*")
			}
		case "**":
			{
				fmt.Println("*")
			}
		default:
			{
				fmt.Println("default")
			}
		}
	}
	return workspace, nil
}
