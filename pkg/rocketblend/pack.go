package rocketblend

import (
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketpack"
)

func (d *Driver) DescribePackByReference(ref reference.Reference) (*rocketpack.RocketPack, error) {
	pack, err := d.pack.DescribeByReference(ref)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (d *Driver) InstallPackByReference(ref reference.Reference, force bool) error {
	err := d.pack.InstallByReference(ref, force)
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) UninstallPackByReference(ref reference.Reference) error {
	err := d.pack.UninstallByReference(ref)
	if err != nil {
		return err
	}

	return nil
}
