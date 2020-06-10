package torrent

import (
	"reflect"
	"testing"

	"github.com/hekmon/cunits/v2"
	"github.com/hekmon/transmissionrpc"
	"github.com/shric/go-trpc/pkg/config"
)

func MakeTorrent(SizeWhenDone cunits.Bits, Eta int64, RecheckProgress float64,
	LeftUntilDone int64) (torrent *Torrent) {
	ID := int64(1)
	Error := int64(0)
	Name := "Torrent 1"
	RateUpload := int64(0)
	RateDownload := int64(0)
	UploadedEver := int64(0)
	BandwidthPriority := int64(0)
	Status := transmissionrpc.TorrentStatus(0)
	trpcTorrent := transmissionrpc.Torrent{
		SizeWhenDone:      &SizeWhenDone,
		ID:                &ID,
		Name:              &Name,
		Error:             &Error,
		Eta:               &Eta,
		RateUpload:        &RateUpload,
		RateDownload:      &RateDownload,
		RecheckProgress:   &RecheckProgress,
		LeftUntilDone:     &LeftUntilDone,
		Status:            &Status,
		UploadedEver:      &UploadedEver,
		BandwidthPriority: &BandwidthPriority,
	}
	conf := config.Config{
		Trackernames: map[string]string{
			"foo": "foo-tracker",
		},
	}
	torrent = NewFrom(&trpcTorrent, &conf)
	return
}

func TestEta(t *testing.T) {
	type test struct {
		input *Torrent
		want  string
	}
	tests := []test{
		{input: MakeTorrent(cunits.Bits(1), 1, 0.0, 1), want: "1 secs"},
		{input: MakeTorrent(cunits.Bits(1), 180, 0.0, 1), want: "180 secs"},
		{input: MakeTorrent(cunits.Bits(1), 240, 0.0, 1), want: "4 mins"},
		{input: MakeTorrent(cunits.Bits(1), 10000, 0.0, 1), want: "166 mins"},
		{input: MakeTorrent(cunits.Bits(1), 20000, 0.0, 1), want: "5 hours"},
		{input: MakeTorrent(cunits.Bits(1), 100000, 0.0, 1), want: "27 hours"},
		{input: MakeTorrent(cunits.Bits(1), 1000000, 0.0, 1), want: "11 days"},
		{input: MakeTorrent(cunits.Bits(1), 10000000, 0.0, 1), want: "3 months"},
		{input: MakeTorrent(cunits.Bits(1), 100000000, 0.0, 1), want: "3 years"},
		{input: MakeTorrent(cunits.Bits(1), -1, 0.0, 0), want: "Done"},
		{input: MakeTorrent(cunits.Bits(1), 1, 0.0, 0), want: "Done"},
		{input: MakeTorrent(cunits.Bits(1), -1, 0.0, 1), want: "Unknown"},
	}
	for _, tc := range tests {
		got := tc.input.eta()
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
