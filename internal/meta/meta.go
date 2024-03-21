package meta

import (
	"fmt"
	"github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var coverCache = make(map[string][]byte)

type Metadata struct {
	Title       string
	Album       string
	Artist      string
	Description string
	Date        string
	Cover       string
}

func (m *Metadata) GetTypeOfTitle() mp4.BoxType {
	return mp4.BoxType{0xA9, 'n', 'a', 'm'}
}

func (m *Metadata) GetTypeOfAlbum() mp4.BoxType {
	return mp4.BoxType{0xA9, 'a', 'l', 'b'}
}

func (m *Metadata) GetTypeOfArtist() mp4.BoxType {
	return mp4.BoxType{0xA9, 'A', 'R', 'T'}
}

func (m *Metadata) GetTypeOfDescription() mp4.BoxType {
	return mp4.StrToBoxType("desc")
}

func (m *Metadata) GetTypeOfDate() mp4.BoxType {
	return mp4.BoxType{0xA9, 'd', 'a', 'y'}
}

func (m *Metadata) GetTypeOfCover() mp4.BoxType {
	return mp4.StrToBoxType("covr")
}

func (m *Metadata) AddMeta(w *mp4.Writer, ctx mp4.Context) error {
	if err := m.AddTitle(w, ctx); err != nil {
		return err
	}
	if err := m.AddAlbum(w, ctx); err != nil {
		return err
	}
	if err := m.AddArtist(w, ctx); err != nil {
		return err
	}
	if err := m.AddDescription(w, ctx); err != nil {
		return err
	}
	if err := m.AddDate(w, ctx); err != nil {
		return err
	}
	if err := m.AddCover(w, ctx); err != nil {
		return err
	}
	return nil
}

func (m *Metadata) AddTitle(w *mp4.Writer, ctx mp4.Context) error {
	return addMeta(w, ctx, m.GetTypeOfTitle(), &mp4.Data{Data: []byte(m.Title), DataType: mp4.DataTypeStringUTF8})
}

func (m *Metadata) AddAlbum(w *mp4.Writer, ctx mp4.Context) error {
	return addMeta(w, ctx, m.GetTypeOfAlbum(), &mp4.Data{Data: []byte(m.Album), DataType: mp4.DataTypeStringUTF8})
}

func (m *Metadata) AddArtist(w *mp4.Writer, ctx mp4.Context) error {
	return addMeta(w, ctx, m.GetTypeOfArtist(), &mp4.Data{Data: []byte(m.Artist), DataType: mp4.DataTypeStringUTF8})
}

func (m *Metadata) AddDescription(w *mp4.Writer, ctx mp4.Context) error {
	return addMeta(w, ctx, m.GetTypeOfDescription(), &mp4.Data{Data: []byte(m.Description), DataType: mp4.DataTypeStringUTF8})
}

func (m *Metadata) AddDate(w *mp4.Writer, ctx mp4.Context) error {
	return addMeta(w, ctx, m.GetTypeOfDate(), &mp4.Data{Data: []byte(m.Date), DataType: mp4.DataTypeStringUTF8})
}

func (m *Metadata) AddCover(w *mp4.Writer, ctx mp4.Context) error {
	pic, err := getCover(m.Cover)
	if err != nil {
		return err
	}
	return addMeta(w, ctx, m.GetTypeOfCover(), &mp4.Data{Data: pic, DataType: mp4.DataTypeStringJPEG})
}

func WriteMetadata(file string, metadata Metadata, rewrite bool) error {
	input, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func(input *os.File) {
		_ = input.Close()
	}(input)

	out, err := os.Create(file + ".tmp")
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()

		if rewrite {
			err = os.Remove(file)
			if err != nil {
				panic(err)
			}
			err = os.Rename(file+".tmp", file)
			if err != nil {
				panic(err)
			}
		} else {
			err = os.Rename(file+".tmp", strings.Replace(file, filepath.Ext(file), "[meta].m4a", 1))
			if err != nil {
				panic(err)
			}
		}
	}(out)

	var ilstExists bool
	var mdatOffsetDiff int64
	var stcoOffsets []int64

	r := bufseekio.NewReadSeeker(input, 1024*1024, 3)
	w := mp4.NewWriter(out)

	_, err = mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		switch h.BoxInfo.Type {
		case mp4.BoxTypeMoov(),
			mp4.BoxTypeTrak(),
			mp4.BoxTypeMdia(),
			mp4.BoxTypeMinf(),
			mp4.BoxTypeStbl(),
			mp4.BoxTypeUdta(),
			mp4.BoxTypeMeta(),
			mp4.BoxTypeIlst():
			_, err := w.StartBox(&h.BoxInfo)
			if err != nil {
				return nil, err
			}
			if _, err := h.Expand(); err != nil {
				return nil, err
			}

			if h.BoxInfo.Type == mp4.BoxTypeMoov() && !ilstExists {
				path := []mp4.BoxType{mp4.BoxTypeUdta(), mp4.BoxTypeMeta()}
				for _, boxType := range path {
					if _, err := w.StartBox(&mp4.BoxInfo{Type: boxType}); err != nil {
						return nil, err
					}
				}
				ctx := h.BoxInfo.Context
				ctx.UnderUdta = true
				if _, err := w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeHdlr()}); err != nil {
					return nil, err
				}
				hdlr := &mp4.Hdlr{
					HandlerType: [4]byte{'m', 'd', 'i', 'r'},
				}
				if _, err := mp4.Marshal(w, hdlr, ctx); err != nil {
					return nil, err
				}
				if _, err := w.EndBox(); err != nil {
					return nil, err
				}
				if _, err := w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeIlst()}); err != nil {
					return nil, err
				}
				ctx.UnderIlst = true
				if err := metadata.AddMeta(w, ctx); err != nil {
					return nil, err
				}
				if _, err := w.EndBox(); err != nil {
					return nil, err
				}
				for range path {
					if _, err := w.EndBox(); err != nil {
						return nil, err
					}
				}
			}

			if h.BoxInfo.Type == mp4.BoxTypeIlst() {
				ctx := h.BoxInfo.Context
				ctx.UnderIlst = true
				if err := metadata.AddMeta(w, ctx); err != nil {
					return nil, err
				}
				ilstExists = true
			}
			if _, err = w.EndBox(); err != nil {
				return nil, err
			}
		default:
			if h.BoxInfo.Type == mp4.BoxTypeStco() {
				offset, _ := w.Seek(0, io.SeekCurrent)
				stcoOffsets = append(stcoOffsets, offset)
			}
			if h.BoxInfo.Type == mp4.BoxTypeMdat() {
				iOffset := int64(h.BoxInfo.Offset)
				oOffset, _ := w.Seek(0, io.SeekCurrent)
				mdatOffsetDiff = oOffset - iOffset
			}
			if err := w.CopyBox(r, &h.BoxInfo); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	// if mdat box is moved, update stco box
	if mdatOffsetDiff != 0 {
		for _, stcoOffset := range stcoOffsets {
			// seek to stco box header
			if _, err := out.Seek(stcoOffset, io.SeekStart); err != nil {
				panic(err)
			}
			// read box header
			bi, err := mp4.ReadBoxInfo(out)
			if err != nil {
				panic(err)
			}
			// read stco box payload
			var stco mp4.Stco
			if _, err := mp4.Unmarshal(out, bi.Size-bi.HeaderSize, &stco, bi.Context); err != nil {
				panic(err)
			}
			// update chunk offsets
			for i := range stco.ChunkOffset {
				stco.ChunkOffset[i] += uint32(mdatOffsetDiff)
			}
			// seek to stco box payload
			if _, err := bi.SeekToPayload(out); err != nil {
				panic(err)
			}
			// write stco box payload
			if _, err := mp4.Marshal(out, &stco, bi.Context); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func addMeta(w *mp4.Writer, ctx mp4.Context, boxType mp4.BoxType, data *mp4.Data) error {
	if _, err := w.StartBox(&mp4.BoxInfo{Type: boxType}); err != nil {
		return err
	}
	if _, err := w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeData()}); err != nil {
		return err
	}
	dataCtx := ctx
	dataCtx.UnderIlstMeta = true
	if _, err := mp4.Marshal(w, data, dataCtx); err != nil {
		return err
	}
	if _, err := w.EndBox(); err != nil {
		return err
	}
	_, err := w.EndBox()
	return err
}

func getCover(url string) ([]byte, error) {
	if c, ok := coverCache[url]; ok {
		return c, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// 读取整个响应体
	pic, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	coverCache[url] = pic
	return pic, nil
}
