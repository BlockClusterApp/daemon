package tasks

import "github.com/BlockClusterApp/daemon/src/helpers"

func ValidateLicence() {
	helpers.UpdateLicence()

	bc := helpers.GetBlockclusterInstance()
	bc.Licence.Key = helpers.GetLicence().Key
	bc.FetchLicenceDetails()
}
