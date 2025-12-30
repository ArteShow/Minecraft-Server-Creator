package eulaacceptor

import "os"

func WriteEULA(serverDir string) error {
	return os.WriteFile(
		serverDir+"/eula.txt",
		[]byte("eula=true\n"),
		0644,
	)
}
