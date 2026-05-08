package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"

	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

func generate_var_files(env string) {

	org := settings.String("general.organization")

	//Get file templates
	templates, _ := os.ReadDir(paths["templates"])

	//Initializing file_loader
	file_loader, _ := loaders.NewFileSystemLoader(paths["templates"])

	cfg := config.New()
	cfg.StrictUndefined = true

	//Get all the required fields
	sites := settings.Strings("sites.ordered")
	firebase_credentials := settings.String("environments.temporary.secrets.firebase_credentials")
	mail_password := settings.String("environments.temporary.secrets.mail_password")

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

	logger.Debug(fmt.Sprintf("Sites: %v", sites))

	os.Chdir(paths["templates"])

	//For each template
	for _, template := range templates {

		if !template.IsDir() {

			logger.Debug("Template file: " + template.Name())

			temp, err := exec.NewTemplate(template.Name(), cfg, file_loader, environment)
			if err != nil {
				logger.Error("Error Opening the template " + template.Name())
			}

			output, err := temp.ExecuteToBytes(environment.Context)
			if err != nil {
				logger.Error("Error rendering template " + template.Name())
			}

			name := strings.ReplaceAll(template.Name(), ".j2", "")

			//If the name is env
			if name == "env" {
				//Create new file and write the template output there
				f, _ := os.Create(paths["variables"] + "/" + org + "-" + env + ".tfvars")

				l, err := f.Write(output)
				if l == 0 && err != nil {
					logger.Error(fmt.Sprint("An error has occurred writing to the file: ", err))
					f.Close()
					os.Exit(1)
				}
				err = f.Close()
			} else {
				//Create new file and write the template output there
				f, _ := os.Create(paths["variables"] + "/" + org + "-" + env + "-" + name + ".tfvars")

				l, err := f.Write(output)
				if l == 0 && err != nil {
					logger.Error(fmt.Sprint("An error has occurred writing to the file: ", err))
					f.Close()
					os.Exit(1)
				}
				err = f.Close()
			}
		}
	}

}

func generate_password(size int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
	charactersSplit := strings.Split(characters, "")
	password := ""

	logger.Debug(fmt.Sprintf("Size value: %d", size))

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
