package models

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type reader struct {
}

func (f reader) Read(bytes []byte) (int, error) {
	return 0, fmt.Errorf("reader failed")
}

func TestCreateNewPolicyWithValidJson(t *testing.T) {
	Convey("When a policy has a valid json body, a new policy is returned", t, func() {
		reader := strings.NewReader(`{"entities": ["e1", "e2"], "roles": ["r1", "r2"], "conditions": [{"attributes": ["a1"], "operator": "and", "values": ["v1"]}]}`)

		policy, err := CreateNewPolicy(reader)

		So(err, ShouldBeNil)
		So(policy.Entities, ShouldResemble, []string{"e1", "e2"})
		So(policy.Roles, ShouldResemble, []string{"r1", "r2"})
		So(policy.Conditions, ShouldResemble, []Condition{
			{Attributes: []string{"a1"}, Values: []string{"v1"}, Operator: "and"}},
		)
	})
}

func TestCreateNewPolicyWithNoBody(t *testing.T) {
	Convey("When a policy message has no body, an error is returned", t, func() {

		policy, err := CreateNewPolicy(reader{})

		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, ErrorReadingBody)
		So(policy, ShouldBeNil)
	})
}

func TestCreateNewPolicyWithInvalidJson(t *testing.T) {
	Convey("When a policy message is missing entities or roles fields, an error is returned", t, func() {
		policy, err := CreateNewPolicy(strings.NewReader(`{"conditions": [{"attributes": ["a1"], "operator": "and", "values": ["v1"]}]}`))
		So(err, ShouldBeNil)

		err = policy.ValidateNewPolicy()
		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, fmt.Errorf("missing mandatory fields: entities, roles"))
	})

	Convey("When a policy message has empty entities fields, an error is returned", t, func() {
		policy, err := CreateNewPolicy(strings.NewReader(`{"entities": [], "roles": [], "conditions": [{"attributes": ["a1"], "operator": "and", "values": ["v1"]}]}`))
		So(err, ShouldBeNil)

		err = policy.ValidateNewPolicy()
		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, fmt.Errorf("missing mandatory fields: entities, roles"))
	})
}
