package integration

import (
	"bytes"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/samalba/dockerclient"

	. "gopkg.in/check.v1"
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

func RandStr(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}
	return string(b)
}

type RunSuite struct {
	command  string
	projects []string
}

var _ = Suite(&RunSuite{
	command: "/home/darren/src/libcompose/test",
})

func (s *RunSuite) CreateProjectFromText(c *C, input string) string {
	return s.ProjectFromText(c, "create", input)
}

func (s *RunSuite) ProjectFromText(c *C, command, input string) string {
	projectName := "test-project-" + RandStr(7)
	return s.FromText(c, projectName, command, input)
}

func (s *RunSuite) FromText(c *C, projectName, command, input string) string {
	args := []string{"--verbose", "-p", projectName, "-f", "-", command}
	if command == "up" {
		args = append(args, "-d")
	} else if command == "down" {
		args = append(args, "--timeout", "0")
	}

	cmd := exec.Command(s.command, args...)
	cmd.Stdin = bytes.NewBufferString(strings.Replace(input, "\t", "  ", -1))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		logrus.Errorf("Failed to run %s %v: %v\n with input:\n%s", s.command, err, args, input)
	}

	c.Assert(err, IsNil)

	return projectName
}

func GetClient(c *C) dockerclient.Client {
	context := docker.Context{}
	err := context.CreateClient()

	c.Assert(err, IsNil)

	return context.Client
}

func (s *RunSuite) GetContainerByName(c *C, name string) *dockerclient.ContainerInfo {
	client := GetClient(c)
	container, err := docker.GetContainerByName(client, name)

	c.Assert(err, IsNil)

	if container == nil {
		return nil
	}

	info, err := client.InspectContainer(container.Id)

	c.Assert(err, IsNil)

	return info
}
