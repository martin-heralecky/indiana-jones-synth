package main

import (
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"
)

// in seconds
type time float64

// in [-1,1]
type sample float64

type sound interface {
	sample(t time) sample
}

type noSound struct{}

func (s noSound) sample(t time) sample {
	return 0
}

type wave interface {
	// t is in [0,1). 0 = wave start. 1 = wave end.
	sample(t time) sample
}

type sinWave struct{}

func (w sinWave) sample(t time) sample {
	return sample(math.Sin(float64(t) * 2 * math.Pi))
}

type squareWave struct{}

func (w squareWave) sample(t time) sample {
	if t < .5 {
		return -1
	} else {
		return 1
	}
}

type sawWave struct{}

func (w sawWave) sample(t time) sample {
	return sample(math.Abs(4*float64(t)-2) - 1)
}

// infinite repeating wave at given frequency
type waveSound struct {
	w wave
	f float64 // hz
}

func (s waveSound) sample(t time) sample {
	waveLength := 1 / s.f
	return s.w.sample(time(math.Mod(float64(t), waveLength) * s.f))
}

// adds fade-in and fade-out effect to another sound
type fadeInFadeOutDecorator struct {
	s sound
	soundDuration time
	fadeDuration time
}

func (d fadeInFadeOutDecorator) sample(t time) sample {
	s := d.s.sample(t)

	if t <= d.fadeDuration {
		// fade in
		s = sample(float64(s) * float64(t * 1/d.fadeDuration))
	} else if t >= d.soundDuration - d.fadeDuration {
		// fade out
		s = sample(float64(s) * float64(-(t - (d.soundDuration - d.fadeDuration)) * 1/d.fadeDuration + 1))
	}

	return s
}

func play(s sound, duration float64, rate uint, file *os.File) {
	samples := int(duration * float64(rate))
	for t := 0; t < samples; t++ {
		time := time(float64(t) / float64(rate))
		sample := s.sample(time)

		encodedSample := int16(math.Round(float64(sample+1)*(1<<15-.5)) - 1<<15)
		err := binary.Write(file, binary.LittleEndian, encodedSample)
		if err != nil {
			panic(err)
		}
	}
}

func makeSong(song string, wave wave, speed float64, filename string, rate uint) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sounds := map[string]sound{
		"-":   noSound{},
		"c0":  waveSound{wave, 16.35},
		"c#0": waveSound{wave, 17.32},
		"db0": waveSound{wave, 17.32},
		"d0":  waveSound{wave, 18.35},
		"d#0": waveSound{wave, 19.45},
		"eb0": waveSound{wave, 19.45},
		"e0":  waveSound{wave, 20.60},
		"f0":  waveSound{wave, 21.83},
		"f#0": waveSound{wave, 23.12},
		"gb0": waveSound{wave, 23.12},
		"g0":  waveSound{wave, 24.50},
		"g#0": waveSound{wave, 25.96},
		"ab0": waveSound{wave, 25.96},
		"a0":  waveSound{wave, 27.50},
		"a#0": waveSound{wave, 29.14},
		"bb0": waveSound{wave, 29.14},
		"b0":  waveSound{wave, 30.87},
		"c1":  waveSound{wave, 32.70},
		"c#1": waveSound{wave, 34.65},
		"db1": waveSound{wave, 34.65},
		"d1":  waveSound{wave, 36.71},
		"d#1": waveSound{wave, 38.89},
		"eb1": waveSound{wave, 38.89},
		"e1":  waveSound{wave, 41.20},
		"f1":  waveSound{wave, 43.65},
		"f#1": waveSound{wave, 46.25},
		"gb1": waveSound{wave, 46.25},
		"g1":  waveSound{wave, 49.00},
		"g#1": waveSound{wave, 51.91},
		"ab1": waveSound{wave, 51.91},
		"a1":  waveSound{wave, 55.00},
		"a#1": waveSound{wave, 58.27},
		"bb1": waveSound{wave, 58.27},
		"b1":  waveSound{wave, 61.74},
		"c2":  waveSound{wave, 65.41},
		"c#2": waveSound{wave, 69.30},
		"db2": waveSound{wave, 69.30},
		"d2":  waveSound{wave, 73.42},
		"d#2": waveSound{wave, 77.78},
		"eb2": waveSound{wave, 77.78},
		"e2":  waveSound{wave, 82.41},
		"f2":  waveSound{wave, 87.31},
		"f#2": waveSound{wave, 92.50},
		"gb2": waveSound{wave, 92.50},
		"g2":  waveSound{wave, 98.00},
		"g#2": waveSound{wave, 103.83},
		"ab2": waveSound{wave, 103.83},
		"a2":  waveSound{wave, 110.00},
		"a#2": waveSound{wave, 116.54},
		"bb2": waveSound{wave, 116.54},
		"b2":  waveSound{wave, 123.47},
		"c3":  waveSound{wave, 130.81},
		"c#3": waveSound{wave, 138.59},
		"db3": waveSound{wave, 138.59},
		"d3":  waveSound{wave, 146.83},
		"d#3": waveSound{wave, 155.56},
		"eb3": waveSound{wave, 155.56},
		"e3":  waveSound{wave, 164.81},
		"f3":  waveSound{wave, 174.61},
		"f#3": waveSound{wave, 185.00},
		"gb3": waveSound{wave, 185.00},
		"g3":  waveSound{wave, 196.00},
		"g#3": waveSound{wave, 207.65},
		"ab3": waveSound{wave, 207.65},
		"a3":  waveSound{wave, 220.00},
		"a#3": waveSound{wave, 233.08},
		"bb3": waveSound{wave, 233.08},
		"b3":  waveSound{wave, 246.94},
		"c4":  waveSound{wave, 261.63},
		"c#4": waveSound{wave, 277.18},
		"db4": waveSound{wave, 277.18},
		"d4":  waveSound{wave, 293.66},
		"d#4": waveSound{wave, 311.13},
		"eb4": waveSound{wave, 311.13},
		"e4":  waveSound{wave, 329.63},
		"f4":  waveSound{wave, 349.23},
		"f#4": waveSound{wave, 369.99},
		"gb4": waveSound{wave, 369.99},
		"g4":  waveSound{wave, 392.00},
		"g#4": waveSound{wave, 415.30},
		"ab4": waveSound{wave, 415.30},
		"a4":  waveSound{wave, 440.00},
		"a#4": waveSound{wave, 466.16},
		"bb4": waveSound{wave, 466.16},
		"b4":  waveSound{wave, 493.88},
		"c5":  waveSound{wave, 523.25},
		"c#5": waveSound{wave, 554.37},
		"db5": waveSound{wave, 554.37},
		"d5":  waveSound{wave, 587.33},
		"d#5": waveSound{wave, 622.25},
		"eb5": waveSound{wave, 622.25},
		"e5":  waveSound{wave, 659.25},
		"f5":  waveSound{wave, 698.46},
		"f#5": waveSound{wave, 739.99},
		"gb5": waveSound{wave, 739.99},
		"g5":  waveSound{wave, 783.99},
		"g#5": waveSound{wave, 830.61},
		"ab5": waveSound{wave, 830.61},
		"a5":  waveSound{wave, 880.00},
		"a#5": waveSound{wave, 932.33},
		"bb5": waveSound{wave, 932.33},
		"b5":  waveSound{wave, 987.77},
		"c6":  waveSound{wave, 1046.50},
		"c#6": waveSound{wave, 1108.73},
		"db6": waveSound{wave, 1108.73},
		"d6":  waveSound{wave, 1174.66},
		"d#6": waveSound{wave, 1244.51},
		"eb6": waveSound{wave, 1244.51},
		"e6":  waveSound{wave, 1318.51},
		"f6":  waveSound{wave, 1396.91},
		"f#6": waveSound{wave, 1479.98},
		"gb6": waveSound{wave, 1479.98},
		"g6":  waveSound{wave, 1567.98},
		"g#6": waveSound{wave, 1661.22},
		"ab6": waveSound{wave, 1661.22},
		"a6":  waveSound{wave, 1760.00},
		"a#6": waveSound{wave, 1864.66},
		"bb6": waveSound{wave, 1864.66},
		"b6":  waveSound{wave, 1975.53},
		"c7":  waveSound{wave, 2093.00},
		"c#7": waveSound{wave, 2217.46},
		"db7": waveSound{wave, 2217.46},
		"d7":  waveSound{wave, 2349.32},
		"d#7": waveSound{wave, 2489.02},
		"eb7": waveSound{wave, 2489.02},
		"e7":  waveSound{wave, 2637.02},
		"f7":  waveSound{wave, 2793.83},
		"f#7": waveSound{wave, 2959.96},
		"gb7": waveSound{wave, 2959.96},
		"g7":  waveSound{wave, 3135.96},
		"g#7": waveSound{wave, 3322.44},
		"ab7": waveSound{wave, 3322.44},
		"a7":  waveSound{wave, 3520.00},
		"a#7": waveSound{wave, 3729.31},
		"bb7": waveSound{wave, 3729.31},
		"b7":  waveSound{wave, 3951.07},
		"c8":  waveSound{wave, 4186.01},
		"c#8": waveSound{wave, 4434.92},
		"db8": waveSound{wave, 4434.92},
		"d8":  waveSound{wave, 4698.63},
		"d#8": waveSound{wave, 4978.03},
		"eb8": waveSound{wave, 4978.03},
		"e8":  waveSound{wave, 5274.04},
		"f8":  waveSound{wave, 5587.65},
		"f#8": waveSound{wave, 5919.91},
		"gb8": waveSound{wave, 5919.91},
		"g8":  waveSound{wave, 6271.93},
		"g#8": waveSound{wave, 6644.88},
		"ab8": waveSound{wave, 6644.88},
		"a8":  waveSound{wave, 7040.00},
		"a#8": waveSound{wave, 7458.62},
		"bb8": waveSound{wave, 7458.62},
		"b8":  waveSound{wave, 7902.13},
	}

	for _, tone := range strings.Split(song, " ") {
		parts := strings.Split(tone, "x")
		duration, _ := strconv.ParseFloat(parts[0], 64)
		realDuration := (1 / speed) * duration
		s := fadeInFadeOutDecorator{sounds[parts[1]], time(realDuration), 0.04}
		play(s, realDuration, rate, f)
	}
}

func main() {
	const outFilename = "indiana_jones.pcm"
	const rate = 44100
	const speed = 7.3

	wave := sawWave{}

	song := "4x- 3xe4 1xf4 2xg4 10xc5 3xd4 1xe4 12xf4 3xg4 1xa4 2xb4 10xf5 3xa4 1xb4 4xc5 4xd5 4xe5 3xe4 1xf4 2xg4 10xc5 3xd5 1xe5 12xf5 3xg4 1xg4 4xe5 3xd5 1xg4 4xe5 3xd5 1xg4 4xf5 3xe5 1xd5 2xc5 6x-"

	makeSong(song, wave, speed, outFilename, rate)
}
