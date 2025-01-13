package veado

import (
	"io"
	"reflect"

	"github.com/bake/bin"
)

func Read(r io.Reader) (*Veado, error) {
	vr := bin.NewReader(r)
	v := &Veado{}
	if err := vr.Read(v); err != nil {
		return nil, err
	}
	return v, nil
}

type String struct {
	Length uint64 `bin:"uvarint,ref"`
	Value  string `bin:",size=.Length"`
}

type Header struct {
	Magic string `bin:",size=9"`
}

type Meta struct {
	Software    String
	Credits     String
	Description String
}

func (c Meta) Skip(v reflect.Value) bool {
	return v.String() != "META"
}

type Mlst struct {
	ChunkIDs []uint32 `bin:",size=EOF"`
}

func (c Mlst) Skip(v reflect.Value) bool {
	return v.String() != "MLST"
}

type EffectFlag uint8

const (
	EffectFlagActive         EffectFlag = 0x1
	EffectFlagUsePreset      EffectFlag = 0x2
	EffectFlagUsePresetChunk EffectFlag = 0x4
)

type EffectChunkID struct {
	Value uint32
}

func (e EffectChunkID) Skip(v reflect.Value) bool {
	return uint8(v.Uint())&0x04 == 0
}

type EffectPresetID struct {
	Value String
}

func (e EffectPresetID) Skip(v reflect.Value) bool {
	return uint8(v.Uint())&0x02 == 0
}

type Effect struct {
	ID        String
	Flag      uint8          `bin:",ref"` // EffectFlag
	ChunkID   EffectChunkID  `bin:",skip=.Flag"`
	PresetID  EffectPresetID `bin:",skip=.Flag"`
	NumValues uint64         `bin:"uvarint,ref"`
	Values    []float64      `bin:",size=.NumValues"`
}

type Effects struct {
	NumEffects uint64   `bin:"uvarint,ref"`
	Effects    []Effect `bin:",size=.NumEffects"`
}

type StateFlag uint32

const (
	StateFlagPixelated StateFlag = 0x1
	StateFlagBlink     StateFlag = 0x2
	StateFlagStart     StateFlag = 0x4
)

type Signal struct {
	Source String
	Name   String
}

type ShortcutMode string

const (
	ShortcutModePress        ShortcutMode = "PRES"
	ShortcutModeRelease      ShortcutMode = "RLSE"
	ShortcutModeWhilePressed ShortcutMode = "PRED"
)

type Msta struct {
	Name                         String
	State                        uint32 // StateFlag
	ThumbnailClosedMouth         uint32
	ThumbnailOpenMouth           uint32
	ThumbnailBlinkingClosedMouth uint32
	ThumbnailBlinkingOpenMouth   uint32
	ClosedMouth                  uint32
	OpenMouth                    uint32
	BlinkingClosedMouth          uint32
	BlinkingOpenMouth            uint32
	BlinkDuration                float64
	MinBlinkInterval             float64
	MaxBlinkInterval             float64
	ClosedMouthEffects           Effects
	OpenMouthEffects             Effects
	OnOpenMouthEffects           Effects
	OnCloseMouthEffects          Effects
	NumSignals                   uint64   `bin:"uvarint,ref"`
	Signals                      []Signal `bin:",size=.NumSignals"`
	ShortcutMode                 string   `bin:",size=4"` // ShortcutMode
}

func (c Msta) Skip(v reflect.Value) bool {
	return v.String() != "MSTA"
}

type AsfdMetadata struct {
	Type        string `bin:",size=4"`
	NumMetadata uint64 `bin:"uvarint,ref"`
	Metadata    []byte `bin:",size=.NumMetadata"`
}

type AsfdEntry struct {
	EntryName   String
	ChunkID     uint32
	NumMetadata uint64         `bin:"uvarint,ref"`
	Metadata    []AsfdMetadata `bin:",size=.NumMetadata"`
}

type Asfd struct {
	RootCode string      `bin:",size=4"`
	Entries  []AsfdEntry `bin:",size=EOF"`
}

func (c Asfd) Skip(v reflect.Value) bool {
	return v.String() != "ASFD"
}

type Thmb struct {
	Data []byte `bin:",size=EOF"`
}

func (c Thmb) Skip(v reflect.Value) bool {
	return v.String() != "THMB"
}

type AimgFrame struct {
	ChunkID  uint32
	OffsetX  int32
	OffsetY  int32
	Duration float64
}

type AimgNumLoops struct {
	Value uint64 `bin:"uvarint"`
}

func (a AimgNumLoops) Skip(v reflect.Value) bool {
	return v.Uint() <= 1
}

type Aimg struct {
	Width     uint32
	Height    uint32
	NumFrames uint64       `bin:"uvarint,ref"`
	NumLoops  AimgNumLoops `bin:",skip=.NumFrames"`
	Frames    []AimgFrame  `bin:",size=.NumFrames"`
}

func (c Aimg) Skip(v reflect.Value) bool {
	return v.String() != "AIMG"
}

type Vdd []uint32

func (c Vdd) Skip(v reflect.Value) bool {
	return v.Uint() >= 0xffffff00
}

type Abmp struct {
	Width  uint32
	Height uint32
	Format string `bin:",size=4"`
	NumAs  uint32
	NumRs  uint32
	NumGs  uint32
	NumBs  uint32
	As     Vdd `bin:",size=.NumAs,skip=.NumAs"`
	Rs     Vdd `bin:",size=.NumRs,skip=.NumRs"`
	Gs     Vdd `bin:",size=.NumGs,skip=.NumGs"`
	Bs     Vdd `bin:",size=.NumBs,skip=.NumBs"`
}

func (c Abmp) Skip(v reflect.Value) bool {
	return v.String() != "ABMP"
}

type Chunk struct {
	ID     uint32
	Type   string `bin:",ref,size=4"`
	Length uint32 `bin:",ref"`
	Meta   Meta   `bin:",size=.Length,skip=.Type"`
	Mlst   Mlst   `bin:",size=.Length,skip=.Type"`
	Msta   Msta   `bin:",size=.Length,skip=.Type"`
	Asfd   Asfd   `bin:",size=.Length,skip=.Type"`
	Thmb   Thmb   `bin:",size=.Length,skip=.Type"`
	Aimg   Aimg   `bin:",size=.Length,skip=.Type"`
	Abmp   Abmp   `bin:",size=.Length,skip=.Type"`
}

type Veado struct {
	Header
	Chunks []Chunk `bin:",size=EOF"`
}
