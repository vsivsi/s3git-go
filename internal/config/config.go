package config

import (
	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
)

const S3GIT_CONFIG = ".s3git.config"
const CONFIG = "config"

var Config ConfigObject

type ConfigObject struct {
	S3gitVersion     int    `json:"s3gitVersion"`
	S3gitType        string `json:"s3gitType"` // config

	S3gitCasPath	 string `json:"s3gitCasPath"`

	S3gitS3Bucket    string `json:"s3gitS3Bucket"`
	S3gitS3Region    string `json:"s3gitS3Region"`
	S3gitS3AccessKey string `json:"s3gitS3AccessKey"`
	S3gitS3SecretKey string `json:"s3gitS3SecretKey"`

	S3gitMinioEndpoint string `json:"s3gitMinioEndpoint"`
	S3gitMinioInsecure bool   `json:"s3gitMinioInsecure"`

	S3gitAcdRefreshToken string `json:"s3gitAcdRefreshToken"`
}

func getConfigFile(dir string) string {
	return dir + "/" + S3GIT_CONFIG
}

func LoadConfig(dir string) (bool, error) {

	configFile := getConfigFile(dir)
	byteConfig, err := ioutil.ReadFile(configFile)
	if err != nil { // No config found, fine, ignore
		return false, nil
	}

	dec := json.NewDecoder(strings.NewReader(string(byteConfig)))
	if err := dec.Decode(&Config); err != nil {
		return false, err
	}

	return true, nil
}

func SaveConfig(dir string) error {

	bucket := getDefaultValue("test", "S3GIT_S3_BUCKET")
	region := getDefaultValue("us-east-1", "S3GIT_S3_REGION")
	accessKey := getDefaultValue("", "S3GIT_S3_ACCESS_KEY")
	secretKey := getDefaultValue("", "S3GIT_S3_SECRET_KEY")

	configObject := ConfigObject{S3gitVersion: 1, S3gitType: CONFIG, S3gitCasPath:dir,
		S3gitS3Bucket: bucket, S3gitS3Region: region, S3gitS3AccessKey: accessKey, S3gitS3SecretKey: secretKey, S3gitMinioEndpoint: "localhost:9000", S3gitMinioInsecure: true }

	buf := new(bytes.Buffer)

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(configObject); err != nil {
		return err
	}

	err := ioutil.WriteFile(getConfigFile(dir), buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func getDefaultValue(def, envName string) string {

	val := def

	envVal := os.Getenv(envName)
	if envVal != "" {
		val = envVal
	}
	return val
}