package jsonapi

type TaskAttachIncluded struct {
	Parent       TaskResultRecords
	II           *IncludeInstructions
	Output       chan chan *Output
	ActualOutput *Output
	OutputType   OutputType
	Linkname     string
}

type OutputType int

const (
	OutputTypeResources OutputType = iota
	OutputTypeLinkages
)

func NewTaskAttachIncluded(parent TaskResultRecords, ii *IncludeInstructions, outputtype OutputType, linkname string) *TaskAttachIncluded {
	return &TaskAttachIncluded{
		Parent:     parent,
		II:         ii,
		Output:     make(chan chan *Output),
		OutputType: outputtype,
		Linkname:   linkname,
	}
}

func (t *TaskAttachIncluded) Work(r *Request) {
	parent_result := t.Parent.GetResult()
	r.API.Logger.Debugf("PARENT: %#v\n", parent_result)
	res := NewOutput()

	output_primary := []*Record{}
	output_included := []*Record{}
	var output_relationship *ORelationship = nil

	queue := parent_result.Records
	primary_data_count := len(queue)
	for {
		if len(queue) == 0 {
			break
		}
		var next *Record
		next, queue = queue[0], queue[1:] // queue pop

		r.API.Logger.Infof("MAIN LOOP HANDLING: %#v\n", next)
		//relationships := next.GetRelationships();
		if primary_data_count > 0 {
			r.API.Logger.Infof("MAIN LOOP INCLUDING PRIMARY: %#v\n", next)
			primary_data_count--
			output_primary = append(output_primary, next)
			rel := next.GetRelationships().Relationships.GetRelationshipByName(t.Linkname)
			if rel != nil {
				if output_relationship == nil {
					output_relationship = rel
				} else {
					output_relationship.Data = append(output_relationship.Data, rel.Data...)
				}
			}
		} else if next.ShouldInclude {
			r.API.Logger.Infof("MAIN LOOP INCLUDING: %#v\n", next)
			output_included = append(output_included, next)
		}
		rels := next.GetRelationships()
		queue = append(queue, rels.Included...)
		next.Relationships = rels.Relationships
	}

	if t.OutputType == OutputTypeResources {
		res.Data = &ORecords{
			IsSingle: parent_result.IsSingle,
			Records:  output_primary,
		}
	} else {
		res.Data = output_relationship
	}
	res.Included = output_included
	t.ActualOutput = res
}

func (w *TaskAttachIncluded) ResponseWorker(has_paniced bool) {
	go func() {
		for req := range w.Output {
			req <- w.ActualOutput
		}
	}()
}

func (w *TaskAttachIncluded) Cleanup(r *Request) {
	r.API.Logger.Debugf("TaskAttachIncluded.Cleanup\n")
	close(w.Output)
}

func (w *TaskAttachIncluded) GetResult() *Output {
	r := make(chan *Output)
	defer close(r)
	w.Output <- r
	return <-r
}
