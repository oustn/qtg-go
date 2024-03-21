package download

import (
	mapset "github.com/deckarep/golang-set"
	teacommon "github.com/oustn/qtg/internal/ui/common"
	"sync"
)

type status int16

const (
	Update      status = -1
	Finished    status = 0
	Errored     status = 1
	Downloading status = 2
)

type Status struct {
	Id string
	S  status
}

type Progress struct {
	Pending     mapset.Set
	Downloading mapset.Set
	Finished    mapset.Set
	Error       mapset.Set
}

func (p *Progress) Copy() Progress {
	return Progress{
		Pending:     p.Pending,
		Downloading: p.Downloading,
		Finished:    p.Finished,
		Error:       p.Error,
	}
}

func (p *Progress) Total() int {
	if p.Pending == nil {
		return 0
	}
	return p.Pending.Cardinality() + p.Downloading.Cardinality() + p.Finished.Cardinality() + p.Error.Cardinality()
}

type lockedProgress struct {
	sync.RWMutex
	Progress
}

func (l *lockedProgress) Update(s Status) {
	switch s.S {
	case Downloading:
		l.SetDownloading(s.Id)
	case Finished:
		l.SetFinished(s.Id)
	case Errored:
		l.SetError(s.Id)
	}
}

func (l *lockedProgress) SetError(id string) {
	l.Lock()
	defer l.Unlock()
	l.Downloading.Remove(id)
	l.Error.Add(id)
}

func (l *lockedProgress) SetFinished(id string) {
	l.Lock()
	defer l.Unlock()
	l.Downloading.Remove(id)
	l.Finished.Add(id)
}

func (l *lockedProgress) SetDownloading(id string) {
	l.Lock()
	defer l.Unlock()
	l.Pending.Remove(id)
	l.Downloading.Add(id)
}

func (l *lockedProgress) add(id string) {
	if l.Pending == nil {
		l.Pending = mapset.NewSet()
		l.Downloading = mapset.NewSet()
		l.Finished = mapset.NewSet()
		l.Error = mapset.NewSet()
	}
	l.Pending.Add(id)
}

type ProgramProgress struct {
	Progress
	Programs map[string]teacommon.Program
}

type Channel struct {
	teacommon.Channel
	Programs ProgramProgress
}

type ChannelProgress struct {
	Progress
	Channels map[string]Channel
}

type lockedProgramProcess struct {
	lockedProgress
	Programs map[string]teacommon.Program
}

func (p *lockedProgramProcess) AddProgram(program teacommon.Program) {
	p.Lock()
	defer p.Unlock()
	if p.Programs == nil {
		p.Programs = make(map[string]teacommon.Program)
	}
	p.Programs[program.StringId()] = program
	p.add(program.StringId())
}

type lockedChannelProcess struct {
	lockedProgress
	Channels map[string]Channel
}

func (p *lockedChannelProcess) AddChannel(channel teacommon.Channel) {
	p.Lock()
	defer p.Unlock()
	if p.Channels == nil {
		p.Channels = make(map[string]Channel)
	}
	p.Channels[channel.Id] = Channel{Channel: channel}
	p.add(channel.Id)
}

func (p *lockedChannelProcess) HasChannel(id string) bool {
	_, ok := p.Channels[id]
	return ok
}

func (p *lockedChannelProcess) UpdatePrograms(channelId string, s ProgramProgress) {
	p.Lock()
	defer p.Unlock()
	if p.Channels == nil {
		return
	}
	channel, ok := p.Channels[channelId]
	if !ok {
		return
	}
	channel.Programs = ProgramProgress{
		Progress: s.Progress,
		Programs: s.Programs,
	}
	p.Channels[channelId] = channel
}
