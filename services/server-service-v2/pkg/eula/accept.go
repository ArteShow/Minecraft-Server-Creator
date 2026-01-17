package eula

func Accept() ([]byte, error) {
	data := []byte("eula=true\n")
	return data, nil
}
