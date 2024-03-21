package teacommon

import "strconv"

type VipInfo struct {
	Vip    bool   `json:"vip"`
	Expire string `json:"expire_time"`
}

type PrivateInfo struct {
	VipInfo VipInfo `json:"vip_info"`
}

type UserInfo struct {
	Nickname    string      `json:"nick_name"`
	AvatarUrl   string      `json:"avatar"`
	PrivateInfo PrivateInfo `json:"private_info"`
}

type Podcaster struct {
	Name string `json:"name"`
}
type Channel struct {
	Id        string    `json:"id"`
	Name      string    `json:"title"`
	Finished  bool      `json:"is_finished"`
	Desc      string    `json:"desc"`
	Count     int       `json:"program_count"`
	Podcaster Podcaster `json:"podcaster"`
}

func (c Channel) FilterValue() string {
	return c.Name
}

func (c Channel) Title() string {
	return c.Name
}

func (c Channel) Description() string {
	return c.Desc
}

func (c Channel) Author() string {
	return c.Podcaster.Name
}

type SearchResult struct {
	Keyword  string
	Type     string
	Page     int       `json:"page"`
	PageSize int       `json:"pagesize"`
	Total    int       `json:"total"`
	Data     []Channel `json:"data"`
}

type DetailChannel struct {
	Id    int `json:"id"`
	Count int `json:"program_count"`
}

type Program struct {
	Id         int    `json:"id"`
	Name       string `json:"title"`
	Cover      string `json:"cover"`
	UpdateTime string `json:"update_time"`
	Index      int
	NamePrefix string
}

func (p *Program) Date() string {
	return p.UpdateTime[:4]
}

func (p *Program) StringId() string {
	return strconv.Itoa(p.Id)
}

type Edition struct {
	Format  string   `json:"format"`
	Bitrate int      `json:"bitrate"`
	Size    int      `json:"size"`
	Urls    []string `json:"urls"`
}

type BySize []Edition

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size > a[j].Size }
