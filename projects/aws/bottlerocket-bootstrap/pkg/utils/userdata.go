package utils

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const userDataFile = "/.bottlerocket/host-containers/current/user-data"

type WriteFile struct {
	Path        string
	Owner       string
	Permissions string
	Content     string
}

type UserData struct {
	WriteFiles []WriteFile `yaml:"write_files"`
	RunCmd     string      `yaml:"runcmd"`
}

func ParseUserData() (*UserData, error) {
	fmt.Println("Parsing userdata")
	// read userdata from the file
	data, err := ioutil.ReadFile(userDataFile)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading user data file")
	}

	userData := &UserData{}
	err = yaml.Unmarshal(data, userData)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling user data")
	}
	fmt.Printf("\n%+v\n", userData)
	return userData, nil
}

func WriteUserDataFiles(userData *UserData) error {
	fmt.Println("Writing userdata write files")
	for _, file := range userData.WriteFiles {
		if file.Permissions == "" {
			file.Permissions = "0640"
		}
		perm, err := strconv.ParseInt(file.Permissions, 8, 64)
		if err != nil {
			return errors.Wrap(err, "Error converting string to int for permissions")
		}
		dir := filepath.Dir(file.Path)
		err = os.MkdirAll(dir, 0640)
		if err != nil {
			return errors.Wrap(err, "Error creating directories")
		}
		err = ioutil.WriteFile(file.Path, []byte(file.Content), fs.FileMode(perm))
		if err != nil {
			return errors.Wrapf(err, "Error creating file: %s", file.Path)
		}
		// get owner
		owners := strings.Split(file.Owner, ":")
		owner := owners[0]
		userDetails, err := user.Lookup(owner)
		if err != nil {
			return errors.Wrap(err, "Error getting user/group details ")
		}
		uid, _ := strconv.Atoi(userDetails.Uid)
		gid, _ := strconv.Atoi(userDetails.Gid)
		err = syscall.Chown(file.Path, uid, gid)
		if err != nil {
			return errors.Wrap(err, "Error running chown to set owners/groups")
		}
	}
	return nil
}
