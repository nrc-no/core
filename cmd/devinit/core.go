package devinit

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

func (c *Config) makeCore() error {

	var err error

	c.coreDbUsername = "core"
	c.coreDbName = "core"
	c.coreDbPassword, err = getOrCreateRandomSecretStr(32, CoreDir, "db-password")
	if err != nil {
		return err
	}

	c.dbUsers = append(c.dbUsers, dbUser{
		username: c.coreDbUsername,
		password: c.coreDbPassword,
		database: c.coreDbName,
	})

	const proxyUrl = "https://localhost:8443"

	coreConfig := map[string]interface{}{
		"serve": map[string]interface{}{
			"admin": map[string]interface{}{
				"cache": map[string]interface{}{
					"redis": map[string]interface{}{
						"password": c.redisPassword,
					},
				},
				"secrets": map[string]interface{}{
					"hash": []string{
						c.coreAdminApiHashKey,
					},
					"block": []string{
						c.coreAdminApiBlockKey,
					},
				},
				"urls": map[string]interface{}{
					"self": proxyUrl,
				},
			},
			"public": map[string]interface{}{
				"cache": map[string]interface{}{
					"redis": map[string]interface{}{
						"password": c.redisPassword,
					},
				},
				"secrets": map[string]interface{}{
					"hash": []string{
						c.coreApiHashKey,
					},
					"block": []string{
						c.coreApiBlockKey,
					},
				},
				"urls": map[string]interface{}{
					"self": proxyUrl,
				},
			},
			"login": map[string]interface{}{
				"cache": map[string]interface{}{
					"redis": map[string]interface{}{
						"password": c.redisPassword,
					},
				},
				"secrets": map[string]interface{}{
					"hash": []string{
						c.loginHashKey,
					},
					"block": []string{
						c.loginBlockKey,
					},
				},
				"urls": map[string]interface{}{
					"self": proxyUrl,
				},
			},
			"auth": map[string]interface{}{
				"urls": map[string]interface{}{
					"self": proxyUrl,
				},
			},
		},
		"dsn": fmt.Sprintf("host=localhost port=5433 user=%s password=%s dbname=%s sslmode=disable", c.coreDbUsername, c.coreDbPassword, c.coreDbName),
		"hydra": map[string]interface{}{
			"admin": map[string]interface{}{
				"host":      "localhost:8443",
				"base_path": "hydra-admin/",
				"schemes":   []string{"https"},
			},
			"public": map[string]interface{}{
				"host":      "localhost:8443",
				"base_path": "hydra/",
				"schemes":   []string{"https"},
			},
		},
	}

	yamlBytes, err := yaml.Marshal(coreConfig)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(CoreDir, "config.yaml"), yamlBytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}