package config

var jwtKey = []byte("SOMESCERET")

func GetJwtKey() []byte {
	return jwtKey
}
