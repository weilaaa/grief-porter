package main

import "fmt"

type source struct {
	// Addr the image source address
	Addr string `json:"addr"`

	// Remark just remark in config
	Remark string `json:"remark"`

	// NewTag a new tag when push to target registry
	NewTag string `json:"new_tag"`

	// Platform will pull image like 'docker pull --platform=xxx'
	Platform string `json:"platform"`

	// Skip will count it out if true
	Skip bool `json:"skip"`
}

func (s *source) pull() cmdInstruction {
	return makeCmdInstruction("docker pull %s", s.Addr)
}

func (s *source) pullByPlatform() cmdInstruction {
	return makeCmdInstruction("docker pull --platform=%s %s", s.Platform, s.Addr)
}

func (s *source) tag() cmdInstruction {
	return makeCmdInstruction("docker tag %s %s", s.Addr, s.NewTag)
}

func (s *source) push() cmdInstruction {
	return makeCmdInstruction("docker push %s", s.NewTag)
}

func (s *source) inspect() cmdInstruction {
	return makeCmdInstruction("docker manifest inspect --verbose %s", s.Addr)
}

// singleMove just do one image lift
func singleMove(s *source) error {
	if s.Skip {
		return nil
	}

	var err error

	if len(s.Platform) > 0 {
		err = s.pullByPlatform().doExec()
	} else {
		err = s.pull().doExec()
	}
	if err != nil {
		return fmt.Errorf("pull image %v failed: %v", s.Addr, err)
	}

	err = s.tag().doExec()
	if err != nil {
		return fmt.Errorf("tag image %v failed: %v", s.Addr, err)
	}

	err = s.push().doExec()
	if err != nil {
		return fmt.Errorf("push image %v to %v failed: %v", s.Addr, s.NewTag, err)
	}

	return nil
}
