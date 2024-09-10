package hlsdl

import (
  "bytes"
  "errors"
  "log"
  "net/url"
  "strings"
  "time"

  "github.com/go-resty/resty/v2"
  "github.com/grafov/m3u8"
)

func parseHlsSegments(hlsURL string, headers map[string]string) ([]*Segment, error) {
  baseURL, err := url.Parse(hlsURL)
  if err != nil {
    return nil, errors.New("Invalid m3u8 url")
  }

  p, t, err := getM3u8ListType(hlsURL, headers)
  if err != nil {
    log.Println("parse segemnt failed", err)
    return nil, err
  }
  if t != m3u8.MEDIA {
    return nil, errors.New("No support the m3u8 format")
  }

  mediaList := p.(*m3u8.MediaPlaylist)
  segments := []*Segment{}
  for _, seg := range mediaList.Segments {
    if seg == nil {
      continue
    }

    if !strings.Contains(seg.URI, "http") {
      segmentURL, err := baseURL.Parse(seg.URI)
      if err != nil {
        return nil, err
      }

      seg.URI = segmentURL.String()
    }

    if seg.Key == nil && mediaList.Key != nil {
      seg.Key = mediaList.Key
    }

    if seg.Key != nil && !strings.Contains(seg.Key.URI, "http") {
      keyURL, err := baseURL.Parse(seg.Key.URI)
      if err != nil {
        return nil, err
      }

      seg.Key.URI = keyURL.String()
    }

    segment := &Segment{MediaSegment: seg}
    segments = append(segments, segment)
  }

  return segments, nil
}

func getM3u8ListType(url string, headers map[string]string) (m3u8.Playlist, m3u8.ListType, error) {
  client := resty.New()
  client.SetRetryCount(5).SetRetryWaitTime(time.Second)
  resp, err := client.R().SetHeaders(headers).Get(url)
  if err != nil {
    return nil, 0, err
  }
  if resp.StatusCode() != 200 {
    return nil, 0, errors.New(resp.Status())
  }
  p, t, err := m3u8.DecodeFrom(bytes.NewReader(resp.Body()), false)
  if err != nil {
    return nil, 0, err
  }
  return p, t, nil
}
