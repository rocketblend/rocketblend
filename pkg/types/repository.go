package types

type (
	Repository interface {
		PackageRepository
		InstallationRepository
	}
)
