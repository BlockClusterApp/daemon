package tasks

import "github.com/BlockClusterApp/daemon/src/helpers"

func ValidateLicence() {
	licenceKey := helpers.GetLicenceKey()

	var bc helpers.BlockCluster
	bc.Licence.Key = licenceKey
	bc.FetchLicenceDetails()
}
