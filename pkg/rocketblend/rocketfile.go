package rocketblend

import (
	"sort"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketfile"
)

func (d *Driver) InstallDependencies(dir string, ref *reference.Reference, force bool) error {
	rkt, err := rocketfile.Load(dir)
	if err != nil {
		return err
	}

	if ref != nil {
		pack, err := d.DescribePackByReference(*ref)
		if err != nil {
			return err
		}

		if pack.Addon != nil {
			rkt.Addons = RemoveDuplicateStr(append(rkt.Addons, ref.String()))
		}

		if pack.Build != nil {
			rkt.Build = ref.String()
		}
	}

	deps := append(rkt.Addons, rkt.Build)

	for _, dep := range deps {
		ref, err := reference.Parse(dep)
		if err != nil {
			return err
		}

		err = d.InstallPackByReference(ref, force)
		if err != nil {
			return err
		}
	}

	// Save the rocketfile. This will only happen if a pack was added.
	if ref != nil {
		err := rocketfile.Save(dir, rkt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Driver) UninstallDependencies(dir string, ref reference.Reference) error {
	rkt, err := rocketfile.Load(dir)
	if err != nil {
		return err
	}

	if rkt.Build == ref.String() {
		rkt.Build = ""
	}

	rkt.Addons = RemoveStr(rkt.Addons, ref.String())

	err = rocketfile.Save(dir, rkt)
	if err != nil {
		return err
	}

	return nil
}

func RemoveStr(strs []string, param string) []string {
	sort.Strings(strs)
	for i := len(strs) - 1; i > 0; i-- {
		if strs[i] == strs[i-1] && strs[i] == param {
			strs = append(strs[:i], strs[i+1:]...)
		}
	}

	return strs
}

func RemoveDuplicateStr(strs []string) []string {
	sort.Strings(strs)
	for i := len(strs) - 1; i > 0; i-- {
		if strs[i] == strs[i-1] {
			strs = append(strs[:i], strs[i+1:]...)
		}
	}

	return strs
}
