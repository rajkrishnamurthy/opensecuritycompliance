package rule

import (
	"cowlibrary/applications"
	"cowlibrary/constants"
	rule "cowlibrary/rule"
	cowlibutils "cowlibrary/utils"
	"cowlibrary/vo"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"cowctl/utils"
	"cowctl/utils/validationutils"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.NoArgs,

		Use:   "rule",
		Short: "Export the rule",
		Long:  "Export the rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd)
		},
	}

	// We can reset the flag by nonPersistentFlag.Changed = false

	cmd.Flags().String("rules-path", "", "path of the rules folder.")
	cmd.Flags().String("rule-name", "", "rule name.")
	cmd.Flags().String("downloads-path", "", "path of the downloads folder.")
	cmd.Flags().String("tasks-path", "", "path of the tasks.")
	cmd.Flags().String("config-path", "", "path for the configuration file.")
	cmd.Flags().String("export-file-type", "tar", "export file type(zip/tar)")
	cmd.Flags().String("publish-name", "", "this name will be used to publish")
	cmd.Flags().String("publish-description", "", "this field will be used to publish as description")
	cmd.Flags().String("client-id", "", "Id which is generated by the user")
	cmd.Flags().String("client-secret", "", "Secret key which is generated by user")
	cmd.Flags().String("sub-domain", "", "where to publish? dev/partner. default partner ")
	cmd.Flags().String("user-domain", "", "Where should the content be published if the user is connected to multiple domains? Default: Primary domain.")
	cmd.Flags().String("host", "", "where to publish? (eg:dev.compliancecow.live) ")
	cmd.Flags().String("catalog", "", "search in globalcatalog/rules only for the rule")
	cmd.Flags().Bool("can-override", false, "rule already exists in system")
	cmd.Flags().Bool("binary", false, "whether using cowctl binary")
	cmd.Flags().Bool("publish-app", false, "whether using cowctl binary")
	cmd.Flags().Bool("publish-rule", false, "whether using cowctl binary")
	cmd.Flags().String("app-name", "", "application name")
	cmd.Flags().String("rule-path", "", "path of the rule")

	return cmd
}

func runE(cmd *cobra.Command) error {
	additionalInfo, err := utils.GetAdditionalInfoFromCmd(cmd)
	if err != nil {
		return err
	}
	rulesPath := ``
	downloadsPath := additionalInfo.PolicyCowConfig.PathConfiguration.DownloadsPath
	localcatalogPath := additionalInfo.PolicyCowConfig.PathConfiguration.LocalCatalogPath

	additionalInfo.RulePublisher = &vo.RulePublisher{}
	var binaryEnabled bool
	var appPublish bool
	var appName string
	rulePublish := true

	if cmd.Flags().HasFlags() {
		binaryEnabled, _ = cmd.Flags().GetBool("binary")
		rulesPath = utils.GetFlagValueAndResetFlag(cmd, "rule-path", "")
		additionalInfo.RuleName = utils.GetFlagValueAndResetFlag(cmd, "rule-name", "")
		additionalInfo.DownloadsPath = utils.GetFlagValueAndResetFlag(cmd, "downloads-path", downloadsPath)
		additionalInfo.ExportFileType = utils.GetFlagValueAndResetFlag(cmd, "export-file-type", "tar")
		additionalInfo.RulePublisher.Name = utils.GetFlagValueAndResetFlag(cmd, "publish-name", "")
		additionalInfo.RulePublisher.Description = utils.GetFlagValueAndResetFlag(cmd, "publish-description", "")
		additionalInfo.ClientID = utils.GetFlagValueAndResetFlag(cmd, "client-id", "")
		additionalInfo.ClientSecret = utils.GetFlagValueAndResetFlag(cmd, "client-secret", "")
		additionalInfo.SubDomain = utils.GetFlagValueAndResetFlag(cmd, "sub-domain", "")
		additionalInfo.UserDomain = utils.GetFlagValueAndResetFlag(cmd, "user-domain", "")
		additionalInfo.Host = utils.GetFlagValueAndResetFlag(cmd, "host", "")

		if currentFlag := cmd.Flags().Lookup("can-override"); currentFlag != nil && currentFlag.Changed {
			if flagValue := currentFlag.Value.String(); cowlibutils.IsNotEmpty(flagValue) {
				currentFlag.Value.Set("false")
				additionalInfo.CanOverride, _ = strconv.ParseBool(flagValue)
			}
		}
		if currentFlag := cmd.Flags().Lookup("publish-rule"); currentFlag != nil && currentFlag.Changed {
			if flagValue := currentFlag.Value.String(); cowlibutils.IsNotEmpty(flagValue) {
				currentFlag.Value.Set("false")
				rulePublish, _ = strconv.ParseBool(flagValue)
			}
		}
		appName = utils.GetFlagValueAndResetFlag(cmd, "app-name", "")
		appPublish, _ = cmd.Flags().GetBool("publish-app")
	}

	defaultConfigPath := cowlibutils.IsDefaultConfigPath(constants.CowDataDefaultConfigFilePath)

	if cowlibutils.IsNotEmpty(additionalInfo.ClientID) || cowlibutils.IsNotEmpty(additionalInfo.ClientSecret) || cowlibutils.IsNotEmpty(additionalInfo.SubDomain) || cowlibutils.IsNotEmpty(additionalInfo.Host) {
		if !cowlibutils.IsValidCredentials(additionalInfo) {
			return errors.New("not a valid credentials")
		}
	}

	if !utils.IsValidExportFileType(additionalInfo.ExportFileType) {
		exportFileType, err := utils.GetConfirmationFromCmdPromptWithOptions("not a valid file type to export. Choose the type to provide the file type(default:tar):", "tar", []string{"tar", "zip"})
		if err != nil {
			return err
		}
		additionalInfo.ExportFileType = exportFileType
	}

	if cowlibutils.IsEmpty(additionalInfo.RuleName) {
		if !defaultConfigPath {
			return errors.New("Set the rule name by using the flag 'rule-name'")
		}
		pathPrefixes := []string{filepath.Join(additionalInfo.PolicyCowConfig.PathConfiguration.RulesPath, "*", "rule.json"),
			filepath.Join(additionalInfo.PolicyCowConfig.PathConfiguration.RulesPath, "*", "rule.yaml")}
		if cowlibutils.IsEmpty(rulesPath) {
			if !additionalInfo.GlobalCatalog {
				pathPrefixes = append(pathPrefixes, filepath.Join(localcatalogPath, "rules", "*", "rule.json"))
				pathPrefixes = append(pathPrefixes, filepath.Join(localcatalogPath, "rules", "*", "rule.yaml"))
				pathPrefixes = append(pathPrefixes, filepath.Join(localcatalogPath, "*", "rules", "*", "rule.json"))
				pathPrefixes = append(pathPrefixes, filepath.Join(localcatalogPath, "*", "rules", "*", "rule.yaml"))
			}

		}

		name, err := utils.GetValueAsFolderNameFromCmdPromptInCatalogs("Select a rule :", true, pathPrefixes, utils.ValidateString, additionalInfo)

		if err != nil {
			return err
		}

		if cowlibutils.IsEmpty(name) {
			return errors.New("rule name cannot be empty")
		}
		additionalInfo.RuleName = name
	}

	if cowlibutils.IsEmpty(rulesPath) {
		rulesPath = cowlibutils.GetRulePathFromCatalog(additionalInfo, additionalInfo.RuleName)
	}
	if cowlibutils.IsNotValidRulePath(rulesPath) {
		return fmt.Errorf("%s not valid rule path", rulesPath)
	}
	if cowlibutils.IsNotEmpty(rulesPath) {
		if cowlibutils.IsFolderNotExist(rulesPath) {
			pathFromCmd, err := utils.GetValueAsFilePathFromCmdPrompt("Enter a valid file path", true, utils.ValidateFilePath)
			if err != nil || cowlibutils.IsEmpty(pathFromCmd) {
				return err
			}
			rulesPath = pathFromCmd

		}
		additionalInfo.Path = rulesPath
	}
	if cowlibutils.IsEmpty(additionalInfo.RulePublisher.Name) {
		additionalInfo.RulePublisher.Name = additionalInfo.RuleName
	}

	err = GetValidRuleName("Unable to publish with the same name", additionalInfo, binaryEnabled)
	if err != nil {
		return err
	}

	if additionalInfo.TerminateFlow {
		return nil
	}

	if cowlibutils.IsEmpty(additionalInfo.RulePublisher.Description) && !additionalInfo.CanOverride {
		if !defaultConfigPath || binaryEnabled {
			return errors.New("Give the rule description to publish by using the flag 'publish-description'")
		}
		synthesizerNameFromCmd, err := utils.GetValueAsStrFromCmdPrompt("Rule Description to publish", false, utils.ValidateString)
		if err != nil {
			return err
		}
		additionalInfo.RulePublisher.Description = synthesizerNameFromCmd
	}

	isRuleAlreadyPublished := additionalInfo.CanOverride

	inputYAMLFileByts, err := os.ReadFile(filepath.Join(rulesPath, constants.TaskInputYAMLFile))
	if err == nil {
		var appInfo vo.TaskInputV2
		err = yaml.Unmarshal(inputYAMLFileByts, &appInfo)
		if err != nil {
			return fmt.Errorf("not a valid rule input structure. error :%s", err.Error())
		}
		applicationName := appInfo.UserObject.App.ApplicationName

		namePointer := &vo.CowNamePointersVO{}
		namePointer.Name = applicationName
		if binaryEnabled {
			namePointer.Name = appName
		}
		if cowlibutils.IsNotEmpty(namePointer.Name) {
			if err := publishApplication(namePointer, additionalInfo, binaryEnabled, appPublish); err != nil {
				return fmt.Errorf("error during publishing application : %s", err)
			}
		}

	}

	if !isRuleAlreadyPublished {
		additionalInfo.CanOverride = false
	}

	if rulePublish {
		err = rule.PublishRule(rulesPath, additionalInfo)
	}

	additionalInfo.GlobalCatalog = false

	return err

}

func GetValidRuleName(errorDesc string, additionalInfo *vo.AdditionalInfo, binaryEnabled bool) error {
	defaultConfigPath := cowlibutils.IsDefaultConfigPath(constants.CowDataDefaultConfigFilePath)

	isRuleAlreadyPresent, err := rule.IsRuleAlreadyPresent(additionalInfo.RulePublisher.Name, additionalInfo)
	if err != nil {
		return err
	}

	if isRuleAlreadyPresent && !additionalInfo.CanOverride {
		if !defaultConfigPath || binaryEnabled && !additionalInfo.CanOverride {
			return errors.New("Rule name is already present in the system. To want to override with new implementation, set the 'can-override' flag as true")
		}

		isConfirmed, err := utils.GetConfirmationFromCmdPrompt("The rule name is already present in the system, and it will be overridden with a new implementation. Do you want to go ahead?")

		if err != nil {
			return err
		}

		ruleErrorDesc := fmt.Sprintf("%s. Give a different name", errorDesc)

		if isConfirmed {
			additionalInfo.CanOverride = true
			return nil
		} else {
			isConfirmed, err := utils.GetConfirmationFromCmdPrompt("Do you want to publish it under a different name?")
			if err != nil {
				return err
			}
			if !isConfirmed {
				additionalInfo.TerminateFlow = true
				return nil
			}
			ruleErrorDesc = "Provide an alternative rule name "
		}

		ruleNameFromCmd, err := utils.GetValueAsStrFromCmdPrompt(ruleErrorDesc, true, validationutils.ValidateAlphaNumeric)
		if err != nil || cowlibutils.IsEmpty(ruleNameFromCmd) {
			return err
		}
		additionalInfo.RulePublisher.Name = ruleNameFromCmd
		return GetValidRuleName("Rule name already present", additionalInfo, binaryEnabled)
	}

	return nil

}

func publishApplication(namePointer *vo.CowNamePointersVO, additionalInfo *vo.AdditionalInfo, binaryEnabled bool, appPublish bool) error {
	appPath := additionalInfo.PolicyCowConfig.PathConfiguration.AppConnectionPath
	packageName := strings.ToLower(namePointer.Name)
	language, err := cowlibutils.GetApplicationLanguageFromRule(additionalInfo.RuleName, additionalInfo)
	if err != nil || cowlibutils.IsEmpty(language) {
		return fmt.Errorf("failed to get application language")
	}
	additionalInfo.Language = language

	if !binaryEnabled {
		linkedApplications, _ := applications.GetLinkedApplications(namePointer, additionalInfo)
		if len(linkedApplications) > 0 {
			for _, linkedApp := range linkedApplications {
				newNamePointer := &vo.CowNamePointersVO{Name: linkedApp.Name}
				if err := publishApplication(newNamePointer, additionalInfo, binaryEnabled, appPublish); err != nil {
					return err
				}
			}
		}
	}

	applicationValidatorResp, err := applications.GetAvailableApplications(namePointer.Name, additionalInfo)
	if err != nil {
		return err
	}
	appConnectionPath := filepath.Join(appPath, additionalInfo.Language, packageName)
	if additionalInfo.Language == "python" {
		appConnectionPath = filepath.Join(appPath, additionalInfo.Language, filepath.Base(appPath), packageName)
	}

	if !applicationValidatorResp.Valid {
		if cowlibutils.IsFolderExist(appConnectionPath) {
			errorDetails := applications.PublishApplication(namePointer, additionalInfo)
			if len(errorDetails) > 0 {
				return errors.New(errorDetails[0].Issue)
			}
		}
	} else {
		if binaryEnabled && !appPublish {
			return fmt.Errorf("The application '%s' is already published. Do you want to publish the application again?", namePointer.Name)
		}
		isConfirmed := true
		if !binaryEnabled {
			isConfirmed, err = utils.GetConfirmationFromCmdPrompt(fmt.Sprintf("The application '%s' is already published. Do you want to publish the application again?", namePointer.Name))
			if err != nil {
				return err
			}
		}

		if isConfirmed {
			if cowlibutils.IsFolderExist(appConnectionPath) {
				additionalInfo.CanOverride = true
				errorDetails := applications.PublishApplication(namePointer, additionalInfo)
				if len(errorDetails) > 0 {
					return errors.New(errorDetails[0].Issue)
				}
				d := color.New(color.FgCyan, color.Bold)
				d.Println("Hurray!.. Application Configuration has been published on behalf of you")
			}
		}
	}

	return nil
}
