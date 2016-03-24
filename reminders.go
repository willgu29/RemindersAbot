package reminders

import(
	'log'

)

var p *dt.Plugin

type reminder struct {
	UserID uint64
	Time time.Time
	Content string
}


func init() {
	trigger := &nlp.StructuredInput{
		Commands: []string{"remind", "tell"},
		Objects: []string{"me"},

	}

	fns := &dt.PluginFns{Run: Run, FollowUp: FollowUp}
	var err error
	p, err = plugin.New("github.com/willgu29/RemindersAbot", trigger, fns)
	if err != nil {
		log.Fatal(err)
	}

	p.Vocab = dt.NewVocab(
		dt.VocabHandler{
			Fn: kwSetReminder,
			Trigger: &nlp.StructuredInput{
				Commands: []string{"remind", "tell"},
				Objects: []string{"me"},
			},

		},
	)
}

func kwSetReminder(in *dt.Msg) (resp string) {
	var toFound, commandFound bool

	var cmd, time string
	prepositions := convertMapToSlice(language.Prepositions)

	for _, token := range in.Tokens {
		if (token == "to") {
			toFound = true
			continue
		}
		if (!toFound) {
			continue
		}
		if (language.Contains(prepositions, token)) {
			commandFound = true
			continue
		}
		cmd += token + " "
		if (!commandFound) {
			continue
		}
		time += token + " "


	}
	atTime := convertToDate(time)
	reminder, err := makeReminder(in,cmd[:len(cmd)-1], atTime, recurring)
	
	sm := makeStateMachine(in)
	//TODO
	sm.MakeReminder(reminder, false)

	return "Okay, I'll remind you to " + cmd + "at " + time[:len(time)-1]
}

func makeReminder(in *dt.Msg, s string,t time.Time) ([]byte, error) {
	r := reminder{
		UserID: in.User.ID,
		Time: t,
		Content: s,
	}

	b, err := json.Marshal(r)
	if (err != nil) {
		return nil, err
	}
	return b, nil
}

func convertToDate(s string) (time.Time, error){
	ts, err := time.Parse(s)
	if (err != nil) {
		return nil, err
	}
	if (len(ts) == 0) {
		return nil, nil
	}
	return ts[0], nil
}

func convertMapToSlice(m map[string]struct{}) []string {
	slice := []string{}
	for k, _ := range m {
		slice = append(slice, k)
	}
	return slice
}

func makeStateMachine(in *dt.Msg) (*dt.StateMachine) {
	sm := dt.NewStateMachine(p)
	sm.SetStates([]dt.State{})
	sm.LoadState(in)
	sm.SetMemory(in, "reminder")
	return sm
}
