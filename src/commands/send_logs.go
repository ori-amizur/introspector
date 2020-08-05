package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
)

const logsDir = "/var/log"

var logsTmpFilesDir = path.Join(logsDir, "upload")

func getJournalLogsWithTag(tag string, since string, outputFilePath string) error {
	log.Infof("Running journalctl with tag %s", tag)
	cmd := exec.Command("journalctl", "-D", "/var/log/journal/",
		fmt.Sprintf("TAG=%s", tag), "--since", since, "--all")

	var stdoutBuf bytes.Buffer
	// open the out file for writing
	outfile, err := os.Create(outputFilePath)
	if err != nil {
		log.WithError(err).Errorf("Failed to create output file %s", outputFilePath)
		return err
	}
	defer outfile.Close()
	cmd.Stdout = outfile
	cmd.Stderr = &stdoutBuf

	err = cmd.Run()
	if err != nil {
		log.WithError(err).Errorf("Failed to run journalctl command, output %s", stdoutBuf)
		return err
	}
	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Infof("Creating dir %s", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.WithError(err).Errorf("Failed to create directory %s", dir)
			return err
		}
	}
	return nil
}

func archiveFilesInFolder(inputPath string, outputFile string, removeFiles bool) error {
	log.Infof("Archiving %s and creating %s", inputPath, outputFile)
	args := []string{"-czvf", outputFile, inputPath}
	if removeFiles {
		args = append(args, "--remove-files")
	}
	cmd := exec.Command("tar", args...)
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Errorf("Failed to run to archive %s", inputPath)
		return err
	}
	return nil
}

func SendLogs(tags []string, since string) error {
	log.Infof("Start gathering journalctl logs with tags %s", tags)
	if err := createDirIfNotExist(logsTmpFilesDir); err != nil {
		return err
	}
	for _, tag := range tags {
		_ = getJournalLogsWithTag(tag, since, path.Join(logsTmpFilesDir, fmt.Sprintf("%s.logs", tag)))
	}
	return archiveFilesInFolder(logsTmpFilesDir, fmt.Sprintf("%s/logs.tar.gz", logsDir), false)
}
