package genusers

import (
	"github.com/rskvp/qb-core/qb_utils"
)

/**
Download templates from GitHub or Bitbucket
*/

const (
	repoRoot = "https://bitbucket.org/digi-sense/qb-core/raw/master/qb_generators/genusers/data/"
)

// ---------------------------------------------------------------------------------------------------------------------
// 	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func DownloadTemplates(dirTarget string) ([]string, []error) {
	// download
	session := qb_utils.IO.NewDownloadSession(templateActions(dirTarget))
	return session.DownloadAll(false)
}

// ---------------------------------------------------------------------------------------------------------------------
// 	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func templateActions(dirTarget string) []*qb_utils.DownloaderAction {
	response := make([]*qb_utils.DownloaderAction, 0)

	// data
	response = append(response, qb_utils.IO.NewDownloaderAction("", repoRoot+"names.csv", "", dirTarget))
	response = append(response, qb_utils.IO.NewDownloaderAction("", repoRoot+"surnames.csv", "", dirTarget))
	response = append(response, qb_utils.IO.NewDownloaderAction("", repoRoot+"country_codes.csv", "", dirTarget))

	return response
}
