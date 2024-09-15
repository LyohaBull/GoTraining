package main

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func getFileAudio(req_str, filename, size string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", req_str, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, identity")
	req.Header.Set("Referer", "https://kinescope.io/")
	//req.Header.Set("range", "bytes=0-"+size)
	req.Header.Set("Origin", "https://kinescope.io")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	binary.Write(f, binary.BigEndian, bodyText)
}

type key_arr struct {
	Kty string
	K   string
	Kid string
}

func (key key_arr) String() string {
	data, err := base64.StdEncoding.DecodeString(key.Kid[:20])
	if err != nil {
		fmt.Println("decode error:", err)
		return "decode error: " + err.Error()
	}
	dst1 := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst1, data)
	data1, err1 := base64.StdEncoding.DecodeString(key.K[:20])
	if err1 != nil {
		fmt.Println("decode error:", err)
		return "decode error: " + err.Error()
	}
	dst2 := make([]byte, hex.EncodedLen(len(data1)))
	hex.Encode(dst2, data1)
	return fmt.Sprintf("%s3d:%s3d\n", dst1, dst2)
}

type license struct {
	Keys  []key_arr
	Token string
	Type  string
}

func getKey(req_str, default_kid string) string {
	client := &http.Client{}
	var data = strings.NewReader(`{"kids":["` + default_kid + `"],"type":"temporary"}`)
	req, err := http.NewRequest("POST", req_str, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://kinescope.io/")
	req.Header.Set("Origin", "https://kinescope.io")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	lic := new(license)
	err = json.Unmarshal(bodyText, lic)
	if err != nil {
		log.Fatal(err)
	}
	return lic.Keys[0].String()
}

type Initialization struct {
	SourceURL string `xml:"sourceURL,attr"`
}

type ContentProtection struct {
	DefaultKID string `xml:"default_KID,attr"`
	Clearkey   string `xml:",any"`
}
type Result struct {
	Initialization    []Initialization    `xml:"Period>AdaptationSet>Representation>SegmentList>Initialization"`
	ContentProtection []ContentProtection `xml:"Period>AdaptationSet>ContentProtection"`
}

func getXML(request string) parsingXML {
	client := &http.Client{}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")

	req.Header.Set("Accept", "*/*")

	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://kinescope.io/embed/202386575")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	/*f, err := os.OpenFile("master.mpd", os.O_EXCL|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, _ := io.ReadAll(f)
	fmt.Println(len(bodyText))*/
	/*data := `
	<SegmentList timescale='1000' duration='4000'>
	<Initialization sourceURL='https://edge-msk-1.kinescopecdn.net/7cb7949b-23bf-466e-a8cd-fe37eafdb877/videos/aa6403aa-6ea8-4c8d-ba99-f379f6fd6671/assets/52d43af0-90d6-47ea-af3a-9ffc9ce1721a/720p.mp4?expires=1712074700&amp;kinescope_project_id=83553c7a-45df-4732-95f4-8a3e714f7899&amp;sign=e214e8d257cf09fc4e81ec2f36e6b59c' range='40-894'/>
	</SegmentList>
	`
	bodyText := []byte{}
	bodyText = append(bodyText, []byte(data)...)*/
	x := Result{}
	err = xml.Unmarshal([]byte(bodyText), &x)
	if err != nil {
		log.Fatal(err)
	}
	aud := ""
	vid := ""
	for _, s := range x.Initialization {
		if strings.Contains(s.SourceURL, "audio_0.mp4") {
			aud = s.SourceURL
		}
		if strings.Contains(s.SourceURL, "720p.mp4") {
			vid = s.SourceURL
		}
	}
	return parsingXML{
		license:     x.ContentProtection[1].Clearkey,
		default_kid: toBase64(x.ContentProtection[0].DefaultKID),
		video:       vid,
		audio:       aud,
	}
}

type parsingXML struct {
	license     string
	default_kid string
	audio       string
	video       string
}

func toBase64(str string) string {
	str = strings.Replace(str, "-", "", -1)
	src := []byte(str)

	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(dst))
	return encoded
}
func main() {
	fmt.Println("input filename:")
	in := bufio.NewReader(os.Stdin)
	filename, _ := in.ReadString('\n')
	fmt.Println("input mpd link:")
	xml := ""
	fmt.Scanf("%s \n", &xml)
	data := getXML(xml)
	//data := getXML("https://kinescope.io/2ee47009-1ad1-4bcf-beef-e1a3597aa5ff/master.mpd?expires=1711988589&kinescope_project_id=83553c7a-45df-4732-95f4-8a3e714f7899&sign=bd0a3d82b520b6ec0f7ecc69f4b05c0e")
	/*fmt.Println("input license url:")
	key := ""
	df := ""
	fmt.Scanf("%s %s\n", &key, &df)*/
	data.default_kid = data.default_kid[:len(data.default_kid)-2]
	keys := getKey(data.license, data.default_kid)
	commands := []exec.Cmd{
		*exec.Command("mp4decrypt", "--key", keys, "video_en.mp4", "video_de.mp4"),
		*exec.Command("mp4decrypt", "--key", keys, "audio_en.m4a", "audio_de.m4a"),
		*exec.Command("ffmpeg", "-i", "video_de.mp4", "-i", "audio_de.m4a", "-c", "copy", "fullvideo.mp4"),
	}
	for i := 0; i < 2; i++ {
		/*a := ""
		size := ""
		fmt.Println("input link and size:")
		fmt.Scanf("%s %s\n", &a, &size)*/
		if i == 0 {
			getFileAudio(data.video, "video_en.mp4", "")
		} else {
			getFileAudio(data.audio, "audio_en.m4a", "")
		}
	}
	for _, cmd := range commands {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	os.Remove("video_de.mp4")
	os.Remove("video_en.mp4")
	os.Remove("audio_en.m4a")
	os.Remove("audio_de.m4a")
	os.Rename("fullvideo.mp4", filename)

}
