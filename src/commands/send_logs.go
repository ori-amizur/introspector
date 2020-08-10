package commands

import (
	"bytes"
	"fmt"
	"github.com/ori-amizur/introspector/src/config"
	"os"
	"os/exec"
	"path"

	"github.com/go-openapi/strfmt"
	"github.com/openshift/assisted-service/client/installer"
	"github.com/ori-amizur/introspector/src/session"
	"github.com/ori-amizur/introspector/src/util"

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
		log.WithError(err).Errorf("Failed to run journalctl command, output %s", stdoutBuf.String())
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

	_, err, execCode := util.Execute("tar", args...)

	if execCode != 0 {
		log.Errorf("Failed to run to archive %s. output", inputPath, execCode)
		return fmt.Errorf(err)
	}
	return nil
}

func uploadLogs(filepath string, clusterID strfmt.UUID, hostId strfmt.UUID, removeAfterUpload bool, inventoryUrl string) error {
	uploadFile, err := os.Open(filepath)
	defer uploadFile.Close()
	if err != nil {
		log.WithError(err).Errorf("Failed to open file %s for upload", uploadFile)
		return err
	}

	fmt.Println("PARAMS ", hostId, clusterID, inventoryUrl)
	invSession := session.New(inventoryUrl)
	params := installer.UploadHostLogsParams{
		Upfile:    uploadFile,
		ClusterID: clusterID,
		HostID:    hostId,
	}

	_, err = invSession.Client().Installer.UploadHostLogs(invSession.Context(), &params)

	if err != nil {
		log.WithError(err).Errorf("Failed to upload file %s to assisted-service", filepath)
		return err
	}

	if removeAfterUpload {
		_ = os.Remove(filepath)
	}
	return nil
}

func SendLogs() error {
	tags := config.LogsSenderConfig.Tags.ToStringList()

	log.Infof("Start gathering journalctl logs with tags %s", tags)
	archivePath := fmt.Sprintf("%s/logs.tar.gz", logsDir)
	if err := createDirIfNotExist(logsTmpFilesDir); err != nil {
		return err
	}
	for _, tag := range tags {
		_ = getJournalLogsWithTag(tag, config.LogsSenderConfig.Since, path.Join(logsTmpFilesDir, fmt.Sprintf("%s.logs", tag)))
	}

	if err := archiveFilesInFolder(logsTmpFilesDir, archivePath, config.LogsSenderConfig.CleanWhenDone); err != nil {
		return err
	}
    // strfmt.UUID(config.LogsSenderConfig.ClusterID)
	return uploadLogs(archivePath, strfmt.UUID(config.LogsSenderConfig.ClusterID),
		strfmt.UUID(config.LogsSenderConfig.HostID), config.LogsSenderConfig.CleanWhenDone, config.LogsSenderConfig.TargetURL)
}
