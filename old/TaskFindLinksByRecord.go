package jsonapi

type TaskFindLinksByRecordResult struct {
	Relationships *ORelationships
	Included      []*Record
}

type TaskFindLinksByRecord struct {
	Record *Record
	II     *IncludeInstructions
	Output chan chan *TaskFindLinksByRecordResult
	Result *TaskFindLinksByRecordResult
}

func NewTaskFindLinksByRecord(r *Record, ii *IncludeInstructions) *TaskFindLinksByRecord {
	return &TaskFindLinksByRecord{
		II:     ii,
		Record: r,
		Output: make(chan chan *TaskFindLinksByRecordResult),
	}
}

func (t *TaskFindLinksByRecord) Work(r *Request) {
	result := &TaskFindLinksByRecordResult{
		Relationships: &ORelationships{
			Relationships: []*ORelationship{},
		},
		Included: []*Record{},
	}
	for linkname, relationship := range r.API.GetRelationshipsByResource(t.Record.Type) {
		shouldFetch := t.II.ShouldFetch(linkname)
		r.API.Logger.Debugf("SHOULDFETCH: %#v %#v %v %v\n", linkname, t.II, shouldFetch, t.II.ShouldInclude(linkname))
		or, included := relationship.Resolve(r, t.Record, shouldFetch, t.II)
		or.RelatedBase = r.GetBaseURL() + t.Record.Type + "/" + t.Record.Id
		or.RelationshipName = linkname
		result.Relationships.Relationships = append(result.Relationships.Relationships, or)
		if t.II.ShouldInclude(linkname) {
			for _, record := range included {
				record.PrepareRelationships(r, t.II.GetChild(linkname))
				record.ShouldInclude = true
			}
		}
		result.Included = append(result.Included, included...)
	}

	t.Result = result
}

func (t *TaskFindLinksByRecord) ResponseWorker(has_paniced bool) {
	go func() {
		for r := range t.Output {
			r <- t.Result
		}
	}()
}

func (t *TaskFindLinksByRecord) Cleanup(r *Request) {
	close(t.Output)
}

func (t *TaskFindLinksByRecord) GetResult() *TaskFindLinksByRecordResult {
	r := make(chan *TaskFindLinksByRecordResult)
	defer close(r)
	t.Output <- r
	res := <-r
	return res
}
