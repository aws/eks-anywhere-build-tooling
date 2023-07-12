package builder

type File struct {
	Source      string `json:"src"`
	Destination string `json:"dest"`
	Owner       string `json:"owner"`
	Group       string `json:"group"`
	Mode        int    `json:"mode"`
}

type AdditionalFilesConfig struct {
	FileVars
	FilesAnsibleConfig
}

type FileVars struct {
	AdditionalFiles     string `json:"additional_files"`
	AdditionalFilesList []File `json:"additional_files_list"`
}

type FilesAnsibleConfig struct {
	AnsibleExtraVars string `json:"ansible_extra_vars"`
	CustomRole       string `json:"custom_role"`
	CustomRoleNames  string `json:"custom_role_names"`
}

type AdditionalFiles interface {
	ProcessAdditionalFiles()
}

func (afc *AdditionalFilesConfig) ProcessAdditionalFiles() {
	if len(afc.FileVars.AdditionalFilesList) != 0 {
		afc.FileVars.AdditionalFiles = "true"
		afc.FilesAnsibleConfig.CustomRole = "true"
		afc.FilesAnsibleConfig.CustomRoleNames = DefaultAMICustomRoleNames
		if afc.FilesAnsibleConfig.AnsibleExtraVars == "" {
			afc.FilesAnsibleConfig.AnsibleExtraVars = DefaultAMIAnsibleExtraVars
		}
	}
}
