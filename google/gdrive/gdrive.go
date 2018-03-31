package drive

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type GoogleDrive struct {
	config  *oauth2.Config
	code    string
	service *drive.Service
	c       context.Context
}

func (p *GoogleDrive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.code = r.URL.Query().Get("code")
	if p.code != "" {
		p.top(w, r)
	} else {
		p.auth(w, r)
	}
}

func (p *GoogleDrive) createFile(name string) {
	_, err := p.service.Files.Create(&drive.File{Name: name}).Do()
	if err != nil {
		log.Errorf(p.c, "Error Create File %v", err)
	}
}

func (p *GoogleDrive) createDir(name string) {
	_, err := p.service.Files.Create(&drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}).Do()
	if err != nil {
		log.Errorf(p.c, "Error Create Dir %v", err)
	}
}

func (p *GoogleDrive) deleteFile(name string) {
	p.service.Files.Delete(name).Do()
}

func (p *GoogleDrive) printFile(num int64, w http.ResponseWriter) {
	d, err := p.service.Files.List().PageSize(num).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Errorf(p.c, "Unable to retrieve files. %v", err)
	}

	if len(d.Files) > 0 {
		for _, i := range d.Files {
			fmt.Fprint(w, "<div>"+i.Name+":"+i.Id+"</div>")
			p.deleteFile(i.Id)
		}
	} else {
		log.Infof(p.c, "No files found.")
	}
}

func (p *GoogleDrive) top(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var err error
	p.service, err = p.getService(r)
	if err != nil {
		log.Errorf(c, "Unable to retrieve drive Client %v", err)
	}

	p.createDir("alal")
	// p.printFile(40, w)
}

func (p *GoogleDrive) getService(r *http.Request) (*drive.Service, error) {
	c := appengine.NewContext(r)

	tok, err := p.config.Exchange(c, p.code)
	if err != nil {
		log.Errorf(c, "Unable to retrieve token from web %v", err)
	}

	client := p.config.Client(c, tok)
	return drive.New(client)
}

func (p *GoogleDrive) auth(w http.ResponseWriter, r *http.Request) {
	p.c = appengine.NewContext(r)
	log.Infof(p.c, "auth start")

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Errorf(p.c, "Unable to read client secret file: %v", err)
	}

	// config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	p.config, err = google.ConfigFromJSON(b, drive.DriveMetadataScope, drive.DriveFileScope, drive.DriveAppdataScope, drive.DriveFileScope)
	if err != nil {
		log.Errorf(p.c, "Unable to parse client secret file to config: %v", err)
	}
	url := p.config.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
