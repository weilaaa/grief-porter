package main

import "fmt"

type brick struct {
	/*
		1. len(sources) == 1
		   sources will be tagged which destination defined and push then
		2. len(source) > 1
		   sources will be combined into one manifest and push then
	*/

	// sources can be multi of an image
	Sources []*source `json:"sources"`

	// manifest generic private register like: foo.register/bar-v1.0
	// used to combined multi digest to one manifest
	Manifest string `json:"manifest"`
}

type source struct {
	Addr   string `json:"addr"`
	Remark string `json:"remark"`
	NewTag string `json:"new_tag"`
	Skip   bool   `json:"skip"`
}

func (s *source) pull() error {
	cmd := fmt.Sprintf("docker pull %s", s.Addr)
	return doExec(cmd, "")
}

func (s *source) tag() error {
	cmd := fmt.Sprintf("docker tag %s %s", s.Addr, s.NewTag)
	return doExec(cmd, "")
}

func (s *source) push() error {
	cmd := fmt.Sprintf("docker push %s", s.NewTag)
	return doExec(cmd, "")
}

func (b *brick) createManifest() error {
	digests := ""

	for _, s := range b.Sources {
		if s.Skip {
			continue
		}

		digests += fmt.Sprintf(" %s", s.NewTag)
	}

	if len(digests) < 1 {
		return fmt.Errorf("digest less than 1, can not create manifest %s", b.Manifest)
	}

	cmd := fmt.Sprintf("docker manifest create %s%s", b.Manifest, digests)
	return doExec(cmd, "")
}

// pushManifest would failed if remote register already had such manifest
func (b *brick) pushManifest() error {
	cmd := fmt.Sprintf("docker manifest push %s", b.Manifest)
	return doExec(cmd, "")
}

func (b *brick) moving() error {
	if len(b.Sources) == 1 {
		return singleMove(b.Sources[0])
	}

	return multiMove(b)
}

// singleMove just do one image lift
func singleMove(s *source) error {
	if s.Skip {
		return nil
	}

	err := s.pull()
	if err != nil {
		return fmt.Errorf("pull image %v failed: %v", s.Addr, err)
	}

	err = s.tag()
	if err != nil {
		return fmt.Errorf("tag image %v failed: %v", s.Addr, err)
	}

	err = s.push()
	if err != nil {
		return fmt.Errorf("push image %v to %v failed: %v", s.Addr, s.NewTag, err)
	}

	return nil
}

// multiMove can move multi images at same time
// if manifest specified, would combine multi images
// into one manifest
func multiMove(b *brick) error {
	var err error

	for _, s := range b.Sources {
		err = singleMove(s)
		if err != nil {
			return err
		}
	}

	if len(b.Manifest) > 0 {
		err = b.createManifest()
		if err != nil {
			return err
		}

		err = b.pushManifest()
		if err != nil {
			return err
		}
	}

	return nil
}
