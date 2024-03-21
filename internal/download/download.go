package download

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/oustn/qtg/internal/api"
	"github.com/oustn/qtg/internal/meta"
	teacommon "github.com/oustn/qtg/internal/ui/common"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

type downloadPayload struct {
	channel teacommon.Channel
	program teacommon.Program
}

type Downloader struct {
	api               *api.QingTingApi
	channelConcurrent int
	programConcurrent int
	task              chan teacommon.Channel // 下载任务 chan
	taskResult        chan Status            // 下载任务结果 chan
	progress          lockedChannelProcess   // 下载进度
	wg                sync.WaitGroup
	client            *grab.Client
	Callback          func(progress ChannelProgress)
}

var dl *Downloader

func NewDownloader(
	qtApi *api.QingTingApi,
	channelConcurrent int,
	programConcurrent int,
) *Downloader {
	var s sync.Once
	s.Do(func() {
		dl = &Downloader{
			api:               qtApi,
			channelConcurrent: channelConcurrent,
			programConcurrent: programConcurrent,
			task:              make(chan teacommon.Channel, channelConcurrent),
			taskResult:        make(chan Status),
			client:            grab.NewClient(),
		}
		go dl.start()
	})
	return dl
}

func (d *Downloader) start() {
	for i := 0; i < d.channelConcurrent; i++ {
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for channel := range d.task {
				d.download(channel)
			}
		}()
	}
	for task := range d.taskResult {
		d.progress.Update(task)
		if d.Callback != nil {
			d.Callback(ChannelProgress{
				Progress: Progress{
					Pending:     d.progress.Pending,
					Downloading: d.progress.Downloading,
					Finished:    d.progress.Finished,
					Error:       d.progress.Error,
				},
				Channels: d.progress.Channels,
			})
		}
	}
	d.wg.Wait()
}

func (d *Downloader) DownloadChannel(channel teacommon.Channel) {
	if d.progress.HasChannel(channel.Id) {
		return
	}
	// println("创建 channel 下载任务" + channel.Name)
	d.progress.AddChannel(channel)
	d.task <- channel
	d.taskResult <- Status{
		Id: channel.Id,
		S:  Update,
	}
}

func (d *Downloader) download(channel teacommon.Channel) {
	// println("开始下载 channel: " + channel.Name)
	d.taskResult <- Status{
		Id: channel.Id,
		S:  Downloading,
	}
	s := lockedProgramProcess{}
	channel = d.api.GetChannelInfo(channel)
	programCh := make(chan downloadPayload, channel.Count)
	resultCh := make(chan Status, channel.Count*2)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		d.downloadPrograms(programCh, resultCh)
	}()

	go func() {
		programs := d.fetchPrograms(channel)
		for _, program := range programs {
			s.AddProgram(program)

			programCh <- downloadPayload{
				channel: channel,
				program: program,
			}
		}
		close(programCh)
		wg.Wait()
		close(resultCh)
	}()

	for res := range resultCh {
		s.Update(res)
		d.progress.UpdatePrograms(channel.Id, ProgramProgress{
			Progress: s.Progress,
			Programs: s.Programs,
		})
		d.taskResult <- Status{
			Id: channel.Id,
			S:  Update,
		}
	}
	d.taskResult <- Status{
		Id: channel.Id,
		S:  Finished,
	}
	// println("完成 channel 下载: " + channel.Name)
}

func (d *Downloader) downloadPrograms(ch <-chan downloadPayload, resultCh chan<- Status) {
	var wg sync.WaitGroup

	for i := 0; i < d.programConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for info := range ch {
				d.downloadProgram(info, resultCh)
			}
		}()
	}

	wg.Wait()
}

func (d *Downloader) fetchPrograms(channel teacommon.Channel) []teacommon.Program {
	// println("开始获取 channel " + channel.Name + " 的 programs")

	var programs []teacommon.Program
	if channel.Count > 0 {
		programs = make([]teacommon.Program, channel.Count)
	}

	pages := int(math.Ceil(float64(channel.Count) / (100.0)))
	index := 0

	strLen := strings.Count(strconv.Itoa(channel.Count), "") - 1

	for page := range pages {
		p := d.api.FetchPrograms(channel, page)
		for _, program := range p {
			program.Index = index
			program.NamePrefix = leftPad2Len(strconv.Itoa(index+1), "0", strLen)
			if channel.Count > 0 {
				programs[index] = program
			} else {
				programs = append(programs, program)
			}
			index++
		}
	}

	// println("完成获取 channel " + channel.Name + " 的 programs")
	return programs
}

func (d *Downloader) downloadProgram(payload downloadPayload, resultCh chan<- Status) {
	// println("开始下载 Program " + payload.program.Name)
	resultCh <- Status{
		Id: payload.program.StringId(),
		S:  Downloading,
	}
	handleErr := func(err error) {
		resultCh <- Status{
			Id: payload.program.StringId(),
			S:  Errored,
		}
		panic(err)
	}
	editions := d.api.GetProgramEditions(payload.channel.Id, payload.program.StringId())
	sort.Sort(teacommon.BySize(editions))
	e := editions[0]
	url := e.Urls[0]
	homeDir, _ := os.UserHomeDir()
	downloadDir := filepath.Join(homeDir, "Downloads", payload.channel.Name+"-qtg")
	if os.Mkdir(downloadDir, os.ModePerm) != nil {
		if os.IsNotExist(os.Mkdir(downloadDir, os.ModePerm)) {
			panic("无法创建下载目录")
		}
	}
	req, err := grab.NewRequest(downloadDir, url)
	req.Filename = filepath.Join(downloadDir, fmt.Sprintf("%s.%s.%s", payload.program.NamePrefix, payload.program.Name, e.Format))
	req.NoResume = true
	if err != nil {
		handleErr(err)
	}
	resp := d.client.Do(req)
Loop:
	for {
		select {
		case <-resp.Done:
			// // println("完成下载 Program " + payload.program.Name)
			break Loop
		}
	}
	if err := resp.Err(); err != nil {
		handleErr(err)
	}
	err = d.writeTag(resp.Filename, payload)
	if err != nil {
		handleErr(err)
	}
	resultCh <- Status{
		Id: payload.program.StringId(),
		S:  Finished,
	}
}

func (d *Downloader) Wait() {
	d.wg.Wait()
}

func (d *Downloader) writeTag(file string, payload downloadPayload) error {
	m := meta.Metadata{
		Title:       payload.program.Name,
		Album:       payload.channel.Name,
		Artist:      payload.channel.Podcaster.Name,
		Description: payload.channel.Description(),
		Date:        payload.program.Date(),
		Cover:       payload.program.Cover,
	}
	return meta.WriteMetadata(file, m, true)
}
