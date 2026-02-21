package tui

import (
	"fmt"
	"io"

	"github.com/gopxl/beep/v2"
	fdkaac "github.com/qrtc/fdk-aac-go"
	mp4 "github.com/yapingcat/gomedia/go-mp4"
)

type aacStream struct {
	f          io.ReadSeeker
	demuxer    *mp4.MovDemuxer
	decoder    *fdkaac.AacDecoder
	format     beep.Format
	sampleRate int
	numChans   int

	pcmBuffer []byte
	pcmPos    int

	pos int
	len int
}

func DecodeAAC(f io.ReadSeeker) (beep.StreamSeekCloser, beep.Format, error) {
	demuxer := mp4.CreateMp4Demuxer(f)
	tracks, err := demuxer.ReadHead()
	if err != nil {
		return nil, beep.Format{}, err
	}

	var audioTrack *mp4.TrackInfo
	for i := range tracks {
		if tracks[i].Cid == mp4.MP4_CODEC_AAC {
			audioTrack = &tracks[i]
			break
		}
	}

	if audioTrack == nil {
		return nil, beep.Format{}, fmt.Errorf("no AAC track found")
	}

	totalSamples := int(audioTrack.SampleCount) * 1024

	decoder, err := fdkaac.CreateAccDecoder(&fdkaac.AacDecoderConfig{
		TransportFmt: fdkaac.TtMp4Raw,
	})
	if err != nil {
		return nil, beep.Format{}, err
	}

	decoder.Close()
	decoder, err = fdkaac.CreateAccDecoder(&fdkaac.AacDecoderConfig{
		TransportFmt: fdkaac.TtMp4Adts,
	})
	if err != nil {
		return nil, beep.Format{}, err
	}

	format := beep.Format{
		SampleRate:  beep.SampleRate(audioTrack.SampleRate),
		NumChannels: int(audioTrack.ChannelCount),
		Precision:   2,
	}

	return &aacStream{
		f:          f,
		demuxer:    demuxer,
		decoder:    decoder,
		format:     format,
		sampleRate: int(audioTrack.SampleRate),
		numChans:   int(audioTrack.ChannelCount),
		len:        totalSamples,
	}, format, nil
}

func (s *aacStream) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		if s.pcmPos >= len(s.pcmBuffer) {
			err := s.decodeNextFrame()
			if err == io.EOF {
				return i, i > 0
			}
			if err != nil {
				return i, i > 0
			}
			if len(s.pcmBuffer) == 0 {
				return i, i > 0
			}
		}

		if s.numChans >= 2 {
			if s.pcmPos+3 < len(s.pcmBuffer) {
				l := int16(s.pcmBuffer[s.pcmPos]) | (int16(s.pcmBuffer[s.pcmPos+1]) << 8)
				r := int16(s.pcmBuffer[s.pcmPos+2]) | (int16(s.pcmBuffer[s.pcmPos+3]) << 8)
				s.pcmPos += 4
				samples[i][0] = float64(l) / 32768
				samples[i][1] = float64(r) / 32768
			} else {
				s.pcmPos = len(s.pcmBuffer)
				continue
			}
		} else {
			if s.pcmPos+1 < len(s.pcmBuffer) {
				v := int16(s.pcmBuffer[s.pcmPos]) | (int16(s.pcmBuffer[s.pcmPos+1]) << 8)
				s.pcmPos += 2
				vf := float64(v) / 32768
				samples[i][0] = vf
				samples[i][1] = vf
			} else {
				s.pcmPos = len(s.pcmBuffer)
				continue
			}
		}
		s.pos++
	}
	return len(samples), true
}

func (s *aacStream) decodeNextFrame() error {
	for {
		pkt, err := s.demuxer.ReadPacket()
		if err != nil {
			return err
		}
		if pkt.Cid != mp4.MP4_CODEC_AAC {
			continue
		}

		target := make([]byte, 16384)
		decodedBytes, err := s.decoder.DecodeFrame(pkt.Data, target)

		if err == fdkaac.DecNotEnoughBits {
			continue
		}
		if err != nil {
			continue
		}

		if decodedBytes > 0 {
			s.pcmBuffer = target[:decodedBytes]
			s.pcmPos = 0
			return nil
		}
	}
}

func (s *aacStream) Err() error {
	return nil
}

func (s *aacStream) Len() int {
	return s.len
}

func (s *aacStream) Position() int {
	return s.pos
}

func (s *aacStream) Seek(p int) error {
	if p < 0 || p > s.len {
		return fmt.Errorf("seek out of bounds")
	}

	dtsOffset := uint64(p) * 1000 / uint64(s.sampleRate)
	err := s.demuxer.SeekTime(dtsOffset)
	if err != nil {
		return err
	}

	s.pos = p
	s.pcmBuffer = nil
	s.pcmPos = 0

	s.decoder.Flush()

	return nil
}

func (s *aacStream) Close() error {
	if s.decoder != nil {
		s.decoder.Close()
	}
	if c, ok := s.f.(io.Closer); ok {
		c.Close()
	}
	return nil
}
