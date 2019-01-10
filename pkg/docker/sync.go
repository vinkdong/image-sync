package docker

import (
	"github.com/vinkdong/gox/slice"
	"regexp"
)

type NamedSync struct {
	ApiVersion string `yaml:"apiVersion"`
	Sync       Sync   `yaml:"sync"`
}

type Sync struct {
	From  Docker      `yaml:"from"`
	To    Docker      `yaml:"to"`
	Names []string    `yaml:"names"`
	Rules []NameValue `yaml:"rules"`
}

type NameValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (s *Sync) Do() error {
	s.From.Login()
	s.To.Login()
	for _, name := range s.Names {
		if err := s.syncTags(name); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sync) syncTags(name string) error {
	fromImage, err := s.From.listTags(name)
	if err != nil {
		return err
	}

	toImage, err := s.To.listTags(name)
	if err != nil {
		return err
	}
	diffTags := slice.Difference(fromImage.Tags, toImage.Tags)
	for _, tag := range diffTags {
		if s.matchRules(tag) {
			if err := s.sync(name, tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Sync) matchRules(tagName string) bool {
	for _, rule := range s.Rules {
		reg := regexp.MustCompile(rule.Value)
		if !reg.MatchString(tagName) {
			return false
		}
	}
	return true
}

func (s *Sync) sync(name, tag string) error {
	if err := s.From.pullImage(name, tag); err != nil {
		return err
	}
	if err := s.From.pushImage(name, tag); err != nil {
		return err
	}
	return nil
}
