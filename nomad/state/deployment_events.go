package state

import (
	"fmt"

	"github.com/hashicorp/nomad/nomad/stream"
	"github.com/hashicorp/nomad/nomad/structs"
)

func DeploymentEventFromChanges(msgType structs.MessageType, tx ReadTxn, changes Changes) ([]stream.Event, error) {
	var events []stream.Event

	var deployment *structs.Deployment
	var job *structs.Job
	var eval *structs.Evaluation

	var eventType string
	switch msgType {
	case structs.DeploymentStatusUpdateRequestType:
		eventType = TypeDeploymentUpdate
	case structs.DeploymentPromoteRequestType:
		eventType = TypeDeploymentPromotion
	case structs.DeploymentAllocHealthRequestType:
		eventType = TypeDeploymentAllocHealth
	}

	for _, change := range changes.Changes {
		switch change.Table {
		case "deployment":
			after, ok := change.After.(*structs.Deployment)
			if !ok {
				return nil, fmt.Errorf("transaction change was not a Deployment")
			}

			deployment = after
			event := stream.Event{
				Topic:      TopicDeployment,
				Type:       eventType,
				Index:      changes.Index,
				Key:        deployment.ID,
				FilterKeys: []string{deployment.JobID},
				Payload: &DeploymentEvent{
					Deployment: deployment,
				},
			}

			events = append(events, event)
		case "jobs":
			after, ok := change.After.(*structs.Job)
			if !ok {
				return nil, fmt.Errorf("transaction change was not a Job")
			}

			job = after
			event := stream.Event{
				Topic:      TopicJob,
				Type:       eventType,
				Index:      changes.Index,
				Key:        job.ID,
				FilterKeys: []string{},
				Payload: &JobEvent{
					Job: job,
				},
			}

			events = append(events, event)
		case "evals":
			after, ok := change.After.(*structs.Evaluation)
			if !ok {
				return nil, fmt.Errorf("transaction change was not an Evaluation")
			}

			eval = after

			event := stream.Event{
				Topic:      TopicEval,
				Type:       eventType,
				Index:      changes.Index,
				Key:        eval.ID,
				FilterKeys: []string{eval.DeploymentID, eval.JobID},
				Payload: &EvalEvent{
					Eval: eval,
				},
			}

			events = append(events, event)
		}
	}

	return events, nil
}
