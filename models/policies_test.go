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
		reader := strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1", "condition": {"attribute": "a1", "operator": "StringEquals", "values": ["v1"]}}`)

		policy, err := CreatePolicy(reader)

		So(err, ShouldBeNil)
		So(policy.Entities, ShouldResemble, []string{"e1", "e2"})
		So(policy.Role, ShouldResemble, "r1")
		So(policy.Condition, ShouldResemble, Condition{
			Attribute: "a1", Values: []string{"v1"}, Operator: OperatorStringEquals},
		)
	})
}

func TestCreateNewPolicyWithNoBody(t *testing.T) {
	Convey("When a policy message has no body, an error is returned", t, func() {
		policy, err := CreatePolicy(reader{})

		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, ErrorReadingBody)
		So(policy, ShouldBeNil)
	})
}

func TestCreateNewPolicyWithInvalidJson(t *testing.T) {
	Convey("When a policy message is missing entities or roles fields, an error is returned", t, func() {
		policy, err := CreatePolicy(strings.NewReader(`{"condition": {"attribute": "a1", "operator": "StringEquals", "values": ["v1"]}}`))
		So(err, ShouldBeNil)

		err = policy.ValidatePolicy()
		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, fmt.Errorf("missing mandatory fields: entities, role"))
	})

	Convey("When a policy message has empty entities fields, an error is returned", t, func() {
		policy, err := CreatePolicy(strings.NewReader(`{"entities": [], "role": "", "condition": {"attribute": "a1", "operator": "StringEquals", "values": ["v1"]}}`))
		So(err, ShouldBeNil)

		err = policy.ValidatePolicy()
		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, fmt.Errorf("missing mandatory fields: entities, role"))
	})

	Convey("When a policy message has an invalid condition operator, an error is returned", t, func() {
		policy, err := CreatePolicy(strings.NewReader(`{"entities": ["e1", "e2"], "role": "r1", "condition": {"attribute": "a1", "operator": "And", "values": ["v1"]}}`))
		So(err, ShouldBeNil)

		err = policy.ValidatePolicy()
		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, fmt.Errorf("invalid field values: condition operator And"))
	})
}
