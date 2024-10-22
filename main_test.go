package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-permissions-api/features/steps"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

const MongoVersion = "4.4.8"
const DatabaseName = "testing"

var componentFlag = flag.Bool("component", false, "perform component tests")

type ComponentTest struct {
	MongoFeature *componenttest.MongoFeature
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	authorizationFeature := componenttest.NewAuthorizationFeature()
	component, err := steps.NewPermissionsComponent(f.MongoFeature.Server.URI())
	if err != nil {
		panic(err)
	}

	ctx.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
		if err := component.Close(); err != nil {
			return nil, err
		}
		authorizationFeature.Close()
		return ctx, nil
	})

	component.RegisterSteps(ctx)
	component.APIFeature.RegisterSteps(ctx)
	authorizationFeature.RegisterSteps(ctx)
}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		f.MongoFeature = componenttest.NewMongoFeature(componenttest.MongoOptions{MongoVersion: MongoVersion, DatabaseName: DatabaseName})
	})
	ctx.AfterSuite(func() {
		f.MongoFeature.Close()
	})
}

func TestComponent(t *testing.T) {
	flag.Parse()
	*componentFlag = false // put this line in if you want to "debug test" this function in vscode IDE
	if *componentFlag {
		status := 0

		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
		}

		f := &ComponentTest{}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		fmt.Println("=================================")
		fmt.Printf("Component test coverage: %.2f%%\n", testing.Coverage()*100)
		fmt.Println("=================================")

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
