package main

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

var audioContext *audio.Context

//go:embed assets/blocked.ogg
var sfx_blocked_byte []byte
var sfx_blocked []*audio.Player

//go:embed assets/extracted.ogg
var sfx_extracted_byte []byte
var sfx_extracted []*audio.Player

//go:embed assets/level_complete.ogg
var sfx_level_complete_byte []byte
var sfx_level_complete []*audio.Player

//go:embed assets/menu.ogg
var sfx_menu_byte []byte
var sfx_menu []*audio.Player

//go:embed assets/game_over.ogg
var sfx_game_over_byte []byte
var sfx_game_over []*audio.Player

//go:embed assets/start.ogg
var sfx_start_byte []byte
var sfx_start []*audio.Player

func soundInit() {
	if audioContext == nil {
		audioContext = audio.NewContext(48000)
	}

	sfx_menu = make([]*audio.Player, 10)
	for i := 0; i < 10; i++ {
		tmpDecoded, _ := vorbis.DecodeWithSampleRate(48000, bytes.NewReader(sfx_menu_byte))
		sfx_menu[i], _ = audioContext.NewPlayer(tmpDecoded)
		//sfx_menu[i].SetVolume(0.2)
	}

	sfx_blocked = make([]*audio.Player, 10)
	for i := 0; i < 10; i++ {
		tmpDecoded, _ := vorbis.DecodeWithSampleRate(48000, bytes.NewReader(sfx_blocked_byte))
		sfx_blocked[i], _ = audioContext.NewPlayer(tmpDecoded)
	}

	sfx_extracted = make([]*audio.Player, 6)
	for i := 0; i < 6; i++ {
		tmpDecoded, _ := vorbis.DecodeWithSampleRate(48000, bytes.NewReader(sfx_extracted_byte))
		sfx_extracted[i], _ = audioContext.NewPlayer(tmpDecoded)
		sfx_extracted[i].SetVolume(0.5)
	}

	sfx_level_complete = make([]*audio.Player, 1)
	for i := 0; i < 1; i++ {
		tmpDecoded, _ := vorbis.DecodeWithSampleRate(48000, bytes.NewReader(sfx_level_complete_byte))
		sfx_level_complete[i], _ = audioContext.NewPlayer(tmpDecoded)
	}

	sfx_game_over = make([]*audio.Player, 1)
	for i := 0; i < 1; i++ {
		tmpDecoded, _ := vorbis.DecodeWithSampleRate(48000, bytes.NewReader(sfx_game_over_byte))
		sfx_game_over[i], _ = audioContext.NewPlayer(tmpDecoded)
		sfx_game_over[i].SetVolume(0.5)
	}

	sfx_start = make([]*audio.Player, 1)
	for i := 0; i < 1; i++ {
		tmpDecoded, _ := vorbis.DecodeWithSampleRate(48000, bytes.NewReader(sfx_start_byte))
		sfx_start[i], _ = audioContext.NewPlayer(tmpDecoded)
		sfx_start[i].SetVolume(0.5)
	}
}

func soundPlay(soundBank []*audio.Player) {
	for i := 0; i < len(soundBank); i++ {
		if !soundBank[i].IsPlaying() {
			soundBank[i].Rewind()
			soundBank[i].Play()
			return
		}
	}
}
