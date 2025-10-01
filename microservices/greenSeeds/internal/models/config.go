package models

type Config struct {
	Server struct {
		Address   string `yaml:"address"`
		PathToLog string `yaml:"path_to_log"`
	}

	Auth struct {
		DefaultPassword string `yaml:"default_password"`
	}

	Database struct {
		Host string `yaml:"host"`
		Port uint64 `yaml:"port"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"password"`
	}

	JWT struct {
		Secret    string `yaml:"client_secret"`
		ExpiresIn int    `yaml:"expires_in"`
	}
}
