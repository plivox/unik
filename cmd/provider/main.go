package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/unik/pkg/compilers"
	"github.com/emc-advanced-dev/unik/pkg/providers/aws"
	"github.com/emc-advanced-dev/unik/pkg/config"
	"github.com/layer-x/layerx-commons/lxlog"
)

func main() {
	os.Setenv("TMPDIR", "/Users/pivotal/tmp/uniktest")
	log.SetLevel(log.DebugLevel)

	r := compilers.RunmpCompiler{
		DockerImage: "rumpcompiler-go-xen",
		CreateImage: compilers.CreateImageAws,
	}
	f, err := os.Open("a.tar")
	if err != nil {
		log.Fatal(err)
	}
	rawimg, err := r.CompileRawImage(f, "", []string{})
	if err != nil {
		log.Fatal(err)
	}
	c := config.Aws{
		Name: "aws-provider",
		AwsAccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAcessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region: os.Getenv("AWS_REGION"),
		Zone: os.Getenv("AWS_AVAILABILITY_ZONE"),
	}
	p := aws.NewAwsProvier(c)
	defer func() {
		err = p.Save()
		if err != nil {
			log.WithError(err).Error("failed to save")
		}
	}()


	logger := lxlog.New("scott")
	logger.SetLogLevel(lxlog.DebugLevel)

	img, err := p.Stage(logger, "test-scott", rawimg, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(img)
	fmt.Println()

	env := make(map[string]string)
	env["FOO"] = "BAR"

	instance, err := p.RunInstance(logger, "test-scott-instance-1", img.Id, nil, env)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(instance)
	fmt.Println()

	images, err := p.ListImages(logger)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(images)
	fmt.Println()

	instances, err := p.ListInstances(logger)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(instances)
	fmt.Println()

	for _, instance := range instances {
		err = p.DeleteInstance(logger, instance.Id)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, image := range images {
		err = p.DeleteImage(logger, image.Id, false)
		if err != nil {
			log.Fatal(err)
		}
	}

}