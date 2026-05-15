package variables

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strings"

	"github.com/l58193/terrabutler/internal/logger"
	"github.com/l58193/terrabutler/internal/settings"
	"github.com/l58193/terrabutler/internal/utils"

	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/spf13/afero"
)

func Generate_var_files(env string, fs afero.Fs) error {

	org := settings.Conf.String("general.organization")

	//Get file templates
	templates, err := afero.ReadDir(fs, utils.Paths["templates"])
	if err != nil {
		return errors.New("Error reading templates directory.")
	}

	// Initializing the data of the templates into a memory map, so it is compatible with tests
	templatesFilesData := make(map[string]string)
	for _, template := range templates {
		if !template.IsDir() {
			data, _ := afero.ReadFile(fs, utils.Paths["templates"]+"/"+template.Name())
			logger.Zap.Info("New templated loaded: " + template.Name() + " Data: " + string(data))
			templatesFilesData["/"+template.Name()] = string(data)
		}
	}

	//Initializing file_loader
	file_loader, err := loaders.NewMemoryLoader(templatesFilesData)
	if err != nil {
		return errors.New("Error initializing file_loader in templates dir, error: " + err.Error())
	}

	cfg := config.New()
	cfg.StrictUndefined = true

	//Get all the required fields
	sites := settings.Conf.Strings("sites.ordered")
	firebase_credentials := settings.Conf.String("environments.temporary.secrets.firebase_credentials")
	mail_password := settings.Conf.String("environments.temporary.secrets.mail_password")

	//Remove inception from sites
	if index := slices.Index(sites, "inception"); index != -1 {
		sites = append(sites[:index], sites[index+1:]...)
	}

	// In the Environment you can add the context
	environment := &exec.Environment{
		Context: exec.NewContext(map[string]any{
			"env":                         env,
			"generate_encrypted_password": generate_encrypted_password,
			"sites":                       sites,
			"mail_password":               mail_password,
			"firebase_credentials":        firebase_credentials}),
		//Its Required to us a ControlStructure or the render fails
		//Tests:             builtins.Tests,
		ControlStructures: builtins.ControlStructures,
		// The other fields or optional
		// Methods:           builtins.Methods,
		// Filters:           builtins.Filters,
	}

	//os.Chdir(utils.Paths["templates"])

	//For each template
	for _, template := range templates {

		if !template.IsDir() {

			logger.Zap.Debug("Template file: " + template.Name())

			temp, err := exec.NewTemplate(template.Name(), cfg, file_loader, environment)
			if err != nil {
				return errors.New("Error opening the template " + template.Name() + ", Error: " + err.Error())
			}

			output, err := temp.ExecuteToBytes(environment.Context)
			if err != nil {
				return errors.New("Error rendering template " + template.Name() + ", Error: " + err.Error())
			}

			name := strings.ReplaceAll(template.Name(), ".j2", "")

			//If the name is env
			if name == "env" {
				//Create new file and write the template output there
				f, _ := fs.Create(utils.Paths["variables"] + "/" + org + "-" + env + ".tfvars")

				l, err := f.Write(output)
				if l == 0 && err != nil {
					f.Close()
					return errors.New("An error has occurred writing to the file: " + err.Error())
				}
				err = f.Close()
			} else {
				//Create new file and write the template output there
				f, _ := fs.Create(utils.Paths["variables"] + "/" + org + "-" + env + "-" + name + ".tfvars")

				l, err := f.Write(output)
				if l == 0 && err != nil {
					f.Close()
					return errors.New("An error has occurred writing to the file: " + err.Error())
				}
				f.Close()
			}
		}
	}

	return nil

}

func generate_password(size int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
	charactersSplit := strings.Split(characters, "")
	password := ""

	logger.Zap.Debug(fmt.Sprintf("Size value: %d", size))

	for i := 0; i < size; i++ {
		password = password + charactersSplit[rand.Intn(len(charactersSplit))]
	}
	return password
}

// The encryption of the AWS is not implemented, so the encoding and decoding is useless for now
func encrypt_password(password string) string {
	//region := settings.String("environments.default.region")
	//key_id := settings.String("general.secrets_key_id")

	//This is Python version code....
	// environment = boto3.session.Session(profile_name=f"{ORG}-dev", region_name=REGION)
	//kms = environment.client("kms")
	//encrypted = kms.encrypt(KeyId=KEY_ID, Plaintext=password)
	//password_encrypted = encrypted[u'CiphertextBlob']

	passEncoded := base64.StdEncoding.EncodeToString([]byte(password))
	passDecoded, _ := base64.StdEncoding.DecodeString(passEncoded)
	return string(passDecoded)
}

func generate_encrypted_password(size int) string {
	password := generate_password(size)
	return encrypt_password(password)
}
