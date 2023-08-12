package qb_exec_bucket

import (
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	BucketResourceDownloader
//----------------------------------------------------------------------------------------------------------------------

type BucketResourceDownloader struct {
	root       string
	remoteRoot string
	downloads  []*qb_utils.DownloaderAction
}

func NewBucketResourceDownloader(root string) (instance *BucketResourceDownloader) {
	instance = new(BucketResourceDownloader)
	instance.root = root
	instance.downloads = make([]*qb_utils.DownloaderAction, 0)

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketResourceDownloader) SetRemoteRoot(remotePath string) {
	if nil != instance {
		instance.remoteRoot = remotePath
	}
}

func (instance *BucketResourceDownloader) AddResource(remoteRawPath, localRelativePath string) {
	remotePath, localPath := instance.normalizePaths(remoteRawPath, localRelativePath)
	download := qb_utils.IO.NewDownloaderAction(
		"",
		remotePath,
		"",
		localPath,
	)
	instance.addDownload(download)
}

func (instance *BucketResourceDownloader) DownloadAll(force bool) ([]string, []error) {
	session := qb_utils.IO.NewDownloadSession(instance.downloads)
	return session.DownloadAll(force)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketResourceDownloader) addDownload(download *qb_utils.DownloaderAction) {
	if nil != instance && nil != download && len(download.Source) > 0 {
		instance.downloads = append(instance.downloads, download)
	}
}

func (instance *BucketResourceDownloader) normalizePaths(remoteRawPath, localRelativePath string) (remotePath, localPath string) {
	remotePath = remoteRawPath
	isRemoteAbsolute := strings.HasPrefix(remotePath, "http")
	remoteDir := ""

	// remote
	if !isRemoteAbsolute {
		remoteDir = qb_utils.Paths.Dir(remotePath)
		if len(instance.remoteRoot) > 0 {
			remotePath = qb_utils.Paths.Concat(instance.remoteRoot, remotePath)
		}
	}

	// locale
	localDir := qb_utils.Paths.Dir(localRelativePath)
	if len(localDir) == 0 || localDir == "." {
		if len(remoteDir) > 0 && remoteDir != "." {
			localRelativePath = qb_utils.Paths.Concat(
				qb_utils.Paths.NormalizePathForOS(remoteDir),
				localRelativePath,
			)
		}
	}
	localPath = qb_utils.Paths.Absolutize(localRelativePath, instance.root)
	_ = qb_utils.Paths.Mkdir(localPath)

	return
}
